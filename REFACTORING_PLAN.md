# HAGG Boilerplate Refactoring Plan

**Ziel:** Das bestehende hagg-Boilerplate modernisieren mit finalen Architektur-Entscheidungen.

**Kontext:** Siehe `ARCHITECTURE.md` für bestehende Architektur-Dokumentation.

**Wichtig:** HAGG ist WIP und ein Gedankenexperiment. Nichts ist in Stein gemeißelt. Dieser Plan dokumentiert den aktuellen Stand nach ausführlicher Diskussion und Evaluation.

---

## 📋 Finalisierte Entscheidungen

### 1. Context-Wrapper

**Entscheidung:** Option C - Felder + minimale Helper-Methods

**Design:**
```go
// hagg-lib/handler/context.go
type Context struct {
    Res http.ResponseWriter  // Explizit als Feld (kein Embedding)
    Req *http.Request        // Explizit als Feld (kein Embedding)

    logger *slog.Logger
    events []Event
}

// Nur die Helper-Methods die wirklich gebraucht werden:
func (c *Context) Render(node g.Node) error
func (c *Context) Toast(msg string) *toast.Builder
func (c *Context) Event(name string, data any)
```

**Handler-Pattern:**
```go
type HandlerFunc func(*Context) error

func (d *Deps) Wrap(h HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := &Context{
            Res: w,
            Req: r,
            logger: d.Logger,
        }

        if err := h(ctx); err != nil {
            // Unified error handling
            d.Logger.Error("handler error", "error", err)
            http.Error(w, "Internal Server Error", 500)
        }
    }
}
```

**Begründung:**
- ✅ Explizit, keine Magie (keine Embedding-Probleme)
- ✅ Kompatibel mit allen Middleware
- ✅ Nur nötige Helpers (kein Overhead)
- ✅ Für Boilerplate: Einfach zu verstehen
- ❌ JSON/Status-Helpers bewusst weggelassen (nicht nötig für HTMX+Gomponents)

---

### 2. Router

**Entscheidung:** Chi v5

**Migration:** Gin → Chi

**Begründung:**
- ✅ Näher am Go-Standard (stdlib-kompatibel)
- ✅ Minimal, kein Framework-Overhead
- ✅ Nur Routing + Middleware (genau was wir brauchen)
- ✅ Idiomatisches Go
- ✅ Große Community, aktiv maintained

**Beispiel:**
```go
// Vorher (Gin)
r := gin.Default()
r.GET("/users", middleware.Auth(), handlers.Users)

// Nachher (Chi)
r := chi.NewRouter()
r.Use(middleware.Logger)
r.Group(func(r chi.Router) {
    r.Use(middleware.Auth)
    r.Get("/users", deps.Wrap(handlers.Users))
})
```

---

### 3. CSS Framework

**Entscheidung:** Tailwind CSS mit Pico.css-Ästhetik

**Migration:** Pico.css → Tailwind CSS

**Ziel:** Volle Flexibilität von Tailwind, aber Pico's cleanen, minimalistischen Look nachbauen

**Design-Prinzipien (von Pico übernehmen):**
- **Sanfte Rundungen** - Subtile border-radius (nicht zu abgerundet)
- **Dezente Schatten** - Soft shadows für Tiefe ohne Aufdringlichkeit
- **Neutral-Palette** - Grau-Töne als Basis, blau als Primary
- **Gute Typography** - Klare, lesbare Schrift mit genug Line-Height
- **Viel Whitespace** - Luftiges Layout, nicht gedrängt
- **Semantic Colors** - Success (grün), Error (rot), Warning (gelb), Info (blau)

