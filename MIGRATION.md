# Migration Guide: hagg v0.9 → v1.0

**Status:** 🚧 Work in Progress (wird während Refactoring erweitert)

Dieser Guide dokumentiert alle Breaking Changes und hilft bei der Migration einer bestehenden hagg-basierten App auf die neue Architektur.

---

## 🎯 Übersicht

**Hauptänderungen:**
- Gin → Chi (Router)
- gin-sessions → SCS (Session Management)
- gin.Context → handler.Context (Custom Context Wrapper)
- notie → toast (Notification System)
- hxevents komplett neu geschrieben (Event System)
- Tailwind CSS v4 (CSS Framework)

**Migration Strategy:** Dual-Server-Modus (Gin + Chi parallel), dann schrittweise Route-Migration

---

## 📋 Breaking Changes

### 1. Router: Gin → Chi v5

**Betrifft:** Alle Route-Definitionen, Middleware, Handler-Signaturen

#### Route Syntax

| Gin | Chi |
|-----|-----|
| `router.GET("/path", handler)` | `r.Get("/path", wrapper.Wrap(handler))` |
| `router.POST("/path", handler)` | `r.Post("/path", wrapper.Wrap(handler))` |
| `router.Group("/api")` | `r.Route("/api", func(r chi.Router) { ... })` |
| `:param` | `{param}` |

**Beispiel:**

```go
// ALT (Gin)
router := gin.New()
auth := router.Group("/auth")
auth.GET("/login", authHandler.LoginPage)
auth.POST("/login", authHandler.Login)

// NEU (Chi)
r := chi.NewRouter()
r.Route("/auth", func(r chi.Router) {
    r.Get("/login", wrapper.Wrap(authHandler.LoginPage))
    r.Post("/login", wrapper.Wrap(authHandler.Login))
})
```

#### URL Parameter Extraction

```go
// ALT (Gin)
func Handler(ctx *gin.Context) {
    id := ctx.Param("id")
    query := ctx.Query("filter")
}

// NEU (Chi)
func Handler(ctx *handler.Context) error {
    id := chi.URLParam(ctx.Req, "id")
    query := ctx.Req.URL.Query().Get("filter")
    return nil
}
```

---

### 2. Context: gin.Context → handler.Context

**Betrifft:** Alle Handler, Middleware, Helper-Funktionen

#### Handler Signature

```go
// ALT (Gin)
func MyHandler(ctx *gin.Context) {
    // ... code
}

// NEU (Chi + handler.Context)
func MyHandler(ctx *handler.Context) error {
    // ... code
    return nil  // oder error
}
```

#### Context API Changes

| Gin | handler.Context |
|-----|-----------------|
| `ctx.Param("id")` | `chi.URLParam(ctx.Req, "id")` |
| `ctx.Query("key")` | `ctx.Req.URL.Query().Get("key")` |
| `ctx.PostForm("key")` | `ctx.Req.FormValue("key")` (nach `ctx.Req.ParseForm()`) |
| `ctx.ShouldBindJSON(&data)` | `json.NewDecoder(ctx.Req.Body).Decode(&data)` |
| `ctx.HTML(200, html)` | `ctx.Render(htmlNode)` |
| `ctx.JSON(200, data)` | `json.NewEncoder(ctx.Res).Encode(data)` |
| `ctx.Status(code)` | `ctx.Res.WriteHeader(code)` |

**Beispiel:**

```go
// ALT (Gin)
func Login(ctx *gin.Context) {
    var req LoginRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // ... login logic
    ctx.JSON(200, gin.H{"success": true})
}

// NEU (handler.Context)
func Login(ctx *handler.Context) error {
    var req LoginRequest
    if err := json.NewDecoder(ctx.Req.Body).Decode(&req); err != nil {
        ctx.Res.WriteHeader(http.StatusBadRequest)
        return json.NewEncoder(ctx.Res).Encode(map[string]string{"error": err.Error()})
    }
    // ... login logic
    ctx.Res.WriteHeader(http.StatusOK)
    return json.NewEncoder(ctx.Res).Encode(map[string]bool{"success": true})
}
```

