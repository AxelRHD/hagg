# Migration Guide: hagg v0.9 → v1.0

**Status:** ✅ Complete (for the hagg project itself)

This guide documents all breaking changes and helps migrate an existing hagg-based app to the new architecture.

**Note:** The `hagg` boilerplate project itself has completed this migration. This document is for users who forked the old version and want to upgrade to the new Chi-based architecture.

---

## 🎯 Overview

**Main Changes:**
- Gin → Chi (Router)
- gin-sessions → SCS (Session Management)
- gin.Context → handler.Context (Custom Context Wrapper)
- notie → toast (Notification System)
- hxevents completely rewritten (Event System)
- Tailwind CSS v4 (CSS Framework)

**Migration Strategy:** Dual-server mode (Gin + Chi in parallel), then incremental route migration

---

## 📋 Breaking Changes

### 1. Router: Gin → Chi v5

**Affects:** All route definitions, middleware, handler signatures

#### Route Syntax

| Gin | Chi |
|-----|-----|
| `router.GET("/path", handler)` | `r.Get("/path", wrapper.Wrap(handler))` |
| `router.POST("/path", handler)` | `r.Post("/path", wrapper.Wrap(handler))` |
| `router.Group("/api")` | `r.Route("/api", func(r chi.Router) { ... })` |
| `:param` | `{param}` |

**Example:**

```go
// OLD (Gin)
router := gin.New()
auth := router.Group("/auth")
auth.GET("/login", authHandler.LoginPage)
auth.POST("/login", authHandler.Login)

// NEW (Chi)
r := chi.NewRouter()
r.Route("/auth", func(r chi.Router) {
    r.Get("/login", wrapper.Wrap(authHandler.LoginPage))
    r.Post("/login", wrapper.Wrap(authHandler.Login))
})
```

#### URL Parameter Extraction

```go
// OLD (Gin)
func Handler(ctx *gin.Context) {
    id := ctx.Param("id")
    query := ctx.Query("filter")
}

// NEW (Chi)
func Handler(ctx *handler.Context) error {
    id := chi.URLParam(ctx.Req, "id")
    query := ctx.Req.URL.Query().Get("filter")
    return nil
}
```

---

### 2. Context: gin.Context → handler.Context

**Affects:** All handlers, middleware, helper functions

#### Handler Signature

```go
// OLD (Gin)
func MyHandler(ctx *gin.Context) {
    // ... code
}

// NEW (Chi + handler.Context)
func MyHandler(ctx *handler.Context) error {
    // ... code
    return nil  // or error
}
```

#### Context API Changes

| Gin | handler.Context |
|-----|-----------------|
| `ctx.Param("id")` | `chi.URLParam(ctx.Req, "id")` |
| `ctx.Query("key")` | `ctx.Req.URL.Query().Get("key")` |
| `ctx.PostForm("key")` | `ctx.Req.FormValue("key")` (after `ctx.Req.ParseForm()`) |
| `ctx.ShouldBindJSON(&data)` | `json.NewDecoder(ctx.Req.Body).Decode(&data)` |
| `ctx.HTML(200, html)` | `ctx.Render(htmlNode)` |
| `ctx.JSON(200, data)` | `json.NewEncoder(ctx.Res).Encode(data)` |
| `ctx.Status(code)` | `ctx.Res.WriteHeader(code)` |

---

### 3. Sessions: gin-sessions → SCS (alexedwards/scs/v2)

**Affects:** Session access, flash messages, authentication

#### Session Setup

```go
// OLD (Gin)
store := cookie.NewStore([]byte(secret))
router.Use(sessions.Sessions("session_name", store))

// NEW (SCS)
session.Manager = scs.New()
session.Manager.Lifetime = 24 * time.Hour
session.Manager.Cookie.Name = "session_name"
// ... (see internal/session/manager.go)

// In Chi Router:
r.Use(session.Manager.LoadAndSave)
```

#### Session API

```go
// OLD (Gin)
func Handler(ctx *gin.Context) {
    session := sessions.Default(ctx)
    userID := session.Get("user_id")
    session.Set("user_id", 123)
    session.Save()
}

// NEW (SCS)
func Handler(ctx *handler.Context) error {
    sessionCtx := ctx.Req.Context()
    userID := session.Manager.GetInt(sessionCtx, "user_id")
    session.Manager.Put(sessionCtx, "user_id", 123)
    // No Save() needed - automatic on response
    return nil
}
```