**Tailwind Config:**
```js
// tailwind.config.js
module.exports = {
  content: [
    "./internal/frontend/**/*.go",
    "./static/**/*.js",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          500: '#3b82f6',  // Pico's blue
          600: '#2563eb',
          700: '#1d4ed8',
        },
        success: '#10b981',
        error: '#ef4444',
        warning: '#f59e0b',
        info: '#3b82f6',
      },
      borderRadius: {
        'pico': '0.375rem',  // Pico's standard radius
      },
      boxShadow: {
        'pico': '0 0 0 1px rgba(0, 0, 0, 0.05), 0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        'pico-lg': '0 0 0 1px rgba(0, 0, 0, 0.05), 0 4px 6px -1px rgba(0, 0, 0, 0.1)',
      },
      fontFamily: {
        sans: [
          'system-ui',
          '-apple-system',
          'Segoe UI',
          'Roboto',
          'sans-serif',
        ],
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
```

**Base Styles (Pico-inspiriert):**
```css
/* static/css/base.css */
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  body {
    @apply bg-gray-50 text-gray-900;
    @apply leading-relaxed;
  }

  /* Pico's form styling */
  input, textarea, select {
    @apply rounded-pico shadow-pico;
    @apply border-gray-300;
    @apply focus:border-primary-500 focus:ring-primary-500;
  }

  /* Pico's button styling */
  button, [role="button"] {
    @apply rounded-pico shadow-pico;
    @apply px-4 py-2;
    @apply font-medium;
    @apply transition-colors;
  }
}

@layer components {
  .btn-primary {
    @apply bg-primary-500 text-white;
    @apply hover:bg-primary-600;
  }

  .btn-secondary {
    @apply bg-gray-200 text-gray-900;
    @apply hover:bg-gray-300;
  }

  .card {
    @apply bg-white rounded-pico shadow-pico-lg;
    @apply p-6;
  }
}
```

**Toast-Styling (Pico-Stil):**
```css
@layer components {
  .toast {
    @apply bg-white rounded-pico shadow-pico-lg;
    @apply border-l-4;
    @apply p-4 mb-2;
    @apply min-w-[300px] max-w-md;
  }

  .toast-success {
    @apply border-success;
  }

  .toast-error {
    @apply border-error;
  }

  .toast-warning {
    @apply border-warning;
  }

  .toast-info {
    @apply border-info;
  }
}
```

**Begründung:**
- ✅ Tailwind: Volle Flexibilität für Custom-Components
- ✅ Pico-Look: Cleanes, professionelles Design
- ✅ @tailwindcss/forms: Besseres Default-Styling für Forms
- ✅ Keine externe CSS-Library (kein CDN)
- ✅ Build-Step: Purging für kleine Bundle-Size

---

### 4. Toast/Notification System

**Entscheidung:** Initial-Events + HTMX Events mit surreal.js

**Kern-Konzept:**
- **HTMX-Requests:** Events über `HX-Trigger` Header (wie bisher mit hxevents)
- **Full-Page-Loads:** Events als `<script id="initial-events">` im HTML gerendert
- **Frontend:** Eine `showToast()` Funktion für beide Fälle
- **Kein Alpine Store** (zu komplex)
- **hxevents bleibt erhalten** (sehr gut!)

**Backend (Go):**
```go
// hagg-lib/toast/toast.go
type Toast struct {
    Message  string `json:"message"`
    Level    string `json:"level"`    // success, error, warning, info
    Timeout  int    `json:"timeout"`  // ms, 0 = stay forever
    Position string `json:"position"` // bottom-right, top-right, etc.
}

func (c *Context) Toast(msg string) *Toast {
    return &Toast{
        Message:  msg,
        Level:    "info",
        Timeout:  5000,
        Position: "bottom-right",
    }
}

func (t *Toast) Success() *Toast { t.Level = "success"; return t }
func (t *Toast) Error() *Toast   { t.Level = "error"; return t }
func (t *Toast) Warning() *Toast { t.Level = "warning"; return t }
func (t *Toast) Info() *Toast    { t.Level = "info"; return t }
func (t *Toast) Stay() *Toast    { t.Timeout = 0; return t }

func (t *Toast) Notify(ctx *Context) {
    ctx.Event("toast", t)  // Nutzt hxevents
}
```