---

### 3. Sessions: gin-sessions → SCS (alexedwards/scs/v2)

**Betrifft:** Session-Zugriff, Flash Messages, Authentication

#### Session Setup

```go
// ALT (Gin)
store := cookie.NewStore([]byte(secret))
router.Use(sessions.Sessions("session_name", store))

// NEU (SCS)
session.Manager = scs.New()
session.Manager.Lifetime = 24 * time.Hour
session.Manager.Cookie.Name = "session_name"
// ... (siehe internal/session/manager.go)

// In Chi Router:
r.Use(session.Manager.LoadAndSave)
```

#### Session API

```go
// ALT (Gin)
func Handler(ctx *gin.Context) {
    session := sessions.Default(ctx)
    userID := session.Get("user_id")
    session.Set("user_id", 123)
    session.Save()
}

// NEU (SCS)
func Handler(ctx *handler.Context) error {
    sessionCtx := ctx.Req.Context()
    userID := session.Manager.GetInt(sessionCtx, "user_id")
    session.Manager.Put(sessionCtx, "user_id", 123)
    // Kein Save() nötig - automatisch beim Response
    return nil
}
```

#### Flash Messages

```go
// ALT (Gin + gin-sessions)
session := sessions.Default(ctx)
session.AddFlash("Message", "success")
session.Save()

// NEU (SCS + shared/flash.go)
shared.SetFlash(ctx, "success", "Message")
// Abrufen:
messages := shared.GetFlashMessages(ctx)
```

---

### 4. Notifications: notie → toast

**Betrifft:** Alle Notification-Aufrufe

```go
// ALT (notie)
notie.NewAlert("User created").Success().Notify(ctx)

// NEU (toast)
ctx.Toast("User created").Success().Notify()
```

**Toast API:**
- `.Success()`, `.Error()`, `.Warning()`, `.Info()`
- `.Stay()` - persistente Notification
- `.SetTimeout(ms)` - custom timeout
- `.SetPosition(pos)` - Position ("bottom-right", "top-right", etc.)

**Frontend:**
- notie.js → toast.js (neue Implementation mit surreal.js)
- Icons: SVG statt Font (siehe toast/icons.go)

---

### 5. Events: hxevents (komplett neue API)

**Betrifft:** Alle Event-Emits (HTMX Triggers)

#### Event Emission

```go
// ALT (hagg-lib/hxevents - alte Version)
// ⚠️ Alte API ist inkompatibel mit neuer Version!

// NEU (Phase 1 hxevents)
// Normale Events (für initial-events + HX-Trigger):
ctx.Event("event-name", payload)

// Events mit Phase (nur HX-Trigger):
hxevents.Add(ctx, hxevents.Immediate, "event-name", payload)
hxevents.Add(ctx, hxevents.AfterSwap, "event-name", payload)
hxevents.Add(ctx, hxevents.AfterSettle, "event-name", payload)
```

**Wichtig:** Events werden automatisch committed via `handler.Wrapper` - **kein manueller Commit** nötig!

#### Frontend Event Processing

**NEU:** Einheitliches Event-System:
- `initial-events` Script (für full-page loads)
- HX-Trigger Headers (für HTMX requests)
- Beide triggern die gleichen Event-Handler (siehe `static/js/events.js`)

---

### 6. Middleware

**Betrifft:** Custom Middleware, Auth Guards

#### Middleware Signature

```go
// ALT (Gin)
func MyMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // ... logic
        ctx.Next()
    }
}

// NEU (Chi)
func MyMiddleware(wrapper *handler.Wrapper) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // ... logic
            next.ServeHTTP(w, r)
        })
    }
}
```

#### Auth Middleware