---

### 4. Notifications: notie → toast

**Affects:** All notification calls

```go
// OLD (notie)
notie.NewAlert("User created").Success().Notify(ctx)

// NEW (toast)
ctx.Toast("User created").Success().Notify()
```

**Toast API:**
- `.Success()`, `.Error()`, `.Warning()`, `.Info()`
- `.Stay()` - persistent notification
- `.SetTimeout(ms)` - custom timeout
- `.SetPosition(pos)` - position ("bottom-right", "top-right", etc.)

---

### 5. Events: hxevents (completely new API)

**Affects:** All event emits (HTMX triggers)

#### Event Emission

```go
// OLD (hagg-lib/hxevents - old version)
// ⚠️ Old API is incompatible with new version!

// NEW (Phase 1 hxevents)
// Normal events (for initial-events + HX-Trigger):
ctx.Event("event-name", payload)

// Events with phase (HX-Trigger only):
hxevents.Add(ctx, hxevents.Immediate, "event-name", payload)
hxevents.Add(ctx, hxevents.AfterSwap, "event-name", payload)
hxevents.Add(ctx, hxevents.AfterSettle, "event-name", payload)
```

**Important:** Events are automatically committed via `handler.Wrapper` - **no manual commit** needed!

---

### 6. Middleware

**Affects:** Custom middleware, auth guards

#### Middleware Signature

```go
// OLD (Gin)
func MyMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // ... logic
        ctx.Next()
    }
}

// NEW (Chi)
func MyMiddleware(wrapper *handler.Wrapper) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // ... logic
            next.ServeHTTP(w, r)
        })
    }
}
```

---

### 7. Rendering: view.Render → handler.Context.Render

**Affects:** All template/gomponents renders

```go
// OLD (hagg-lib/view)
view.Render(ctx, myPage())

// NEW (handler.Context)
return ctx.Render(myPage())
```

**Automatic features:**
- Event commitment (HX-Trigger headers)
- Error handling
- Content-Type header

---

## 🔧 Dependencies

### To Install:

```bash
go get github.com/go-chi/chi/v5
go get github.com/alexedwards/scs/v2
go get github.com/alexedwards/scs/sqlite3store
```

### To Remove (after migration):

```bash
# Remove from go.mod:
github.com/gin-gonic/gin
github.com/gin-contrib/sessions
```

---

## 📝 Migration Checklist (for existing apps)

### Phase 0: Preparation

- [ ] **Create backup** (code + database)
- [ ] **Audit dependencies:** Which hagg-lib features does the app use?
- [ ] **Route inventory:** List all routes
- [ ] **Identify custom code:** Middleware, helpers, extensions
- [ ] **Document tests:** Which features must work after migration?

### Phase 1: Setup (Dual-Server Mode)

- [ ] **Update hagg-lib** (new version with handler, toast, hxevents)
- [ ] **Install dependencies** (Chi, SCS)
- [ ] **Extend config:**
  - Add `SESSION_DB_PATH`
  - Keep Gin config (parallel)
- [ ] **Initialize session manager** (copy `internal/session/manager.go`)
- [ ] **Create Chi server** (copy `internal/server/chi.go`)
- [ ] **Migrate middleware** (`internal/middleware/chi.go` + `auth.go`)
- [ ] **Adapt flash system** (`internal/shared/flash.go`)
- [ ] **Update server.go:** Enable dual-server mode
  - Gin on :8080 (existing)
  - Chi on :8081 (new)
- [ ] **Test compilation:** Project must build
- [ ] **Start both servers:** Gin + Chi run in parallel

### Phase 2: Route Migration (incremental)

**Strategy:** One route at a time

- [ ] **Analyze routing.go:** Document all routes
- [ ] **Migrate test route:** Simple route without auth/params
- [ ] **Auth routes:** Login, logout migration
- [ ] **Public routes:** Homepage, about, etc.
- [ ] **Protected routes:** Dashboard, profile, etc.
- [ ] **API routes:** REST endpoints (if any)