**Frontend (JavaScript mit surreal.js):**
```js
// Initial Events (Full-Page-Load)
document.addEventListener('DOMContentLoaded', () => {
    const initialEventsEl = me('#initial-events')
    if (initialEventsEl) {
        const events = JSON.parse(initialEventsEl.textContent)
        events.forEach(processEvent)
        initialEventsEl.remove()
    }
})

// HTMX Events
htmx.on('toast', (e) => processEvent({
    name: 'toast',
    payload: e.detail
}))

function processEvent(event) {
    if (event.name === 'toast') {
        showToast(event.payload)
    }
}

// Toast mit surreal.js und Destructuring
function showToast({
    message,
    level = 'info',
    timeout = 5000,
    position = 'bottom-right'
}) {
    const container = me(`#toast-container-${position}`)

    // Toast HTML erstellen
    const toastHtml = `
        <div class="toast toast-${level} ...">
            ${message}
        </div>
    `

    container.insertAdjacentHTML('beforeend', toastHtml)
    const toast = container.lastElementChild

    // Animations mit surreal.js
    me(toast).classRemove('opacity-0')
    me(toast).classAdd('opacity-100')

    if (timeout > 0) {
        setTimeout(() => me(toast).remove(), timeout)
    }
}
```

**Layout (Initial-Events Script):**
```go
func renderInitialEvents(events []Event) g.Node {
    if len(events) == 0 {
        return g.Text("")
    }

    jsonData, _ := json.Marshal(events)
    return g.El("script",
        g.Attr("type", "application/json"),
        g.Attr("id", "initial-events"),
        g.Raw(string(jsonData)),
    )
}
```

**Begründung:**
- ✅ Unified Flow: Eine `showToast()` Funktion für HTMX + Full-Page
- ✅ Nutzt bestehendes hxevents (sehr gut, kein neues System)
- ✅ Kein Alpine Store (weniger Komplexität)
- ✅ surreal.js reduziert Boilerplate
- ✅ Tailwind für Styling (volle Kontrolle)
- ❌ notie.js wird entfernt (externe Dependency)

---

### 5. Flash/Session Management

**Entscheidung:** Hybrid - alexedwards/scs mit cookiestore default

**Kern-Idee:**
- **Flash:** Session-basiert (für Redirects)
- **Events:** Für alles andere (HTMX, Initial-Events)
- **Flash → Event:** Bei Layout-Render wird Flash zu Event konvertiert

**Session-Store:**
- **Default:** Cookiestore (zero deps, stateless)
- **Optional:** SQLite, PostgreSQL, Redis (via ENV-Variable wählbar)

**Implementation:**
```go
// Setup (cmd/main.go)
func setupSessionStore(cfg config.Config, db *sql.DB) scs.Store {
    switch cfg.SessionStore {
    case "sqlite":
        return sqlite3store.New(db)
    case "postgres":
        return postgresstore.New(db)
    case "redis":
        return redisstore.New(pool)
    case "memory":
        return memstore.New()
    default:
        // Default: Cookiestore (zero deps)
        log.Println("Using cookie session store (default)")
        return cookiestore.New([]byte(cfg.SecretKey))
    }
}

sessionManager := scs.New()
sessionManager.Store = setupSessionStore(cfg, db)
sessionManager.Lifetime = 24 * time.Hour
```

**Flash-Nutzung:**
```go
// Handler: Flash setzen (bei Redirect)
func LogoutHandler(ctx *Context, deps *Deps) error {
    sessionManager.Put(ctx.Req.Context(), "flash", Flash{
        Message: "Logged out successfully",
        Level:   "info",
    })

    http.Redirect(ctx.Res, ctx.Req, "/", http.StatusFound)
    return nil
}

