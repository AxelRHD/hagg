# HAGG Boilerplate Refactoring Plan

**Goal:** Modernize the existing hagg boilerplate with final architecture decisions.

**Context:** See `ARCHITECTURE.md` for existing architecture documentation.

**Important:** HAGG is WIP and a thought experiment. Nothing is set in stone. This plan documents the current state after extensive discussion and evaluation.

---

## 🎯 Progress Tracking

**Freeze Point:** `v0.9.0` (12713e0) - Project state before refactoring

**Current Status:** 2026-01-01

### ✅ Completed

- **Phase 1, Point 3: Tailwind CSS Setup** (47f19f3)
  - Tailwind v4 standalone binary installed (~/.local/bin)
  - Pico-inspired design system with warm colors
  - Dark Mode support (manual toggle + prefers-color-scheme)
  - Toast components with SVG icons
  - justfile for build automation (css-build, css-watch)
  - Test page: `/static/test-styles.html`
  - REFACTORING_PLAN.md updated (no npm/node)

- **Documentation Update** (7365bc2)
  - README.md completely rewritten (new HAGG-Stack architecture)
  - ARCHITECTURE.md completely rewritten (Chi, surreal.js, patterns)
  - Journey/process emphasized as core of the project

- **Phase 1, Point 1: Context-Wrapper + Event System + Toast** (b5da4ad)
  - handler.Context with event storage
  - handler.Wrapper for stdlib http.HandlerFunc
  - hxevents completely new (events.go, context.go, commit.go, initial.go)
  - toast package with fluent API (replaces notie)
  - Frontend: events.js + toast.js (surreal.js)
  - Integration.md documentation

- **Phase 2 Setup: Chi Router Infrastructure** (ea9b385 + 0900875)
  - Chi v5 + SCS v2 dependencies installed
  - Session Manager (SQLite-backed, persistent)
  - Chi-compatible middleware (Logger, Recovery, CORS, Auth)
  - Chi Server setup (internal/server/chi.go)
  - Flash system for SCS (internal/shared/flash.go)
  - Dual-server mode: Gin :8080 + Chi :8081
  - Old packages temporarily compatible (for Gin routes during migration)

- **Phase 2 Complete: Gin → Chi Migration + Cleanup** (2026-01-01)
  - Login routes fully migrated (GET/, POST/, HTMX endpoints)
  - Gin server removed, Chi is sole server
  - Dual-server mode removed
  - Layout files consolidated (chi_* → standard names)
  - Auth methods cleaned up (Chi suffix removed)
  - Middleware files consolidated (internal/middleware/chi.go)
  - Gin dependencies fully removed
  - ~500 lines of code duplication eliminated
  - Flash system consolidated (SCS-based only)

- **Phase 3 Complete: Boilerplate Pages** (2026-01-01)
  - Homepage at / (HAGG stack presentation)
  - Login moved from / to /login
  - Protected dashboard at /dashboard
  - Navigation with auth-based visibility
  - Auth middleware integrated (RequireAuth)
  - Three complete page examples (public, auth, protected)

### 🔄 In Progress

- **None currently**

### 📅 Next Steps

1. **Phase 1, Point 4: Frontend JavaScript (surreal.js)**
   - Include surreal.js (in addition to Alpine.js)
   - Implement HTMX event handling
   - Toast rendering via surreal
   - Alpine.js remains for client-side state management

3. **Phase 2+**: Session Management, Template Refactoring, Testing

---

## 📋 Finalized Decisions

### 1. Context-Wrapper

**Decision:** Option C - Fields + minimal helper methods

**Design:**
```go
// hagg-lib/handler/context.go
type Context struct {
    Res http.ResponseWriter  // Explicit field (no embedding)
    Req *http.Request        // Explicit field (no embedding)

    logger *slog.Logger
    events []Event
}

// Only the helper methods that are actually needed:
func (c *Context) Render(node g.Node) error
func (c *Context) Toast(msg string) *toast.Builder
func (c *Context) Event(name string, data any)
```

**Handler Pattern:**
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

**Rationale:**
- ✅ Explicit, no magic (no embedding problems)
- ✅ Compatible with all middleware
- ✅ Only necessary helpers (no overhead)
- ✅ For boilerplate: Easy to understand
- ❌ JSON/Status helpers intentionally omitted (not needed for HTMX+Gomponents)

---

### 2. Router

**Decision:** Chi v5

**Migration:** Gin → Chi

**Rationale:**
- ✅ Closer to Go standard (stdlib-compatible)
- ✅ Minimal, no framework overhead
- ✅ Only routing + middleware (exactly what we need)
- ✅ Idiomatic Go
- ✅ Large community, actively maintained

**Example:**
```go
// Before (Gin)
r := gin.Default()
r.GET("/users", middleware.Auth(), handlers.Users)

// After (Chi)
r := chi.NewRouter()
r.Use(middleware.Logger)
r.Group(func(r chi.Router) {
    r.Use(middleware.Auth)
    r.Get("/users", deps.Wrap(handlers.Users))
})
```

---

### 3. CSS Framework

**Decision:** Tailwind CSS with Pico.css aesthetics

