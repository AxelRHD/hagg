# hagg

`hagg` is a deliberately minimal baseline for **server-side rendered Go web applications**.
It is designed for developers who want to build long-lived, maintainable tools — without SPA overhead,
without build pipelines, and without hidden magic.

This project is **not a framework**.
It is a **stable baseline** meant to be copied, adapted, and used for years.

Following the tradition of established stack acronyms such as JAMstack, MERN, or LAMP,
we refer to the stack used by `hagg` as the **HAGG Stack**:

**H**TMX · **A**lpine.js (state) · **G**omponents · **G**o (Chi router)

with **surreal.js** for syntactic sugar in DOM operations.

> **Project status:** This repository is still early **WIP (work in progress)**.
> We are building it as a hobby project to learn by doing — expect rough edges and changes.
>
> **Important:** This is not just about the output — the journey, the process, and the learning
> are core to this project. We document our decisions, discuss trade-offs, and iterate openly.
> If you're looking for a "production-ready framework," this might not be for you.
> If you want to understand *why* things work the way they do, welcome!

---

## Project Goals

- **Server-side HTML** as the default (simple, fast, inspectable)
- **Progressive enhancement** with HTMX (small, explicit interactions)
- **No "frontend build" requirement** for JavaScript (single Go binary + static assets)
- **Clear boundaries** (config, auth, storage, middleware, pages)
- **Replaceable parts** (swap persistence, session store, auth/authorization model)
- **Pragmatic defaults** over academic perfection
- **Learning by doing** — we document decisions and evolve the architecture as we learn

---

## The Stack in One Sentence

Chi routes requests → Gomponents render HTML → HTMX updates parts of the page → surreal.js simplifies DOM operations → Alpine.js adds tiny client-side state.

---

## Stack Components

### Backend