// Layout: Flash → Event konvertieren
func renderLayout(ctx *Context, deps *Deps) g.Node {
    if flash := sessionManager.PopString(ctx.Req.Context(), "flash"); flash != "" {
        ctx.Event("toast", parseFlash(flash))
    }

    return html.HTML5(
        html.Body(
            renderInitialEvents(ctx.events),
            // ... content
        ),
    )
}
```

**Begründung:**
- ✅ Flash nur für Redirects (etabliertes Pattern)
- ✅ Events für HTMX/Initial-Events (flexibler)
- ✅ Cookiestore als Default (zero deps, stateless, cloud-ready)
- ✅ Store austauschbar ohne Code-Änderung (nur ENV)
- ✅ alexedwards/scs: stdlib-kompatibel, aktiv maintained

---

### 6. Projekt-Struktur

**Entscheidung:** Multi-Repo (hagg-lib + hagg)

**Struktur:**
```
hagg-lib/                    # Wiederverwendbare Library
├── handler/                # Context-Wrapper
├── toast/                  # Toast-System
├── hxevents/              # Event-System (bestehend, refactored)
├── flash/                 # Flash-Helpers
├── view/                  # View-Helpers
├── middleware/            # Wiederverwendbare Middleware
└── go.mod

hagg/                       # Boilerplate-Projekt
├── lib/                   # REMOVED (wird zu hagg-lib import)
├── internal/
│   ├── app/              # Deps
│   ├── auth/             # Session-based Auth
│   ├── user/             # User-Model + Store
│   └── frontend/
│       ├── layout/
│       └── pages/
├── cmd/
├── static/
├── migrations/
└── go.mod                # import github.com/axelrhd/hagg-lib
```

**Development (WIP-Phase, jetzt):**
```go
// hagg/go.mod
module github.com/axelrhd/hagg

require (
    github.com/axelrhd/hagg-lib v0.0.0-dev
    // ...
)

replace github.com/axelrhd/hagg-lib => ../hagg-lib
```

**Production (v1.0.0, später):**
```go
// hagg/go.mod
module github.com/axelrhd/hagg

require (
    github.com/axelrhd/hagg-lib v1.2.3  // echte Version
    // ...
)

