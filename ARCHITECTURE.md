# Architecture

This document describes the **concrete architecture** of the `hagg` project.
It explains not only *how* the code is structured, but more importantly *why*.

Throughout this document, **we** refers to the authors and maintainers of this project —
this is our shared baseline.

> **Project status:** `hagg` is still early **WIP (work in progress)**.
> The structure is intentionally simple, and we expect it to evolve as we learn.
>
> **The journey matters:** This architecture is the result of exploration, discussion,
> and iteration. We document our decisions and trade-offs openly.
> If something seems unconventional, there's likely a reason — read on!

---

## Core Idea

The application follows a **server-centric architecture**:

- HTML is rendered on the server (Gomponents)
- Interaction happens over HTTP (HTMX)
- surreal.js provides syntactic sugar for DOM operations
- Alpine.js is used for small client-side state only
- The backend stays explicit and boring (Chi, middleware, plain Go)

This baseline exists to make "normal web apps" feel straightforward again:
forms, tables, and stateful pages — without SPA complexity.

---

## Context (important)

We are **not professional web developers**, nor specialists in frontend or backend frameworks.
This is a **hobby project**, created to explore ideas, deepen understanding,
and learn by building real things.

So this architecture should be read as:

- Honest and practical
- Experience-driven, not academic
- Open to correction and improvement as our understanding grows
- **Process-oriented** — the journey and learning are as important as the output

---

## High-Level Request Flow

1. **Chi** matches a route.
2. **Middleware** runs (HTMX triggers, auth, permissions).
3. **Wrapper** creates a `handler.Context` and calls the handler.
4. A **handler** renders Gomponents into `handler.Context`.
5. The **wrapper** commits events (via `HX-Trigger` header or initial-events script).
6. HTMX endpoints return partial HTML or trigger client-side updates.
7. **Frontend** processes events (toast, auth-changed, etc.) with surreal.js.

The important part: *handlers are simple functions* that work with a lightweight context wrapper.

```mermaid
sequenceDiagram
    participant B as Browser
    participant C as Chi Router
    participant M as Middleware
    participant W as Wrapper
    participant H as Handler
    participant R as Renderer (Gomponents)

    B->>C: HTTP request
    C->>M: HTMX triggers, auth, permissions
    M-->>C: ok or abort
    C->>W: create context
    W->>H: call handler(ctx)
    H->>R: render gomponents tree
    R-->>W: HTML
    W-->>B: HTML + event headers (HX-Trigger)
```

---

## Stack Components in Detail

### Backend

#### Chi Router

We use **Chi v5** as our HTTP router.

**Why Chi?**

- stdlib-compatible (`http.Handler`, `http.HandlerFunc`)
- Minimal, no framework overhead
- Composable middleware
- Idiomatically Go

**Example routing (from `routes.go`):**

```go
func AddRoutes(r chi.Router, wrapper *handler.Wrapper, deps app.Deps) {
    // Public routes
    r.Get("/", wrapper.Wrap(home.Page(deps)))
    r.Get("/login", wrapper.Wrap(login.Page(deps)))

    // HTMX endpoints
    r.Post("/htmx/login", wrapper.Wrap(login.HxLogin(deps)))
    r.Post("/htmx/logout", wrapper.Wrap(login.HxLogout(deps)))

    // Protected routes (require authentication + permission)
    r.Group(func(r chi.Router) {
        r.Use(middleware.RequirePermission(deps.Auth, deps.Users, deps.Perms, "dashboard:view"))
        r.Get("/dashboard", wrapper.Wrap(dashboard.Page(deps)))
    })
}
```

#### Context Wrapper

We use a **lightweight context wrapper** instead of frameworks like Gin or Echo.

**The `handler.Context` type:**

```go
type Context struct {
    Res http.ResponseWriter  // Explicit field (no embedding)
    Req *http.Request        // Explicit field (no embedding)

    logger *slog.Logger
    events []Event
}
```

**Why explicit fields?**

- No embedding = no interface pollution
- Compatible with all stdlib middleware
- Clear ownership and access patterns

**Helper methods (minimal):**

```go
func (c *Context) Render(node g.Node) error
func (c *Context) Toast(msg string) *toast.Builder
func (c *Context) Event(name string, data any)
func (c *Context) Logger() *slog.Logger
```

**Handler pattern:**

