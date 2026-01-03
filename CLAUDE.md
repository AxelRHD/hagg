# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**hagg** is a deliberately minimal baseline for server-side rendered Go web applications using the HAGG stack:
- **H**TMX - HTML over the wire
- **A**lpine.js - client-side state management
- **G**omponents - type-safe HTML components in Go
- **G**o (Chi router) - stdlib-compatible HTTP routing

This is a **work-in-progress hobby project** focused on learning and exploration. The journey matters as much as the destination - decisions are documented, trade-offs are discussed, and the architecture evolves openly.

## Current Status: Phase 2 Complete ✅

The Gin → Chi migration is **complete**. The project now uses Chi router exclusively:
- ✅ All routes migrated to Chi
- ✅ Gin server removed
- ✅ Duplicate code eliminated (~500 lines)
- ✅ Chi is the sole server on configured port
- ✅ All dependencies cleaned up

See `REFACTORING_PLAN.md` for details on completed phases.

## Build & Development Commands

### Running the Application
```bash
# Start the server (Chi server on configured port, default :8080)
go run ./cmd

# Or use the CLI wrapper
go run cmd/main.go serve

# With air for hot-reload
air
```

### CSS Development
```bash
# Build CSS (production, minified)
just css-build

# Watch CSS (development, auto-rebuild)
just css-watch

# Or directly with Tailwind CLI
tailwindcss -i ./static/css/base.css -o ./static/css/styles.css --watch
```

### Testing
```bash
# Run all tests
go test ./...
just test
```

### Building
```bash
# Build binary (outputs to ./bin/)
go build -o bin/hagg cmd/main.go
just build
```

### Database Migrations
Migrations are located in `migrations/` directory and use raw SQL files. They are applied via the internal migration system during server startup.

## Architecture Overview

### Multi-Repo Structure

This project is split into two repositories:
- **hagg-lib** (`github.com/axelrhd/hagg-lib`) - Reusable library for context wrapper, toast system, events, middleware
- **hagg** (this repo) - Boilerplate project that imports hagg-lib

During development (WIP phase), hagg-lib is replaced via `go.mod`:
```go
replace github.com/axelrhd/hagg-lib => ../hagg-lib
```

This allows users who fork `hagg` to update `hagg-lib` without merge conflicts.

### Request Flow

```
HTTP Request → Chi Router → Chi Middleware → handler.Wrapper → handler.Context → Handler → Gomponents → HTML Response
                                                   ↓
                                          (automatic event commit)
```

### Context Wrapper Pattern

The Chi architecture uses a custom `handler.Context` wrapper (from hagg-lib):

```go
// handler.Context wraps http.ResponseWriter and *http.Request
type Context struct {
    Res http.ResponseWriter  // Explicit field (no embedding)
    Req *http.Request        // Explicit field (no embedding)
    // internal fields for events, logger, etc.
}

// Handler signature
type HandlerFunc func(*Context) error

// Wrapper converts HandlerFunc to stdlib http.HandlerFunc
func (w *Wrapper) Wrap(h HandlerFunc) http.HandlerFunc
```

**Key helper methods:**
- `ctx.Render(node g.Node) error` - Render Gomponents HTML
- `ctx.Toast(msg string) *toast.Builder` - Create toast notifications
- `ctx.Event(name string, data any)` - Emit events

### Event System

The project uses a **server-driven event bus** with two delivery mechanisms:

1. **HTMX Events** - Via `HX-Trigger` header for HTMX requests
2. **Initial Events** - Via `<script id="initial-events">` for full page loads

Both trigger the same JavaScript event handlers, creating a unified flow.

**Common events:**
- `toast` - Show notification (success, error, warning, info)
- `auth-changed` - Login/logout triggers UI refresh
- `permission-denied` - Authorization failures

**Backend usage:**
```go
ctx.Toast("Welcome back!").Success().Notify()
ctx.Event("auth-changed", nil)
```

Events are automatically committed by `handler.Wrapper` - no manual commit needed.

### Session Management

The project uses **SCS sessions** (alexedwards/scs/v2) with SQLite backend by default.

The SCS session manager is initialized in `internal/session/manager.go` and uses SQLite (configured via `SESSION_DB_PATH`).

**SCS session usage:**
```go
// Put value
session.Manager.Put(ctx.Req.Context(), "user_id", 123)

// Get value
userID := session.Manager.GetInt(ctx.Req.Context(), "user_id")

// Pop value (flash-style)
msg := session.Manager.PopString(ctx.Req.Context(), "flash_success")
```

### Authentication & Authorization

**Authentication:**
- Session-based auth via SCS
- User ID stored in session key `uid`
- Simple, no JWT or OAuth complexity
- See `internal/auth/auth.go`