// replace entfernt!
```

**User-Projekte (geforkte Boilerplates):**
```bash
# User forkt hagg → my-shop
cd my-shop
go get github.com/axelrhd/hagg-lib@latest  # Update ohne Merge-Konflikt!
```

**Begründung:**
- ✅ User können hagg-lib updaten ohne Merge-Konflikte
- ✅ Klare Trennung: generic (lib) vs. app-spezifisch (hagg)
- ✅ Während WIP: replace directive für schnelles Iterieren
- ✅ Später: Echte Versionen, SemVer
- ✅ hagg-lib unabhängig nutzbar (Nebeneffekt, kein Primärziel)

---

## 🔍 Gap-Analyse: Was ändert sich?

### ✅ Bleibt (gut, übernehmen)

**hagg-lib:**
- ✅ **hxevents** - Event-System (SEHR GUT!)
  - Phasen-Abstraktion (Immediate, AfterSwap, AfterSettle)
  - Commit-Pattern
  - Muss nur refactored werden: gin.Context → http.ResponseWriter

- ✅ **view** - View-Helpers
  - URL-Helpers, etc.
  - Evtl. erweitern

**hagg:**
- ✅ **internal/auth** - Session-based Auth
  - Funktioniert gut
  - Nur Session-Library wechseln (gin-sessions → scs)

- ✅ **internal/user** - User-Model + Store
  - Bleibt wie es ist

- ✅ **internal/frontend/layout** - Layout-Components
  - Bleibt, nur Events-Rendering hinzufügen

- ✅ **Casbin** - Authorization
  - Funktioniert gut
  - Evtl. Middleware anpassen für Chi

### ⚠️ Migrieren (anpassen)

**hagg-lib:**
- ⚠️ **notie** → **toast** (neu bauen)
  - Builder-Pattern übernehmen
  - notie.js entfernen
  - Eigene Tailwind-Lösung

- ⚠️ **flash** → **flash** (refactoren)
  - gin-contrib/sessions → alexedwards/scs
  - API bleibt ähnlich

- ⚠️ **hxevents** (refactoren)
  - gin.Context → http.ResponseWriter/eigener State
  - API bleibt gleich

**hagg:**
- ⚠️ **Pico.css** → **Tailwind CSS** (komplett ersetzen)
  - Pico.css entfernen
  - Tailwind installieren & konfigurieren
  - Base-Styles schreiben (Pico-inspiriert)
  - Alle Components auf Tailwind umstellen

- ⚠️ **routing.go** - Gin → Chi
  - Route-Definitionen umschreiben
  - Middleware anpassen

- ⚠️ **server.go** - Server-Setup
  - Gin → Chi
  - Session-Manager integrieren

- ⚠️ **internal/http/render** - Render-Helpers
  - gin.Context → handler.Context

- ⚠️ **internal/http/middleware** - Middleware
  - gin.HandlerFunc → http.HandlerFunc (Chi-kompatibel)
  - Auth, Permissions, etc. anpassen

- ⚠️ **internal/frontend/pages/*** - Handler
  - gin.Context → handler.Context
  - notie → ctx.Toast()
  - Schrittweise migrieren

### ✨ Neu bauen

**hagg-lib:**
- ✨ **handler/** - Context-Wrapper (komplett neu)
  - Context-Type
  - Wrapper-Funktion
  - Helper-Methods

- ✨ **toast/** - Toast-System (komplett neu)
  - Toast-Type mit Builder
  - JSON-Serialisierung
  - Event-Integration

**hagg:**
- ✨ **static/js/app.js** - Frontend-Logic
  - Initial-Events Handler
  - showToast() Funktion mit surreal.js
  - HTMX Event-Listener

- ✨ **Tailwind Setup** - CSS
  - tailwind.config.js
  - Toast-Styling
  - Pico-inspiriertes Design

---

## 🚀 Implementierungs-Reihenfolge

### Phase 1: Foundations (hagg-lib)

**Ziel:** Neue Core-Packages bauen

1. **handler/context.go** - Context-Wrapper
   - Context-Type
   - Basic Render() Helper

2. **toast/toast.go** - Toast-System
   - Toast-Type
   - Builder-Methods
   - JSON-Tags

3. **hxevents/** - Refactoring
   - gin.Context entfernen
   - Auf http.ResponseWriter umstellen
   - Tests anpassen

**Deliverable:** hagg-lib mit neuen Packages (noch nicht integriert)

---

### Phase 2: Integration (hagg)

**Ziel:** Chi + Context-Wrapper integrieren

1. **Chi Router Setup**
   - routing.go umschreiben
   - server.go anpassen
   - Session-Manager integrieren (scs)

2. **Context-Wrapper Integration**
   - Deps.Wrap() implementieren
   - Ein Test-Handler migrieren (z.B. Login)
   - Testen

3. **Tailwind CSS Setup**
   - Tailwind CLI binary nutzen (NO npm/node!)
   - tailwind.config.js erstellen (Pico-inspirierte Config)
   - static/css/base.css erstellen (Base-Styles + Components)
   - justfile für CSS-Build (just css-watch / just css-build)
   - Pico.css entfernen

4. **Frontend JavaScript**
   - static/js/app.js erstellen
   - Initial-Events Handler
   - Toast-Handler mit surreal.js
   - Event-Processor

**Deliverable:** Funktionierendes hagg mit Chi + Context + Toast (ein Handler)

---

### Phase 3: Migration (hagg)

**Ziel:** Alle Handler migrieren

1. **Middleware migrieren**
   - Auth-Middleware (Chi-kompatibel)
   - Permissions-Middleware
   - Logger, Recovery, etc.

2. **Handler migrieren** (Schritt für Schritt)
   - internal/frontend/pages/login → handler.Context
   - internal/frontend/pages/dashboard → handler.Context
   - Jeweils testen

3. **Render-Helpers anpassen**
   - internal/http/render/page.go
   - gin.Context → handler.Context

**Deliverable:** Komplett migriertes hagg (alle Handler auf neuem Pattern)

---

### Phase 4: Cleanup & Polish

**Ziel:** Alte Packages entfernen, dokumentieren

1. **Cleanup**
   - notie/ Package entfernen (hagg-lib)
   - Gin-Dependencies entfernen (hagg + hagg-lib)
   - Alte flash/-Implementation entfernen

2. **Tests**
   - Unit-Tests für handler/
   - Unit-Tests für toast/
   - Integration-Tests

3. **Dokumentation**
   - README.md updaten
   - ARCHITECTURE.md updaten
   - Code-Kommentare

4. **Versioning**
   - hagg-lib v1.0.0 Tag
   - hagg v1.0.0 Tag
   - replace directive entfernen

**Deliverable:** Production-ready hagg v1.0.0

---

## 📝 Detaillierte Specs

### Context-Wrapper Spec

**File:** `hagg-lib/handler/context.go`

```go
package handler

