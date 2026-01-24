# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**hagg** is a deliberately minimal baseline for server-side rendered Go web applications using the HAGG stack:
- **H**TMX - HTML over the wire
- **A**lpine.js - client-side state management
- **G**omponents - type-safe HTML components in Go
- **G**o (Chi router) - stdlib-compatible HTTP routing

This is a **work-in-progress hobby project** focused on learning and exploration. The journey matters as much as the destination - decisions are documented, trade-offs are discussed, and the architecture evolves openly.

## Current Status: Phase 3 Complete ✅

The project is fully migrated to Chi router and Bootstrap 5.3:
- ✅ All routes migrated to Chi
- ✅ Gin server removed
- ✅ Duplicate code eliminated (~500 lines)
- ✅ Chi is the sole server on configured port
- ✅ Tailwind CSS replaced with Bootstrap 5.3 (CDN)
- ✅ No CSS build step required


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

### CSS
Bootstrap 5.3 is loaded via CDN - no build step required for CSS.

Custom overrides are in `static/css/app.css`.

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
- Enforcer and Perms injected into `app.Deps`

**Middleware:**
- `middleware.RequireAuth(wrapper)` - Requires authenticated user
- `middleware.RequirePermission(deps.Auth, deps.Users, deps.Perms, "action")` - Requires authentication + permission

**Permission format:** Actions as strings (e.g., `dashboard:view`, `user:create`, `user:list`, `user:delete`)

> **Note:** RequirePermission is a reference implementation. Adapt or remove it based on your needs.
> For simple apps where "logged in = full access", just use RequireAuth and remove Casbin entirely.

### Toast Notification System

> ⚠️ **IMPORTANT: Architecture Decision**
>
> The HTMX toast listener **MUST be on the `<body>` element**, not on individual forms!
> See section below for details.

**Backend (hagg-lib/toast):**
```go
// Fluent API for creating toasts
ctx.Toast("User created").Success().Notify()
ctx.Toast("Error occurred").Error().Stay().Notify()
ctx.Toast("Warning").Warning().SetTimeout(3000).Notify()
```

**Frontend - Two Delivery Mechanisms:**

1. **HTMX Requests** → `HX-Trigger` header with `toast` event
2. **Full Page Loads** → `RenderEvents()` generates self-destructing `<script>` tags

**Toast Listener Architecture (CRITICAL!):**

The toast listener is registered **once globally on `<body>`** in `skeleton.go`:

```go
Body(
    // Global HTMX toast listener - catches toast events from ANY HTMX request
    // IMPORTANT: Must be on <body>, not on individual forms!
    hx.On("toast", "showToast(event.detail)"),

    grp,
)
```

**WHY on `<body>` and not on individual forms?**

1. **HTMX doesn't bubble** - Events fire only on the element that made the request
2. **One listener is enough** - The `<body>` is never swapped by HTMX
3. **No duplicates** - Multiple listeners = multiple toasts per event
4. **Future-proof** - New forms don't need their own listener

**WRONG** (causes missing or duplicate toasts):
```go
// ❌ DON'T DO THIS - Listener on each form
Form(
    hx.Post("/htmx/login"),
    hx.On("toast", "showToast(event.detail)"),  // NO!
    ...
)
```

**CORRECT:**
```go
// ✅ DO THIS - One listener on <body> in skeleton.go
Body(
    hx.On("toast", "showToast(event.detail)"),
    grp,
)
```

**Frontend (static/js/toast.js):**
- `showToast({ message, level, timeout, position })` - Displays toast notification
- Bootstrap Icons for success, error, warning, info
- Auto-dismiss after timeout (default: 5000ms)
- Fade-in/out animation

### CSS & Styling

**Bootstrap 5.3** via CDN with minimal custom overrides:
- Bootstrap CSS loaded from CDN in `skeleton.go`
- Bootstrap Icons loaded from CDN
- Custom overrides in `static/css/app.css`
- Dark mode via Bootstrap's native `data-bs-theme` attribute
- Alpine.js persistence for theme preference

**No build step required** - All CSS is loaded directly from CDN

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
  middleware/               # Chi middleware (auth.go, permission.go, chi.go)
  session/                  # SCS session manager
  ucli/                     # CLI commands
  user/                     # User domain model & store
    store_sqlite/           # SQLite implementation
migrations/                 # SQL migrations (numbered files)
static/                     # Static assets (CSS, JS, images)
  css/                      # Custom CSS overrides (app.css)
  js/                       # Frontend JavaScript
model.conf                  # Casbin RBAC model
policy.csv                  # Casbin policies
server.go                   # Server startup, buildRouter()
routes.go                   # Route definitions
```

## Development Guidelines

### When Adding New Routes

1. Add route to Chi router in `routes.go`
2. Use `wrapper.Wrap(handlerFunc)` to wrap your handler
3. Handler signature: `func(*handler.Context) error`
4. Use Chi URL param extraction: `chi.URLParam(ctx.Req, "id")`
5. Test the route

**Example:**
```go
// In routes.go
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

### Default Middleware Stack

The application uses a layered middleware stack configured in `server.go:buildRouter()`:

```go
r.Use(chimiddleware.RealIP)          // Extract real IP from proxy headers
r.Use(chimiddleware.Compress(5))     // Gzip compression
r.Use(session.Manager.LoadAndSave)   // SCS session management (MUST be early)
r.Use(middleware.Recovery(wrapper))  // Panic recovery with logging
r.Use(middleware.Logger(wrapper))    // Request logging
r.Use(middleware.CORS())             // CORS headers for HTMX
r.Use(middleware.RateLimit)          // Rate limiting (placeholder)
r.Use(libmw.Secure)                  // Security headers (from hagg-lib)
```

**Security Headers (libmw.Secure):**
- `X-Frame-Options: DENY` - Prevents clickjacking
- `X-Content-Type-Options: nosniff` - Prevents MIME sniffing
- `X-XSS-Protection: 1; mode=block` - Legacy XSS filter
- `Referrer-Policy: strict-origin-when-cross-origin` - Limits referrer info

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

Use `go run ./cmd config` to print active configuration.

## Dependencies

### Core Backend
- `github.com/go-chi/chi/v5` - HTTP router
- `maragu.dev/gomponents` - Type-safe HTML components
- `github.com/alexedwards/scs/v2` - Session management
- `github.com/casbin/casbin/v2` - Authorization/RBAC
- `github.com/jmoiron/sqlx` - Database utilities
- `github.com/mattn/go-sqlite3` - SQLite driver

### Frontend (CDN, not in go.mod)
- Bootstrap 5.3 (CSS + JS)
- Bootstrap Icons
- HTMX 2.x
- Alpine.js 3.x
- surreal.js (DOM utilities)

### Development Tools
- `github.com/cosmtrek/air` - Hot reload (optional)
- just (task runner, optional)

## Key Files to Understand

- `ARCHITECTURE.md` - Deep architectural documentation
- `server.go` - Server startup logic
- `routes.go` - Chi route definitions
- `internal/session/manager.go` - SCS session manager
- `internal/middleware/auth.go` - RequireAuth, RequireGuest middleware
- `internal/middleware/permission.go` - RequirePermission middleware (Casbin-based)
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
