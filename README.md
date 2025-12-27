# hagg

`hagg` is a deliberately minimal baseline for **server-side rendered Go web applications**.
It is designed for developers who want to build long-lived, maintainable tools — without SPA overhead,
without build pipelines, and without hidden magic.

This project is **not a framework**.
It is a **stable baseline** meant to be copied, adapted, and used for years.

Following the tradition of established stack acronyms such as JAMstack, MERN, or LAMP,
we refer to the stack used by `hagg` as the **HAGG Stack**
(**H**TMX · **A**lpine.js · **G**omponents · **G**in/Go).

> **Project status:** This repository is still early **WIP (work in progress)**.
> We are building it as a hobby project to learn by doing — expect rough edges and changes.

---

## Project Goals

- **Server-side HTML** as the default (simple, fast, inspectable)
- **Progressive enhancement** with HTMX (small, explicit interactions)
- **No “frontend build” requirement** (you can ship a single Go binary + static assets)
- **Clear boundaries** (config, auth, storage, middleware, pages)
- **Replaceable parts** (swap persistence, session store, auth/authorization model)
- **Pragmatic defaults** over academic perfection

---

## The Stack in One Sentence

Gin routes requests → Gomponents render HTML → HTMX updates parts of the page → Alpine adds tiny client-side state.

---

## Quickstart

### Requirements

- Go (see `go.mod`)
- SQLite (default database)
- Optional: `task` (Taskfile) and `air` (hot reload)

### Run (dev)

```bash
task dev
```

This uses `air` and sets `GIN_DEV=true` for development mode (see `Taskfile.yml`).

### Run (prod-ish)

```bash
task build
./bin/gin-gomp
```

---

## Configuration

Configuration is loaded from environment variables (and an optional `.env`).

- Server config is prefixed with `GIN_` (e.g. `GIN_PORT`, `GIN_BASE_PATH`, `GIN_DEV`)
- Session config is prefixed with `SESSION_`
- Database config is prefixed with `DB_` (see `internal/config` for details)

To print the active configuration:

```bash
go run ./cmd -config
```

---

## Authentication

Authentication is intentionally simple:

- a session cookie stores the logged-in user's UID (`internal/auth`, session key `uid`)
- pages / HTMX endpoints use that UID to load the current user from the store

See:

- `internal/auth/auth.go`
- `internal/http/middleware/require_auth.go`
- `internal/frontend/pages/login/*`

---

## Authorization (Permission Layer)

`hagg` includes a permission layer built on **Casbin**.

### Why Casbin?

We do **not** want permissions hard-coded into handlers.
Changing roles or permissions should not require a new release.

Casbin gives us:

- RBAC-style policies (roles, users, inheritance)
- a clean enforcement API (`Enforce(subject, action)`)
- a simple file-based adapter for local projects

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

`internal/authz.MustNewEnforcer()` loads `model.conf` and `policy.csv` from the working directory
and injects a `*casbin.Enforcer` into `app.Deps`.

### Middleware

`internal/http/middleware/require_permission.go` provides:

```go
middleware.RequirePermission(deps, "user:list")
```

Behavior:

- if authentication is missing → flash + redirect (same as `RequireAuth`)
- if the policy/model are broken → 500 + notie error
- if permission is denied → 403 + notie warning

> Tip: In routes, we typically compose middlewares like:
> `RequireAuth` (or `RequirePermission`) → handler.

---

## Project Layout

A simplified overview:

```
cmd/                # composition root (CLI flags + server startup)
internal/
  app/              # dependency container (Deps)
  auth/             # session auth
  authz/            # Casbin enforcer setup
  config/           # env/.env config loading
  http/
    middleware/     # Gin middleware (auth, permissions, base path, HX triggers)
    render/         # DRY page renderer (gomponents)
  frontend/         # layouts + pages (gomponents)
  user/             # domain model + store interface
  user/store_sqlite # SQLite implementation
migrations/         # SQL migrations
static/             # static assets (served, also embedded)
model.conf
policy.csv
```

For deeper reasoning and request flow, see **[ARCHITECTURE.md](ARCHITECTURE.md)**.

---

## Forking & Contributions

You are very welcome to **fork this project** and adapt it to your own needs.
`hagg` is intentionally designed as a baseline that can be copied, reshaped,
and evolved independently.

Please note that this is a **hobby and learning project**.
While we are happy to share ideas and discuss approaches,
we are not necessarily able to implement feature requests,
provide long-term support, or follow a fixed roadmap.