import (
    "log/slog"
    "net/http"

    g "maragu.dev/gomponents"
)

type Context struct {
    Res http.ResponseWriter
    Req *http.Request

    logger *slog.Logger
    events []Event
}

type Event struct {
    Name    string      `json:"name"`
    Payload interface{} `json:"payload"`
}

func (c *Context) Render(node g.Node) error {
    c.Res.Header().Set("Content-Type", "text/html; charset=utf-8")
    return node.Render(c.Res)
}

func (c *Context) Event(name string, payload interface{}) {
    c.events = append(c.events, Event{
        Name:    name,
        Payload: payload,
    })
}

func (c *Context) Toast(msg string) *toast.Builder {
    return toast.New(msg, c)
}

func (c *Context) Logger() *slog.Logger {
    return c.logger
}

func (c *Context) Events() []Event {
    return c.events
}
```

**File:** `hagg-lib/handler/wrapper.go`

```go
package handler

import (
    "log/slog"
    "net/http"
)

type HandlerFunc func(*Context) error

type Wrapper struct {
    logger *slog.Logger
}

func NewWrapper(logger *slog.Logger) *Wrapper {
    return &Wrapper{logger: logger}
}

func (w *Wrapper) Wrap(h HandlerFunc) http.HandlerFunc {
    return func(res http.ResponseWriter, req *http.Request) {
        ctx := &Context{
            Res:    res,
            Req:    req,
            logger: w.logger,
            events: make([]Event, 0),
        }

        if err := h(ctx); err != nil {
            w.logger.Error("handler error",
                "path", req.URL.Path,
                "error", err,
            )
            http.Error(res, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        // Commit events (via HX-Trigger header if HTMX request)
        // This will be done in hxevents integration
    }
}
```

---

### Toast System Spec

**File:** `hagg-lib/toast/toast.go`

```go
package toast

import "github.com/axelrhd/hagg-lib/handler"

type Toast struct {
    Message  string `json:"message"`
    Level    string `json:"level"`
    Timeout  int    `json:"timeout"`
    Position string `json:"position"`

    ctx *handler.Context
}

func New(msg string, ctx *handler.Context) *Toast {
    return &Toast{
        Message:  msg,
        Level:    "info",
        Timeout:  5000,
        Position: "bottom-right",
        ctx:      ctx,
    }
}

func (t *Toast) Success() *Toast {
    t.Level = "success"
    return t
}

func (t *Toast) Error() *Toast {
    t.Level = "error"
    return t
}

func (t *Toast) Warning() *Toast {
    t.Level = "warning"
    return t
}

func (t *Toast) Info() *Toast {
    t.Level = "info"
    return t
}

func (t *Toast) Stay() *Toast {
    t.Timeout = 0
    return t
}

func (t *Toast) SetTimeout(ms int) *Toast {
    t.Timeout = ms
    return t
}

func (t *Toast) SetPosition(pos string) *Toast {
    t.Position = pos
    return t
}

func (t *Toast) Notify() {
    t.ctx.Event("toast", t)
}
```

---

### Session Store Setup Spec

**File:** `hagg/internal/session/store.go`

```go
package session

import (
    "database/sql"
    "log"

    "github.com/alexedwards/scs/v2"
    "github.com/alexedwards/scs/cookiestore"
    "github.com/alexedwards/scs/memstore"
    "github.com/alexedwards/scs/sqlite3store"

    "github.com/axelrhd/hagg/internal/config"
)

func Setup(cfg config.Config, db *sql.DB) *scs.SessionManager {
    manager := scs.New()
    manager.Store = createStore(cfg, db)
    manager.Lifetime = cfg.SessionLifetime
    manager.Cookie.Name = "session_id"
    manager.Cookie.HttpOnly = true
    manager.Cookie.SameSite = http.SameSiteLaxMode
    manager.Cookie.Secure = cfg.Production

    return manager
}

func createStore(cfg config.Config, db *sql.DB) scs.Store {
    switch cfg.SessionStore {
    case "sqlite":
        if db == nil {
            log.Fatal("sqlite store requires database")
        }
        log.Println("Using SQLite session store")
        return sqlite3store.New(db)

    case "memory":
        log.Println("Using memory session store (development only)")
        return memstore.New()

    default:
        log.Println("Using cookie session store (default)")
        return cookiestore.New([]byte(cfg.SecretKey))
    }
}
```

---

### Tailwind CSS Setup Spec

**Installation:** Tailwind CLI (standalone binary, NO node/npm!)

```bash
# Install to ~/.local/bin (einmalig)
curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 \
  -o ~/.local/bin/tailwindcss
chmod +x ~/.local/bin/tailwindcss
```

**justfile:**
```justfile
# justfile
css-build:
	tailwindcss -i ./static/css/base.css -o ./static/css/styles.css --minify

css-watch:
	tailwindcss -i ./static/css/base.css -o ./static/css/styles.css --watch
```

**File:** `hagg/tailwind.config.js`

```js
module.exports = {
  content: [
    "./internal/frontend/**/*.go",
    "./static/**/*.js",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',  // Pico blue
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
        },
        success: '#10b981',
        error: '#ef4444',
        warning: '#f59e0b',
        info: '#3b82f6',
      },
      borderRadius: {
        'pico': '0.375rem',
      },
      boxShadow: {
        'pico': '0 0 0 1px rgba(0, 0, 0, 0.05), 0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        'pico-lg': '0 0 0 1px rgba(0, 0, 0, 0.05), 0 4px 6px -1px rgba(0, 0, 0, 0.1)',
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
```

**File:** `hagg/static/css/base.css`

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  body {
    @apply bg-gray-50 text-gray-900;
    @apply leading-relaxed;
    @apply font-sans;
  }

  /* Pico-inspired form styling */
  input:not([type="checkbox"]):not([type="radio"]),
  textarea,
  select {
    @apply rounded-pico shadow-pico;
    @apply border-gray-300;
    @apply px-3 py-2;
    @apply focus:border-primary-500 focus:ring-1 focus:ring-primary-500;
    @apply transition-colors;
  }

  /* Pico-inspired button styling */
  button,
  [role="button"] {
    @apply rounded-pico shadow-pico;
    @apply px-4 py-2;
    @apply font-medium;
    @apply transition-colors;
    @apply cursor-pointer;
  }

  /* Headings */
  h1, h2, h3, h4, h5, h6 {
    @apply font-semibold;
    @apply mb-4;
  }

  h1 { @apply text-3xl; }
  h2 { @apply text-2xl; }
  h3 { @apply text-xl; }
}

@layer components {
  /* Buttons */
  .btn-primary {
    @apply bg-primary-500 text-white;
    @apply hover:bg-primary-600;
    @apply focus:ring-2 focus:ring-primary-500 focus:ring-offset-2;
  }

  .btn-secondary {
    @apply bg-gray-200 text-gray-900;
    @apply hover:bg-gray-300;
    @apply focus:ring-2 focus:ring-gray-400 focus:ring-offset-2;
  }

  .btn-danger {
    @apply bg-error text-white;
    @apply hover:bg-red-600;
  }

  /* Card */
  .card {
    @apply bg-white rounded-pico shadow-pico-lg;
    @apply p-6;
  }

  /* Container */
  .container-narrow {
    @apply max-w-4xl mx-auto px-4;
  }

  /* Toasts */
  .toast {
    @apply bg-white rounded-pico shadow-pico-lg;
    @apply border-l-4;
    @apply p-4 mb-2;
    @apply min-w-[300px] max-w-md;
    @apply transition-opacity duration-300;
  }

  .toast-success {
    @apply border-success;
  }

  .toast-error {
    @apply border-error;
  }

  .toast-warning {
    @apply border-warning;
  }

  .toast-info {
    @apply border-info;
  }
}
```

**Build-Prozess:**
```bash
# Development (watch mode)
tailwindcss -i ./static/css/base.css -o ./static/css/styles.css --watch

# Production (minified)
tailwindcss -i ./static/css/base.css -o ./static/css/styles.css --minify

# Oder via justfile
just css-watch
just css-build
```

**Im Layout referenzieren:**
```go
// internal/frontend/layout/base.go
func Head() g.Node {
    return html.Head(
        html.TitleEl(g.Text("HAGG App")),
        g.El("link",
            g.Attr("rel", "stylesheet"),
            g.Attr("href", "/static/css/styles.css"),
        ),
        // ... other tags
    )
}
```

---

### Frontend Spec

**File:** `hagg/static/js/app.js`

```js
// Initial Events (Full-Page-Load)
document.addEventListener('DOMContentLoaded', () => {
    const initialEventsEl = me('#initial-events')
    if (initialEventsEl) {
        const events = JSON.parse(initialEventsEl.textContent)
        events.forEach(processEvent)
        initialEventsEl.remove()
    }
})

// HTMX Events
htmx.on('toast', (e) => {
    processEvent({
        name: 'toast',
        payload: e.detail
    })
})

// Central Event Processor
function processEvent(event) {
    switch(event.name) {
        case 'toast':
            showToast(event.payload)
            break
        case 'auth-changed':
            htmx.trigger(me('#nav'), 'refresh')
            break
    }
}

// Toast Function (with surreal.js)
function showToast({
    message,
    level = 'info',
    timeout = 5000,
    position = 'bottom-right'
}) {
    const container = me(`#toast-container-${position}`)
    if (!container) return

    const toastHtml = `
        <div class="toast toast-${level} opacity-0 transition-opacity duration-300">
            <div class="flex items-center gap-3">
                <div class="flex-shrink-0">${getIcon(level)}</div>
                <div class="flex-1">${message}</div>
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

function getIcon(level) {
    const icons = {
        success: '<svg>...</svg>',
        error: '<svg>...</svg>',
        warning: '<svg>...</svg>',
        info: '<svg>...</svg>'
    }
    return icons[level] || icons.info
}
```

---

## 🎯 Success Criteria

**Phase 1 (Foundations) ist erfolgreich wenn:**
- ✅ hagg-lib/handler/ existiert mit Context + Wrapper
- ✅ hagg-lib/toast/ existiert mit Builder-Pattern
- ✅ hagg-lib/hxevents/ ist refactored (kein gin.Context)
- ✅ Tests laufen durch

**Phase 2 (Integration) ist erfolgreich wenn:**
- ✅ hagg nutzt Chi statt Gin
- ✅ Ein Handler funktioniert mit handler.Context
- ✅ Toast-System funktioniert (Frontend + Backend)
- ✅ Initial-Events funktionieren bei Full-Page-Load

**Phase 3 (Migration) ist erfolgreich wenn:**
- ✅ Alle Handler nutzen handler.Context
- ✅ Alle Middleware sind Chi-kompatibel
- ✅ Keine gin-Imports mehr in hagg

**Phase 4 (Cleanup) ist erfolgreich wenn:**
- ✅ notie/ Package entfernt
- ✅ Tests laufen durch
- ✅ Dokumentation aktualisiert
- ✅ hagg-lib v1.0.0 + hagg v1.0.0 getaggt

---

## 📚 Ressourcen

- [Chi Router](https://github.com/go-chi/chi)
- [alexedwards/scs](https://github.com/alexedwards/scs)
- [surreal.js](https://github.com/gnat/surreal)
- [Gomponents](https://www.gomponents.com/)
- [HTMX](https://htmx.org/)

---

## 🔄 Änderungshistorie

- **2024-12-31:** Initiale Version basierend auf Architektur-Diskussion
  - Alle 5 Kern-Entscheidungen finalisiert
  - Gap-Analyse hinzugefügt
  - Implementierungs-Reihenfolge definiert
  - Detaillierte Specs für alle neuen Components