**Authorization:**
- **Casbin** for RBAC/ABAC policies
- Model: `model.conf` in project root
- Policies: `policy.csv` in project root
- Enforcer injected into `app.Deps`

**Middleware:**
- `middleware.RequireAuth(wrapper)` - Requires authenticated user
- `middleware.RequirePermission(wrapper, "action")` - Requires specific permission

**Permission format:** Actions as strings (e.g., `user:create`, `user:list`, `user:delete`)

### Toast Notification System

> ⚠️ **WICHTIG: Architektur-Entscheidung**
>
> Der HTMX Toast-Listener **MUSS auf dem `<body>`-Element** sitzen, nicht auf einzelnen Formularen!
> Siehe Abschnitt unten für Details.

**Backend (hagg-lib/toast):**
```go
// Fluent API for creating toasts
ctx.Toast("User created").Success().Notify()
ctx.Toast("Error occurred").Error().Stay().Notify()
ctx.Toast("Warning").Warning().SetTimeout(3000).Notify()
```

**Frontend - Zwei Delivery-Mechanismen:**

1. **HTMX-Requests** → `HX-Trigger` Header mit `toast` Event
2. **Full-Page-Loads** → `RenderEvents()` erzeugt selbst-zerstörende `<script>`-Tags

**Toast-Listener-Architektur (KRITISCH!):**

Der Toast-Listener wird **einmal global auf dem `<body>`** registriert in `skeleton.go`:

```go
Body(
    // Global HTMX toast listener - catches toast events from ANY HTMX request
    // IMPORTANT: Must be on <body>, not on individual forms!
    hx.On("toast", "showToast(event.detail)"),

    grp,
)
```

**WARUM auf `<body>` und nicht auf einzelnen Forms?**

1. **HTMX bubblet nicht** - Events triggern nur auf dem Element, das den Request gemacht hat
2. **Ein Listener reicht** - Der `<body>` wird nie per HTMX geswappt
3. **Keine Duplikate** - Mehrere Listener = mehrere Toasts pro Event
4. **Zukunftssicher** - Neue Forms brauchen keinen eigenen Listener

**FALSCH** (führt zu fehlenden oder doppelten Toasts):
```go
// ❌ NICHT SO - Listener auf jedem Form
Form(
    hx.Post("/htmx/login"),
    hx.On("toast", "showToast(event.detail)"),  // NEIN!
    ...
)
```

**RICHTIG:**
```go
// ✅ SO - Ein Listener auf <body> in skeleton.go
Body(
    hx.On("toast", "showToast(event.detail)"),
    grp,
)
```

**Frontend (static/js/toast.js):**
- `showToast({ message, level, timeout, position })` - Zeigt Toast an
- SVG-Icons für success, error, warning, info
- Auto-dismiss nach Timeout (default: 5000ms)
- Fade-in/out Animation

### CSS & Styling

**Tailwind CSS** with **Pico.css-inspired aesthetics**:
- Design system defined in `tailwind.config.js`
- Base styles in `static/css/base.css`
- Compiled to `static/css/styles.css`
- Warm color palette, soft shadows, clean typography
- Dark mode support via manual toggle + `prefers-color-scheme`

**No npm/node required** - Uses Tailwind standalone CLI binary installed to `~/.local/bin/tailwindcss`

## Project Structure

```
cmd/                        # CLI entry point
internal/
  app/                      # Dependency container (Deps)
  auth/                     # Session-based authentication
  config/                   # Environment config loading
  frontend/                 # Gomponents layouts & pages
    layout/                 # Shared layout components
    pages/                  # Page components (login, dashboard, etc.)
  middleware/               # Chi middleware (auth, permissions, etc.)
  server/                   # Chi server setup
  session/                  # SCS session manager
  ucli/                     # CLI commands
  user/                     # User domain model & store
    store_sqlite/           # SQLite implementation
migrations/                 # SQL migrations (numbered files)
static/                     # Static assets (CSS, JS, images)
  css/                      # Tailwind source & output
  js/                       # Frontend JavaScript
model.conf                  # Casbin RBAC model
policy.csv                  # Casbin policies
server.go                   # Server startup logic
```

## Development Guidelines

### When Adding New Routes

1. Add route to Chi router in `internal/server/chi_routing.go` (or create new routing file)
2. Use `wrapper.Wrap(handlerFunc)` to wrap your handler
3. Handler signature: `func(*handler.Context) error`
4. Use Chi URL param extraction: `chi.URLParam(ctx.Req, "id")`
5. Test the route