```go
// ALT (Gin)
func RequireAuth() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        session := sessions.Default(ctx)
        if session.Get("user_id") == nil {
            ctx.Redirect(302, "/login")
            ctx.Abort()
            return
        }
        ctx.Next()
    }
}

// NEU (Chi - siehe internal/middleware/auth.go)
func RequireAuth(wrapper *handler.Wrapper) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            sessionCtx := r.Context()
            userID := session.Manager.GetInt(sessionCtx, "user_id")
            if userID == 0 {
                http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

### 7. Rendering: view.Render → handler.Context.Render

**Betrifft:** Alle Template/Gomponents Renders

```go
// ALT (hagg-lib/view)
view.Render(ctx, myPage())

// NEU (handler.Context)
return ctx.Render(myPage())
```

**Automatische Features:**
- Event commitment (HX-Trigger Headers)
- Error handling
- Content-Type Header

---

## 🔧 Dependencies

### Zu installieren:

```bash
go get github.com/go-chi/chi/v5
go get github.com/alexedwards/scs/v2
go get github.com/alexedwards/scs/sqlite3store
```

### Zu entfernen (nach Migration):

```bash
# In go.mod entfernen:
github.com/gin-gonic/gin
github.com/gin-contrib/sessions
```

---

## 📝 Migration Checklist (für echte App)

### Phase 0: Vorbereitung

- [ ] **Backup erstellen** (Code + Datenbank)
- [ ] **Dependencies auditieren:** Welche hagg-lib Features nutzt die App?
- [ ] **Route Inventory:** Liste aller Routes erstellen
- [ ] **Custom Code identifizieren:** Middleware, Helper, Extensions
- [ ] **Tests dokumentieren:** Welche Features müssen nach Migration funktionieren?

### Phase 1: Setup (Dual-Server-Modus)

- [ ] **hagg-lib aktualisieren** (neue Version mit handler, toast, hxevents)
- [ ] **Dependencies installieren** (Chi, SCS)
- [ ] **Config erweitern:**
  - `SESSION_DB_PATH` hinzufügen
  - Gin-Config beibehalten (parallel)
- [ ] **Session Manager initialisieren** (`internal/session/manager.go` kopieren)
- [ ] **Chi Server erstellen** (`internal/server/chi.go` kopieren)
- [ ] **Middleware migrieren** (`internal/middleware/chi.go` + `auth.go`)
- [ ] **Flash System anpassen** (`internal/shared/flash.go`)
- [ ] **server.go aktualisieren:** Dual-Server-Modus aktivieren
  - Gin auf :8080 (bestehend)
  - Chi auf :8081 (neu)
- [ ] **Kompilieren testen:** Projekt muss bauen
- [ ] **Beide Server starten:** Gin + Chi laufen parallel

### Phase 2: Route Migration (schrittweise)

**Strategie:** Eine Route nach der anderen

- [ ] **routing.go analysieren:** Alle Routes dokumentieren
- [ ] **Test-Route migrieren:** Einfache Route ohne Auth/Params
- [ ] **Auth Routes:** Login, Logout migrieren
- [ ] **Public Routes:** Startseite, About, etc.
- [ ] **Protected Routes:** Dashboard, Profile, etc.
- [ ] **API Routes:** REST Endpoints (falls vorhanden)

**Pro Route:**
1. Chi-Route definieren (`r.Get()`, `r.Post()`, etc.)
2. Handler anpassen:
   - Signature: `func(*handler.Context) error`
   - ctx.Param() → chi.URLParam()
   - Session API anpassen
   - Events anpassen
   - Render anpassen
3. Auf :8081 testen (Chi)
4. Regression-Check auf :8080 (Gin)
5. Commit

### Phase 3: Port Swap

- [ ] **Alle Routes migriert?** Checklist durchgehen
- [ ] **Finale Tests:** Alle Features auf :8081 testen
- [ ] **server.go anpassen:** Chi → :8080, Gin entfernen
- [ ] **Smoke Test:** App auf :8080 testen
- [ ] **Monitoring:** Errors, Performance prüfen

### Phase 4: Cleanup

- [ ] **Gin entfernen:**
  - Dependencies aus `go.mod`
  - Gin-Server Code löschen
  - Alte Middleware entfernen
- [ ] **Alte hagg-lib Packages löschen:**
  - `hagg-lib/middleware/hxtriggers.go`
  - `hagg-lib/view/render.go`
  - `hagg-lib/notie/`
- [ ] **Tests aktualisieren:**
  - Test-Helpers für Chi
  - Integration Tests
- [ ] **Dokumentation aktualisieren:**
  - README.md
  - API Docs (falls vorhanden)
- [ ] **go mod tidy:** Dependencies aufräumen

### Phase 5: Verifikation

- [ ] **Alle Features testen:** Manuelle Smoke-Tests
- [ ] **Tests laufen:** `go test ./...`
- [ ] **Deployment Test:** Staging-Umgebung testen
- [ ] **Performance Check:** Keine Regression
- [ ] **Production Deployment**

---

## ⚠️ Häufige Probleme & Lösungen

### Problem: Session geht nach Migration verloren

**Ursache:** SCS nutzt andere Cookie-Namen/Storage als gin-sessions

**Lösung:**
```go
// In session/manager.go:
Manager.Cookie.Name = "my_old_cookie_name"  // Gleicher Name wie vorher
```

Oder: Nutzer müssen sich neu einloggen (Session-Reset)

### Problem: Flash Messages verschwinden nicht

**Ursache:** SCS `PopString()` funktioniert anders als gin-sessions Flash

**Lösung:** Sicherstellen, dass `PopString()` (nicht `GetString()`) verwendet wird:
```go
// RICHTIG:
msg := session.Manager.PopString(ctx, "flash_success")