- **[Chi](https://github.com/go-chi/chi)** — stdlib-compatible HTTP router (v5)
- **[Gomponents](https://www.gomponents.com/)** — type-safe HTML components in Go
- **[alexedwards/scs](https://github.com/alexedwards/scs)** — session management (multiple stores)
- **[Casbin](https://casbin.org/)** — authorization / RBAC
- **SQLite** — default database (swappable)

### Frontend

- **[HTMX](https://htmx.org/)** — HTML over the wire, progressive enhancement
- **[Alpine.js](https://alpinejs.dev/)** — client-side state management (minimal)
- **[surreal.js](https://github.com/gnat/surreal)** — syntactic sugar for DOM operations
- **[Bootstrap 5.3](https://getbootstrap.com/)** — CSS framework via CDN (no build step)

---

## Quickstart

### Requirements

- Go 1.23+ (see `go.mod`)
- SQLite (default database)
- [just](https://github.com/casey/just) (task runner, optional but recommended)

### Setup

```bash
# Clone the repo
git clone https://github.com/axelrhd/hagg.git
cd hagg

# Install dependencies
go mod download

# Run the app
go run ./cmd
```

### Development Mode

Use [air](https://github.com/cosmtrek/air) for Go hot-reload:

```bash
air
```

No CSS build step required — Bootstrap is loaded via CDN.

---

## What's Included

The `hagg` boilerplate includes example pages demonstrating key patterns:

### Pages

- **Homepage (`/`)** — Public landing page with HAGG stack overview
- **Login (`/login`)** — UID-based authentication with session management
- **Dashboard (`/dashboard`)** — Protected page (requires authentication)

### Features Demonstrated

- **Session-based authentication** with SCS (persistent SQLite storage)
- **Protected routes** using Chi middleware (`middleware.RequireAuth`)
- **Conditional navigation** (shows/hides links based on auth status)
- **HTMX interactions** (form submissions, auto-refresh on auth-changed events)
- **Server-to-client events** (toast notifications via HX-Trigger headers)
- **Dark mode toggle** with Alpine.js state persistence
- **Type-safe HTML** rendering with gomponents
- **Responsive design** with Bootstrap 5.3

All pages use server-side rendering — no JavaScript build step required.

---

## Configuration

Configuration is loaded from environment variables (and an optional `.env`).

- Server config is prefixed with `SERVER_` (e.g. `SERVER_PORT`, `SERVER_BASE_PATH`)
- Session config is prefixed with `SESSION_`
- Database config is prefixed with `DB_` (see `internal/config` for details)

To print the active configuration:

```bash
go run ./cmd config
```

---

## Authentication

Authentication is intentionally simple:

- A session cookie stores the logged-in user's UID (`internal/auth`, session key `uid`)
- Pages / HTMX endpoints use that UID to load the current user from the store
- Session storage is pluggable (cookie-based by default, can use SQLite/Postgres/Redis)

See:

- `internal/auth/auth.go`
- `internal/middleware/auth.go` (RequireAuth middleware)
- `internal/frontend/pages/login/*`

---

## Authorization (Permission Layer)

`hagg` includes a permission layer built on **Casbin**.

### Why Casbin?

We do **not** want permissions hard-coded into handlers.
Changing roles or permissions should not require a new release.

Casbin gives us:

- RBAC-style policies (roles, users, inheritance)
- A clean enforcement API (`Enforce(subject, action)`)
- A simple file-based adapter for local projects

### Where it lives

- `model.conf` — the Casbin model
- `policy.csv` — the policy (roles → permissions, users → roles)

Both files live in the **project root** so that the compiled binary can find them at runtime.

### Policy structure

We model permissions as **actions** (strings), e.g.

- `user:create`
- `user:list`
- `user:delete`
- `selfdestroy`

Example (from `policy.csv`):

```csv
p, admin, user:create
p, admin, user:list

g, alice, admin
```

- `p` lines map **role → action**
- `g` lines map **user → role**

### Enforcement

The Casbin enforcer is initialized in `server.go:buildRouter()` using `casbinx.NewFileEnforcer()`
from hagg-lib. It loads `model.conf` and `policy.csv` from the working directory
and injects a `*casbin.Enforcer` into `app.Deps`.

### Middleware

`internal/middleware/auth.go` provides:

```go
middleware.RequireAuth(wrapper)   // Requires logged-in user
middleware.RequireGuest(wrapper)  // Requires NOT logged-in (e.g., login page)
```

Behavior:

- `RequireAuth`: If not authenticated → redirect to login
- `RequireGuest`: If already authenticated → redirect to dashboard

For permission-based access control, use Casbin's enforcer directly in handlers,
or extend the middleware as needed for your use case.

> **Tip:** In routes, compose middlewares like:
> `RequireAuth` → handler (for protected pages)

---

## Event System

`hagg` uses a **lightweight, server-driven event bus** implemented via HTTP response headers.

Instead of a client-side global event system, the backend emits **semantic events** that the frontend reacts to.

### Two types of events:

1. **HTMX Events** — sent via `HX-Trigger` header for HTMX requests
2. **Initial Events** — rendered as `<script id="initial-events">` for full page loads

Both are processed by the same JavaScript event handler, creating a unified flow.

### Typical Events:

- `toast` — show notification (success, error, warning, info)
- `auth-changed` — login or logout (triggers UI refresh)
- `permission-denied` — authorization failures

### Backend Usage:

```go
// Handler factory pattern (returns handler.HandlerFunc)
func HxLogin(deps app.Deps) handler.HandlerFunc {
    return func(ctx *handler.Context) error {
        // ... authentication logic ...

        ctx.Toast("Welcome back!").Success().Notify()
        ctx.Event("auth-changed", nil)

        return ctx.NoContent()  // HTMX will react to events
    }
}
```

See `ARCHITECTURE.md` for detailed event flow diagrams.

---

## Project Layout

A simplified overview:

```
server.go             # Server startup, buildRouter(), middleware stack
routes.go             # Route definitions (AddRoutes)
model.conf            # Casbin RBAC model
policy.csv            # Casbin policies
justfile              # Task runner (dev, build, css-build, css-watch)

cmd/                  # CLI entry point (urfave/cli)
internal/
  app/                # Dependency container (Deps struct)
  auth/               # Session-based authentication (SCS)
  config/             # Environment config loading (.env support)
  db/                 # Database connection setup
  frontend/           # Gomponents UI layer
    layout/           # Shared layout components (skeleton, nav, events)
    pages/            # Page handlers (home, login, dashboard)
  middleware/         # Chi middleware (auth, permissions, logging)
  session/            # SCS session manager (SQLite backend)
  ucli/               # CLI commands (serve, user management)
  user/               # User domain model + store interface
    store_sqlite/     # SQLite implementation

migrations/           # SQL migrations (numbered files)
static/               # Static assets (CSS, JS, images)
  css/                # Custom CSS overrides (app.css)
  js/                 # Frontend logic (app.js, toast.js, etc.)
```

For deeper reasoning and request flow, see **[ARCHITECTURE.md](ARCHITECTURE.md)**.

---

## Multi-Repo Structure (hagg-lib + hagg)

This project is split into two repositories:

- **hagg-lib** — reusable library (context wrapper, toast system, events, middleware, Casbin helpers)
- **hagg** — boilerplate project (imports hagg-lib)

### What hagg-lib provides:

| Package | Purpose |
|---------|---------|
| `handler` | Context wrapper with fluent API (`ctx.Render()`, `ctx.Toast()`, `ctx.Event()`) |
| `hxevents` | Server-to-client event system (HX-Trigger headers + initial-events scripts) |
| `toast` | Toast notification builder (success, error, warning, info) |
| `middleware` | BasePath injection, Security headers |
| `casbinx` | Casbin enforcer helpers |
| `view` | URL helpers (basePath-aware) |

### Why the split?

Users who fork `hagg` can update `hagg-lib` without merge conflicts.
The library provides stable building blocks, while the boilerplate is meant to be customized.

### Development (current WIP phase):

```go
// hagg/go.mod
replace github.com/axelrhd/hagg-lib => ../hagg-lib
```

### Production (v1.0.0+):

```go
// hagg/go.mod
require github.com/axelrhd/hagg-lib v1.2.3
```

---

## CSS & Styling

We use **Bootstrap 5.3** via CDN.

### Why Bootstrap?

- Battle-tested, well-documented
- No build step required (CDN delivery)
- Native dark mode support (`data-bs-theme`)
- Rich component library out of the box

### No Build Step

Bootstrap CSS and JS are loaded directly from CDN in `skeleton.go`.
Custom overrides live in `static/css/app.css`.

---

## Forking & Contributions

You are very welcome to **fork this project** and adapt it to your own needs.
`hagg` is intentionally designed as a baseline that can be copied, reshaped,
and evolved independently.

Please note that this is a **hobby and learning project**.
While we are happy to share ideas and discuss approaches,
we are not necessarily able to implement feature requests,
provide long-term support, or follow a fixed roadmap.

**The journey matters as much as the destination.**
We document our decisions, iterate openly, and learn by building real things.
If you're looking for production-ready stability, wait for v1.0.0.
If you want to learn and explore with us, welcome aboard!

---

## License

MIT