**Example:**
```go
// In internal/server/chi_routing.go
func AddUserRoutes(r chi.Router, wrapper *handler.Wrapper, deps app.Deps) {
    r.Get("/users/{id}", wrapper.Wrap(func(ctx *handler.Context) error {
        id := chi.URLParam(ctx.Req, "id")
        // ... handler logic
        return ctx.Render(userPage)
    }))
}
```

### When Working with Sessions

- Use SCS session manager: `session.Manager` (from `internal/session`)
- Session operations require `context.Context` from `ctx.Req.Context()`
- Flash messages: Use `PopString()` to retrieve and remove in one operation

### When Emitting Events

- Use `ctx.Event(name, payload)` or `ctx.Toast(msg).Success().Notify()`
- Events are automatically committed by `handler.Wrapper`
- **Do NOT** manually call `hxevents.Commit()` - it's handled automatically

### When Creating Middleware

Chi middleware signature:
```go
func MyMiddleware(wrapper *handler.Wrapper) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // ... middleware logic
            next.ServeHTTP(w, r)
        })
    }
}
```

### File Organization

- **Pages/layouts** go in `internal/frontend/`
- **Middleware** go in `internal/middleware/`
- **Domain logic** goes in `internal/<domain>/` (e.g., `internal/user/`)
- **Stores** go in `internal/<domain>/store_<type>/` (e.g., `internal/user/store_sqlite/`)
- **Shared utilities** go in `internal/shared/` (if needed)

## Configuration

Configuration is loaded from environment variables and optional `.env` file:
- `SERVER_*` - Server config (PORT, BASE_PATH, SOCKET)
- `SESSION_*` - Session config (SECRET, MAX_AGE, DB_PATH, COOKIE_NAME)
- `DB_*` - Database config (PATH, DEBUG)
- `CASBIN_*` - Casbin config (MODEL_PATH, POLICY_PATH)

Use `go run ./cmd -config` to print active configuration.

## Dependencies

### Core Backend
- `github.com/go-chi/chi/v5` - HTTP router
- `maragu.dev/gomponents` - Type-safe HTML components
- `github.com/alexedwards/scs/v2` - Session management
- `github.com/casbin/casbin/v2` - Authorization/RBAC
- `github.com/jmoiron/sqlx` - Database utilities
- `github.com/mattn/go-sqlite3` - SQLite driver

### Frontend (CDN, not in go.mod)
- HTMX 2.x
- Alpine.js 3.x
- surreal.js (DOM utilities)

### Development Tools
- `github.com/cosmtrek/air` - Hot reload (optional)
- Tailwind CSS standalone CLI
- just (task runner, optional)

## Key Files to Understand

- `REFACTORING_PLAN.md` - Complete refactoring roadmap and decisions
- `MIGRATION.md` - Migration guide for users upgrading from old Gin-based hagg
- `ARCHITECTURE.md` - Deep architectural documentation
- `server.go` - Server startup logic
- `internal/server/chi_routing.go` - Chi route definitions
- `internal/session/manager.go` - SCS session manager
- `internal/middleware/chi.go` - Chi middleware stack
- `internal/auth/auth.go` - Authentication logic

## Handler Pattern Reference

### URL Parameters
```go
id := chi.URLParam(ctx.Req, "id")
```

### Query Parameters
```go
filter := ctx.Req.URL.Query().Get("filter")
```

### Form Data
```go
if err := ctx.Req.ParseForm(); err != nil {
    return err
}
email := ctx.Req.FormValue("email")
```

### JSON Request/Response
```go
// Decode JSON request
var data MyStruct
if err := json.NewDecoder(ctx.Req.Body).Decode(&data); err != nil {
    return err
}

// Encode JSON response
ctx.Res.Header().Set("Content-Type", "application/json")
return json.NewEncoder(ctx.Res).Encode(responseData)
```

### Rendering HTML (Gomponents)
```go
return ctx.Render(myPage())
```

## Philosophy & Approach

This project prioritizes:
- **Server-side HTML** as the default
- **Progressive enhancement** over JavaScript-heavy SPAs
- **Clear boundaries** between layers
- **Replaceable parts** over tight coupling
- **Pragmatic defaults** over academic perfection
- **Learning by doing** over following a fixed roadmap

The codebase is intentionally minimal and designed to be forked, adapted, and customized. It's not a framework - it's a stable baseline meant to be copied and evolved independently.

## Migration Notes (For Context)

The project has completed migration from Gin to Chi. If you see references to Gin in old commits or discussions, they refer to the previous architecture. The current architecture is Chi-only.

**Key changes in Phase 2:**
- Gin router → Chi router
- gin.Context → handler.Context
- gin-sessions → SCS sessions
- Consolidated all duplicate code
- Removed ~500 lines of duplication