```go
type HandlerFunc func(*Context) error

func (w *Wrapper) Wrap(h HandlerFunc) http.HandlerFunc {
    return func(res http.ResponseWriter, req *http.Request) {
        ctx := &Context{
            Res:    res,
            Req:    req,
            logger: w.logger,
            events: make([]Event, 0),
        }

        if err := h(ctx); err != nil {
            w.logger.Error("handler error", "error", err)
            http.Error(res, "Internal Server Error", 500)
        }

        // Commit events (HX-Trigger or initial-events)
    }
}
```

**Usage in handlers (Factory Pattern):**

Handlers use the factory pattern — a function that takes dependencies and returns a `handler.HandlerFunc`:

```go
func Page(deps app.Deps) handler.HandlerFunc {
    return func(ctx *handler.Context) error {
        user, _ := deps.Auth.CurrentUser(ctx.Req)

        content := Div(
            H1(g.Text("Dashboard")),
            P(g.Textf("Welcome, %s!", user.FullName())),
        )

        return ctx.Render(layout.Page(ctx, deps, content))
    }
}
```

**Why the factory pattern?**

- Dependencies are captured in the closure — no global state
- Handler signature is clean (`func(*Context) error`)
- Easy to test (inject mock dependencies)

#### Session Management (scs)

We use **alexedwards/scs** for session management.

**Why scs?**

- stdlib-compatible
- Multiple store backends (cookie, SQLite, Postgres, Redis)
- Simple API
- Actively maintained

**Default: SQLite-backed sessions (persistent)**

Sessions are stored in SQLite, meaning they persist across server restarts:

```go
// internal/session/manager.go
func Init(dbPath string) error {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return err
    }

    Manager = scs.New()
    Manager.Store = sqlite3store.New(db)
    Manager.Lifetime = 24 * time.Hour
    Manager.Cookie.Name = "hagg_session"
    Manager.Cookie.HttpOnly = true
    Manager.Cookie.Secure = false // Set to true in production (HTTPS)

    return nil
}
```

**Global Manager Pattern:**

The session manager is accessed globally via `session.Manager`:

```go
// Get value
uid := session.Manager.GetString(ctx.Req.Context(), "uid")

// Set value
session.Manager.Put(ctx.Req.Context(), "uid", userUID)

// Flash messages (pop = get and remove)
msg := session.Manager.PopString(ctx.Req.Context(), "flash_success")
```

**Session Middleware:**

The session middleware must be registered early in the middleware stack:

```go
r.Use(session.Manager.LoadAndSave)  // MUST come before auth middleware
```

### Frontend

#### HTMX

HTMX powers progressive enhancement:

- Form submissions without JavaScript
- Partial page updates
- Polling and lazy loading
- WebSocket support

**Example:**

```html
<button hx-post="/logout" hx-swap="none">Logout</button>
```

The backend responds with `HX-Trigger: auth-changed`, and the frontend refreshes the nav.

#### Alpine.js

Alpine.js handles **client-side state** (dropdowns, modals, tabs):

```html
<div x-data="{ open: false }">
    <button @click="open = !open">Toggle</button>
    <div x-show="open">Content</div>
</div>
```

**Important:** Alpine is for **UI state only**, not application state.

#### surreal.js

surreal.js provides **syntactic sugar** for DOM operations:

```js
// Vanilla JS
document.getElementById('toast-container')

// surreal.js
me('#toast-container')
```

**Why surreal.js?**

- Reduces boilerplate in event handlers
- Chainable API (`me(el).classAdd('foo').classRemove('bar')`)
- Small (~2KB minified)
- No build step

**Example (toast rendering):**

```js
function showToast({ message, level = 'info' }) {
    const container = me('#toast-container-bottom-right')

    const toastHtml = `
        <div class="toast toast-${level} opacity-0">
            ${message}
        </div>
    `

    container.insertAdjacentHTML('beforeend', toastHtml)
    const toast = container.lastElementChild

    me(toast).classRemove('opacity-0')

    setTimeout(() => {
        me(toast).classAdd('opacity-0')
        setTimeout(() => me(toast).remove(), 300)
    }, 5000)
}
```

#### Tailwind CSS (Pico-inspired)

We use **Tailwind CSS** with a **Pico.css-inspired design system**.

**Why Tailwind?**

- Full control over styling
- No external CSS dependencies
- Purging for small bundle size

**Why Pico-inspired?**

- Clean, minimal aesthetics
- Good defaults for forms, buttons, cards
- Professional without over-design

**Custom theme:**