// FALSCH:
msg := session.Manager.GetString(ctx, "flash_success")  // Bleibt erhalten!
```

### Problem: URL Parameter sind leer

**Ursache:** Chi nutzt `{param}` statt `:param`

**Lösung:**
```go
// Route Definition:
r.Get("/users/{id}", ...)  // NICHT :id

// Handler:
id := chi.URLParam(ctx.Req, "id")  // NICHT ctx.Param("id")
```

### Problem: Middleware-Reihenfolge

**Ursache:** SCS middleware MUSS vor allen anderen Middlewares kommen, die Sessions nutzen

**Lösung:**
```go
r.Use(session.Manager.LoadAndSave)  // ZUERST!
r.Use(middleware.Recovery(wrapper))
r.Use(middleware.Logger(wrapper))
// ... rest
```

### Problem: Events werden nicht committed

**Ursache:** `handler.Wrapper` committed Events automatisch - manueller Commit überschreibt

**Lösung:** **Kein** manueller `hxevents.Commit()` Aufruf nötig! Der Wrapper macht das automatisch.

### Problem: Form Parsing schlägt fehl

**Ursache:** stdlib erfordert expliziten `ParseForm()` Aufruf

**Lösung:**
```go
if err := ctx.Req.ParseForm(); err != nil {
    return err
}
email := ctx.Req.FormValue("email")
```

---

## 📚 Zusätzliche Ressourcen

### Dokumentation:

- **Chi Router:** https://github.com/go-chi/chi
- **SCS Sessions:** https://github.com/alexedwards/scs
- **handler.Context:** Siehe `hagg-lib/handler/context.go`
- **Toast System:** Siehe `hagg-lib/toast/toast.go`
- **Event System:** Siehe `hagg-lib/hxevents/`

### Beispiel-Code:

- **Refactored Boilerplate:** Dieses Projekt (`hagg/`)
- **Integration Examples:** `hagg-lib/INTEGRATION.md`

### Support:

- Bei Fragen: GitHub Issues oder Diskussionen im Projekt

---

## 🔄 Änderungshistorie

**01.01.2026:** Initial version (während Phase 2 Refactoring)
- Breaking Changes dokumentiert (Gin→Chi, Sessions, Events, etc.)
- Migration Checklist erstellt
- Häufige Probleme gesammelt

*Dieses Dokument wird während des Refactorings erweitert.*

---

**🎯 Nächster Schritt:** Route Migration (Phase 2)