**Migration:** Pico.css → Tailwind CSS

**Goal:** Full flexibility of Tailwind, but recreate Pico's clean, minimalist look

**Design Principles (from Pico):**
- **Soft Rounded Corners** - Subtle border-radius (not too rounded)
- **Subtle Shadows** - Soft shadows for depth without being intrusive
- **Neutral Palette** - Gray tones as base, blue as primary
- **Good Typography** - Clear, readable fonts with enough line-height
- **Lots of Whitespace** - Airy layout, not cramped
- **Semantic Colors** - Success (green), Error (red), Warning (yellow), Info (blue)

---

### 4. Toast/Notification System

**Decision:** Initial-Events + HTMX Events with surreal.js

**Core Concept:**
- **HTMX Requests:** Events via `HX-Trigger` header (as before with hxevents)
- **Full-Page-Loads:** Events as `<script id="initial-events">` in HTML
- **Frontend:** One `showToast()` function for both cases
- **No Alpine Store** (too complex)
- **hxevents remains** (very good!)

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
    ctx.Event("toast", t)  // Uses hxevents
}
```

---

### 5. Flash/Session Management

**Decision:** Hybrid - alexedwards/scs with cookiestore default

**Core Idea:**
- **Flash:** Session-based (for redirects)
- **Events:** For everything else (HTMX, Initial-Events)
- **Flash → Event:** Flash is converted to event during layout render

**Session Store:**
- **Default:** Cookiestore (zero deps, stateless)
- **Optional:** SQLite, PostgreSQL, Redis (selectable via ENV variable)

---

### 6. Project Structure

**Decision:** Multi-Repo (hagg-lib + hagg)

**Structure:**
```
hagg-lib/                    # Reusable library
├── handler/                # Context wrapper
├── toast/                  # Toast system
├── hxevents/              # Event system (existing, refactored)
├── flash/                 # Flash helpers
├── view/                  # View helpers
├── middleware/            # Reusable middleware
└── go.mod

hagg/                       # Boilerplate project
├── lib/                   # REMOVED (becomes hagg-lib import)
├── internal/
│   ├── app/              # Deps
│   ├── auth/             # Session-based auth
│   ├── user/             # User model + store
│   └── frontend/
│       ├── layout/
│       └── pages/
├── cmd/
├── static/
├── migrations/
└── go.mod                # import github.com/axelrhd/hagg-lib
```

**Development (WIP phase, now):**
```go
// hagg/go.mod
module github.com/axelrhd/hagg

require (
    github.com/axelrhd/hagg-lib v0.0.0-dev
    // ...
)

replace github.com/axelrhd/hagg-lib => ../hagg-lib
```

---

## 🚀 Implementation Order

### Phase 1: Foundations (hagg-lib)

**Goal:** Build new core packages

1. **handler/context.go** - Context wrapper
2. **toast/toast.go** - Toast system
3. **hxevents/** - Refactoring

**Deliverable:** hagg-lib with new packages (not yet integrated)

---

### Phase 2: Integration (hagg)

**Goal:** Integrate Chi + Context wrapper

1. **Chi Router Setup**
2. **Context Wrapper Integration**
3. **Tailwind CSS Setup**
4. **Frontend JavaScript**

**Deliverable:** Functional hagg with Chi + Context + Toast (one handler)

---

### Phase 3: Migration (hagg)

**Goal:** Migrate all handlers

1. **Migrate middleware**
2. **Migrate handlers** (step by step)
3. **Adapt render helpers**

**Deliverable:** Fully migrated hagg (all handlers on new pattern)

---

### Phase 4: Cleanup & Polish

**Goal:** Remove old packages, document

1. **Cleanup**
2. **Tests**
3. **Documentation**
4. **Versioning**

**Deliverable:** Production-ready hagg v1.0.0

---

## 🎯 Success Criteria

**Phase 1 (Foundations) is successful when:**
- ✅ hagg-lib/handler/ exists with Context + Wrapper
- ✅ hagg-lib/toast/ exists with builder pattern
- ✅ hagg-lib/hxevents/ is refactored (no gin.Context)
- ✅ Tests pass

**Phase 2 (Integration) is successful when:**
- ✅ hagg uses Chi instead of Gin
- ✅ One handler works with handler.Context
- ✅ Toast system works (frontend + backend)
- ✅ Initial-events work on full-page load
- ✅ **COMPLETED** - All criteria met

**Phase 3 (Migration) is successful when:**
- ✅ All handlers use handler.Context
- ✅ All middleware are Chi-compatible
- ✅ No gin imports in hagg

**Phase 4 (Cleanup) is successful when:**
- ✅ notie/ package removed
- ✅ Tests pass
- ✅ Documentation updated
- ✅ hagg-lib v1.0.0 + hagg v1.0.0 tagged

---

## 🔄 Change History

- **2024-12-31:** Initial version based on architecture discussion
  - All 5 core decisions finalized
  - Gap analysis added
  - Implementation order defined
  - Detailed specs for all new components

- **2026-01-01:** Phase 2 completed
  - Gin → Chi migration complete
  - All duplicate code removed
  - Documentation translated to English