```js
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      colors: {
        primary: { /* blue scale */ },
        success: '#10b981',
        error: '#ef4444',
        warning: '#f59e0b',
        info: '#3b82f6',
      },
      borderRadius: {
        'pico': '0.375rem',
      },
      boxShadow: {
        'pico': '0 0 0 1px rgba(0, 0, 0, 0.05), ...',
      },
    },
  },
}
```

**Build process (no npm/node):**

```bash
# Standalone binary (installed once)
tailwindcss -i ./static/css/base.css -o ./static/css/styles.css --watch
```

See `static/css/base.css` for component styles.

---

## Event System (Detailed)

The event system is **server-driven** and **header-based**.

### Two Event Types

#### 1. HTMX Events (via `HX-Trigger` header)

For HTMX requests, events are sent via the `HX-Trigger` response header:

```http
HX-Trigger: {"toast": {"message": "Success!", "level": "success"}}
```

The frontend listens with:

```js
htmx.on('toast', (e) => {
    showToast(e.detail)
})
```

#### 2. Initial Events (via `<script id="initial-events">`)

For full page loads, events are rendered as JSON in the HTML:

```html
<script type="application/json" id="initial-events">
[
  {"name": "toast", "payload": {"message": "Welcome!", "level": "info"}},
  {"name": "auth-changed", "payload": null}
]
</script>
```

The frontend processes on load:

```js
document.addEventListener('DOMContentLoaded', () => {
    const initialEventsEl = me('#initial-events')
    if (initialEventsEl) {
        const events = JSON.parse(initialEventsEl.textContent)
        events.forEach(processEvent)
        initialEventsEl.remove()
    }
})
```

### Unified Event Processing

Both types feed into the same event processor:

```js
function processEvent(event) {
    switch(event.name) {
        case 'toast':
            showToast(event.payload)
            break
        case 'auth-changed':
            htmx.trigger(me('#nav'), 'refresh')
            break
        case 'permission-denied':
            // ... handle permission denial
            break
    }
}
```

### Backend Event API

**Emit an event:**

```go
ctx.Event("auth-changed", nil)
```

**Toast shorthand:**

```go
ctx.Toast("Operation successful!").Success().Notify()
```

**Under the hood (toast):**

```go
func (t *Toast) Notify() {
    t.ctx.Event("toast", t)
}
```

### Event Flow Diagram

```mermaid
sequenceDiagram
    participant H as Handler
    participant C as Context
    participant W as Wrapper
    participant B as Browser
    participant F as Frontend (JS)

    H->>C: ctx.Toast("Success!").Success().Notify()
    C->>C: append Event{"toast", {...}}
    H->>W: return nil (success)

    alt HTMX Request
        W->>B: HX-Trigger: {"toast": {...}}
        B->>F: htmx.on('toast', ...)
        F->>F: showToast(...)
    else Full Page Load
        W->>B: HTML with <script id="initial-events">
        B->>F: DOMContentLoaded → process events
        F->>F: showToast(...)
    end
```

---

## Toast / Notification System

The toast system is **event-driven** and **unified**.

### Backend (Go)

**Toast builder:**

```go
type Toast struct {
    Message  string `json:"message"`
    Level    string `json:"level"`    // success, error, warning, info
    Timeout  int    `json:"timeout"`  // ms, 0 = stay forever
    Position string `json:"position"` // bottom-right, top-right, etc.
}
```

**Fluent API:**

```go
ctx.Toast("User created")
    .Success()
    .SetTimeout(3000)
    .SetPosition("top-right")
    .Notify()
```

### Frontend (JavaScript)

**Toast rendering (with surreal.js):**

```js
function showToast({ message, level = 'info', timeout = 5000, position = 'bottom-right' }) {
    const container = me(`#toast-container-${position}`)

    const toastHtml = `
        <div class="toast toast-${level} opacity-0 transition-opacity">
            <div class="flex items-center gap-3">
                <div>${getIcon(level)}</div>
                <div>${message}</div>
            </div>
        </div>
    `

    container.insertAdjacentHTML('beforeend', toastHtml)
    const toast = container.lastElementChild

    // Enter animation
    setTimeout(() => me(toast).classRemove('opacity-0'), 10)

    // Auto-remove
    if (timeout > 0) {
        setTimeout(() => {
            me(toast).classAdd('opacity-0')
            setTimeout(() => me(toast).remove(), 300)
        }, timeout)
    }
}
```

### Why This Design?

- **Unified:** One `showToast()` function for HTMX and full-page loads
- **Server-driven:** Backend controls notification logic
- **No external libraries:** Custom implementation with Tailwind
- **Flexible:** Supports different positions, timeouts, levels

---

## Authentication

Authentication is **session-based** and deliberately minimal:

- We store the logged-in user UID in the session (key: `uid`)
- `Auth.CurrentUser(req)` reads the session and loads the user from `user.Store`
- See `internal/auth/auth.go` for implementation

**Auth struct:**

```go
type Auth struct {
    users user.Store
}