**Per route:**
1. Define Chi route (`r.Get()`, `r.Post()`, etc.)
2. Adapt handler:
   - Signature: `func(*handler.Context) error`
   - ctx.Param() → chi.URLParam()
   - Adapt session API
   - Adapt events
   - Adapt rendering
3. Test on :8081 (Chi)
4. Regression check on :8080 (Gin)
5. Commit

### Phase 3: Port Swap

- [ ] **All routes migrated?** Review checklist
- [ ] **Final tests:** Test all features on :8081
- [ ] **Adapt server.go:** Chi → :8080, remove Gin
- [ ] **Smoke test:** Test app on :8080
- [ ] **Monitoring:** Check errors, performance

### Phase 4: Cleanup

- [ ] **Remove Gin:**
  - Dependencies from `go.mod`
  - Gin server code
  - Old middleware
- [ ] **Delete old hagg-lib packages:**
  - `hagg-lib/middleware/hxtriggers.go`
  - `hagg-lib/view/render.go`
  - `hagg-lib/notie/`
- [ ] **Update tests:**
  - Test helpers for Chi
  - Integration tests
- [ ] **Update documentation:**
  - README.md
  - API docs (if any)
- [ ] **go mod tidy:** Clean up dependencies

### Phase 5: Verification

- [ ] **Test all features:** Manual smoke tests
- [ ] **Tests pass:** `go test ./...`
- [ ] **Deployment test:** Test staging environment
- [ ] **Performance check:** No regression
- [ ] **Production deployment**

---

## ⚠️ Common Problems & Solutions

### Problem: Session lost after migration

**Cause:** SCS uses different cookie names/storage than gin-sessions

**Solution:**
```go
// In session/manager.go:
Manager.Cookie.Name = "my_old_cookie_name"  // Same name as before
```

Or: Users must log in again (session reset)

### Problem: Flash messages don't disappear

**Cause:** SCS `PopString()` works differently than gin-sessions Flash

**Solution:** Ensure `PopString()` (not `GetString()`) is used:
```go
// CORRECT:
msg := session.Manager.PopString(ctx, "flash_success")

// WRONG:
msg := session.Manager.GetString(ctx, "flash_success")  // Persists!
```

### Problem: URL parameters are empty

**Cause:** Chi uses `{param}` instead of `:param`

**Solution:**
```go
// Route definition:
r.Get("/users/{id}", ...)  // NOT :id

// Handler:
id := chi.URLParam(ctx.Req, "id")  // NOT ctx.Param("id")
```

### Problem: Middleware order

**Cause:** SCS middleware MUST come before all other middleware that use sessions

**Solution:**
```go
r.Use(session.Manager.LoadAndSave)  // FIRST!
r.Use(middleware.Recovery(wrapper))
r.Use(middleware.Logger(wrapper))
// ... rest
```

### Problem: Events not committed

**Cause:** `handler.Wrapper` commits events automatically - manual commit overwrites

**Solution:** **No** manual `hxevents.Commit()` call needed! The wrapper does it automatically.

### Problem: Form parsing fails

**Cause:** stdlib requires explicit `ParseForm()` call

**Solution:**
```go
if err := ctx.Req.ParseForm(); err != nil {
    return err
}
email := ctx.Req.FormValue("email")
```

---

## 📚 Additional Resources

### Documentation:

- **Chi Router:** https://github.com/go-chi/chi
- **SCS Sessions:** https://github.com/alexedwards/scs
- **handler.Context:** See `hagg-lib/handler/context.go`
- **Toast System:** See `hagg-lib/toast/toast.go`
- **Event System:** See `hagg-lib/hxevents/`

### Example Code:

- **Refactored Boilerplate:** This project (`hagg/`)
- **Integration Examples:** `hagg-lib/INTEGRATION.md`

### Support:

- For questions: GitHub issues or discussions in the project

---

## 🔄 Change History

**2026-01-01:** Initial version (during Phase 2 refactoring)
- Breaking changes documented (Gin→Chi, Sessions, Events, etc.)
- Migration checklist created
- Common problems collected
- Translated to English

---

**🎯 Status:** Migration complete for the hagg project itself. This guide is for users migrating their own forks.