func New(users user.Store) *Auth {
    return &Auth{users: users}
}

func (a *Auth) CurrentUser(req *http.Request) (*user.User, error) {
    rawUID := session.Manager.Get(req.Context(), SessionKeyUID)
    uid, ok := rawUID.(string)
    if !ok || uid == "" {
        return nil, ErrNotAuthenticated
    }
    return a.users.FindByUID(req.Context(), uid)
}

func (a *Auth) IsAuthenticated(req *http.Request) bool {
    rawUID := session.Manager.Get(req.Context(), SessionKeyUID)
    uid, ok := rawUID.(string)
    return ok && uid != ""
}
```

**Middleware (RequireAuth):**

```go
// RequireAuth ensures user is logged in (UID in session)
func RequireAuth(wrapper *handler.Wrapper) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            rawUID := session.Manager.Get(r.Context(), auth.SessionKeyUID)
            uid, ok := rawUID.(string)

            if !ok || uid == "" {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

**RequireGuest middleware:**

For pages that should only be accessible to non-authenticated users (login, register):

```go
// RequireGuest ensures user is NOT logged in
func RequireGuest(wrapper *handler.Wrapper) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            rawUID := session.Manager.Get(r.Context(), auth.SessionKeyUID)
            uid, ok := rawUID.(string)

            if ok && uid != "" {
                http.Redirect(w, r, "/", http.StatusSeeOther)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

---

## Authorization (Permission Layer)

> **Important:** This section describes a **reference implementation** that you can adapt or discard
> based on your application's needs. Not every app needs fine-grained permissions!

### When You Need This

Use Casbin-based authorization when:

- **Multi-user apps** with different roles (admin, editor, viewer)
- **Fine-grained access control** (user:create, user:delete, report:export)
- **Permissions change without code** (policy file edits, no rebuild)

### When You DON'T Need This

Skip Casbin and remove it entirely when:

- **Single-user tools** — just use `RequireAuth`
- **Simple CRUD apps** — resource ownership checks in handlers are enough
- **"Logged in = full access"** — `RequireAuth` is sufficient

If your app only needs "logged in or not", use `RequireAuth` and delete the Casbin files.

### Casbin Integration

We use **Casbin** as the authorization engine:

- `model.conf` defines the evaluation model (RBAC)
- `policy.csv` defines roles/permissions and user-role assignments
- `casbinx.NewFileEnforcer()` from hagg-lib loads both files

**Initialization (in `server.go`):**

```go
enforcer, err := casbinx.NewFileEnforcer(
    cfg.Casbin.ModelPath,  // default: "model.conf"
    cfg.Casbin.PolicyPath, // default: "policy.csv"
)
if err != nil {
    log.Fatal(err)
}

deps := app.Deps{
    // ...
    Enforcer: enforcer,
    Perms:    casbinx.NewPerm(enforcer),
}
```

**Why file-based?**

When we ship a single binary, we can place `model.conf` and `policy.csv` next to it.
No database required for authorization.

### Permissions as Actions

We model permissions as plain action strings:

- `dashboard:view`
- `user:create`
- `user:list`
- `user:delete`

This keeps the policy readable and avoids "object explosion" early on.

**Example policy (`policy.csv`):**

```csv
# Superuser: can do everything
p, superuser, *

# Admin role
p, admin, dashboard:view
p, admin, user:create
p, admin, user:list
p, admin, user:delete

# Viewer role (read-only)
p, viewer, dashboard:view
p, viewer, user:list

# User → Role assignments
g, arudolf, superuser
g, alice, admin
g, worker, viewer
```

### Middleware: RequirePermission

The `RequirePermission` middleware combines authentication and authorization:

```go
func RequirePermission(authService *auth.Auth, users user.Store, perms *casbinx.Perm, action string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Step 1: Check authentication
            rawUID := session.Manager.Get(r.Context(), auth.SessionKeyUID)
            uid, ok := rawUID.(string)
            if !ok || uid == "" {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
            }

            // Step 2: Load user for authorization
            u, err := users.FindByUID(r.Context(), uid)
            if err != nil {
                http.Error(w, "User not found", http.StatusUnauthorized)
                return
            }

            // Step 3: Check permission via Casbin
            // Subject is DisplayName (adapt if your policy uses UID or email)
            if !perms.Can(u.DisplayName, action) {
                http.Error(w, "Permission denied", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

**Usage:**

```go
r.Group(func(r chi.Router) {
    r.Use(middleware.RequirePermission(deps.Auth, deps.Users, deps.Perms, "dashboard:view"))
    r.Get("/dashboard", wrapper.Wrap(dashboard.Page(deps)))
})
```

### Customization Notes

This implementation uses `DisplayName` as the Casbin subject. Depending on your needs:

- **UID as subject** — unique and stable, but less readable in policy.csv
- **Email as subject** — readable, but may change
- **Role in session** — skip user lookup entirely

Adapt the middleware to match your policy design.

### Flow Diagram

```mermaid
sequenceDiagram
    participant B as Browser
    participant M as Middleware
    participant S as Session
    participant U as UserStore
    participant C as Casbin
    participant H as Handler

    B->>M: request
    M->>S: Get(uid)
    alt not logged in
        M-->>B: redirect /login
    else logged in
        M->>U: FindByUID(uid)
        M->>C: Can(displayName, action)
        alt denied
            M-->>B: 403 Forbidden
        else allowed
            M->>H: next.ServeHTTP()
            H-->>B: HTML response
        end
    end
```

---

## Rendering Model

We use **Gomponents** for type-safe HTML rendering.

**Key idea:**

- A **page** is a factory function: `func(deps app.Deps) handler.HandlerFunc`
- The factory captures dependencies and returns a handler
- The handler renders Gomponents into `handler.Context`

**Example page (Factory Pattern):**

```go
// internal/frontend/pages/dashboard/page.go
func Page(deps app.Deps) handler.HandlerFunc {
    return func(ctx *handler.Context) error {
        user, _ := deps.Auth.CurrentUser(ctx.Req)

        content := Div(
            Class("container"),
            H1(g.Text("Dashboard")),
            P(g.Textf("Welcome, %s!", user.FullName())),
        )

        return ctx.Render(layout.Page(ctx, deps, content))
    }
}
```

**Render helper (from hagg-lib):**

```go
func (c *Context) Render(node g.Node) error {
    c.Res.Header().Set("Content-Type", "text/html; charset=utf-8")
    return node.Render(c.Res)
}
```

**Why Gomponents?**

- Type-safe (compile-time checking)
- No templates (no runtime parsing)
- Pure Go (no external DSL)
- Composable (functions returning nodes)

**Why the factory pattern?**

- Dependencies are captured at route registration time
- Handler signature stays clean (`func(*Context) error`)
- Easy to test — inject mock dependencies

---

## Directory Layout (Detailed)

```
server.go             # Server startup, buildRouter(), middleware stack
routes.go             # Route definitions (AddRoutes function)
model.conf            # Casbin RBAC model
policy.csv            # Casbin policies (roles → actions, users → roles)
justfile              # Task runner (dev, build, css-build, css-watch)
tailwind.config.js    # Tailwind CSS configuration

cmd/
  main.go             # CLI entry point (urfave/cli)

internal/
  app/
    deps.go           # Dependency container (Deps struct)

  auth/
    auth.go           # Session-based authentication

  config/
    config.go         # Environment config loading (.env support)

  db/
    sqlite.go         # Database connection setup

  frontend/
    layout/
      skeleton.go     # HTML skeleton (<html>, <head>, <body>)
      page.go         # Page wrapper with layout
      nav.go          # Navigation component
      events.go       # Initial-events renderer

    pages/
      home/           # Homepage
        page.go
      login/          # Login page + HTMX handlers
        page.go
        components.go
        handler.go
      dashboard/      # Protected dashboard
        page.go

  middleware/
    auth.go           # RequireAuth, RequireGuest
    permission.go     # RequirePermission (Casbin-based)
    chi.go            # Logger, Recovery, CORS, RateLimit

  session/
    manager.go        # SCS session manager (SQLite backend)

  ucli/
    serve.go          # CLI serve command
    user.go           # CLI user management

  user/
    model.go          # User domain model
    store.go          # Store interface
    store_sqlite/
      store.go        # SQLite implementation

migrations/
  001_initial.sql     # Initial schema
  002_sessions.sql    # SCS sessions table

static/
  css/
    base.css          # Tailwind input
    styles.css        # Tailwind output (generated, gitignored)

  js/
    app.js            # Main application logic
    toast.js          # Toast notification system
    surreal.min.js    # surreal.js library
    alpine.min.js     # Alpine.js (local copy)
    htmx.min.js       # HTMX (local copy)
```

---

## Composition Root

The composition root is split between `cmd/main.go` (CLI) and `server.go` (server setup).

**CLI (`cmd/main.go`):**

- Parse CLI flags and commands (urfave/cli)
- Load config
- Open database
- Run migrations
- Delegate to `hagg.StartServer()`

**Server (`server.go`):**

- Initialize session manager
- Build Chi router with middleware
- Wire dependencies
- Start HTTP server (TCP or Unix socket)

**buildRouter function (from `server.go`):**

```go
func buildRouter(cfg *config.Config, usrStore user.Store) http.Handler {
    // Create logger
    logger := slog.Default()

    // Create handler wrapper (from hagg-lib)
    wrapper := handler.NewWrapper(logger)

    // Casbin enforcer (from hagg-lib)
    enforcer, err := casbinx.NewFileEnforcer(
        cfg.Casbin.ModelPath,
        cfg.Casbin.PolicyPath,
    )
    if err != nil {
        log.Fatal(err)
    }

    // Dependencies
    deps := app.Deps{
        Users:    usrStore,
        Auth:     auth.New(usrStore),
        Enforcer: enforcer,
        Perms:    casbinx.NewPerm(enforcer),
    }

    // Create Chi router
    r := chi.NewRouter()

    // Middleware stack (order matters!)
    r.Use(chimw.RealIP)                    // Extract real IP from proxy headers
    r.Use(chimw.Compress(5))               // Gzip compression
    r.Use(session.Manager.LoadAndSave)     // SCS sessions (MUST be early!)
    r.Use(middleware.Recovery(wrapper))    // Panic recovery
    r.Use(middleware.Logger(wrapper))      // Request logging
    r.Use(middleware.CORS())               // CORS headers
    r.Use(middleware.RateLimit)            // Rate limiting
    r.Use(libmw.Secure)                    // Security headers

    // Static files
    fs := http.FileServer(http.Dir("./static"))
    r.Handle("/static/*", http.StripPrefix("/static/", fs))

    // Application routes
    AddRoutes(r, wrapper, deps)

    return r
}
```

**Deps struct (`internal/app/deps.go`):**

```go
type Deps struct {
    Users    user.Store
    Auth     *auth.Auth
    Enforcer *casbin.Enforcer
    Perms    *casbinx.Perm
}
```

This separation keeps CLI concerns separate from HTTP wiring.

---

## What We Expect to Change

As we add real pages and features, we expect:

- More page packages under `internal/frontend/pages`
- More store implementations (in-memory/dev vs sqlite vs postgres)
- Richer authorization (maybe object-based rules later)
- Better error handling and observability (structured logging, metrics)
- Improved event system (more event types, better error handling)

The goal is not stability at all costs, but **clarity during growth**.

---

## Why This Architecture?

### What We Optimized For

- **Simplicity** — minimal abstractions, boring code
- **Explicitness** — no magic, clear data flow
- **Replaceability** — swap persistence, sessions, auth without rewriting everything
- **Learning** — we document decisions, iterate, and improve as we learn

### What We Did NOT Optimize For

- **Performance at scale** — this is a baseline for small-to-medium apps
- **Enterprise patterns** — no hexagonal architecture, no CQRS, no event sourcing
- **Frontend richness** — no SPA, no reactive framework, no virtual DOM

### Trade-offs We Accept

- **More boilerplate** than frameworks (but we know what each line does)
- **Less magic** than ORMs (but we control our queries)
- **Manual CSS** with Tailwind (but we have full control)

---

## Influences & Inspirations

- **Go stdlib philosophy** — simple, explicit, composable
- **HTMX** — HTML over the wire, progressive enhancement
- **Pico.css** — clean, minimal design
- **Gomponents** — type-safe templates
- **Chi** — stdlib-compatible routing
- **alexedwards** — session management patterns (scs)
- **surreal.js** — syntactic sugar without build complexity

---

## Final Thoughts

This architecture is **not perfect**, and it's **not finished**.

It's the result of exploration, discussion, and iteration.
We document our decisions, discuss trade-offs, and learn by building.

If you're looking for production-ready stability, wait for v1.0.0.
If you want to understand *why* things work the way they do, welcome aboard!

**The journey matters as much as the destination.**
