package shared

import (
	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg/internal/session"
)

// FlashMessage represents a flash notification message.
// Flash messages are one-time notifications that appear after a redirect.
type FlashMessage struct {
	Level   string `json:"level"`   // success, error, warning, info
	Message string `json:"message"` // The message text
}

// GetFlashMessages retrieves and removes all flash messages from the session.
// This uses SCS session manager with PopString which automatically removes
// the message after reading it (flash behavior).
//
// Flash messages are automatically converted to toast notifications in the layout.
//
// Example:
//
//	messages := shared.GetFlashMessages(ctx)
//	for _, msg := range messages {
//	    ctx.Toast(msg.Message).Level(msg.Level).Notify()
//	}
func GetFlashMessages(ctx *handler.Context) []FlashMessage {
	sessionCtx := ctx.Req.Context()

	// PopString retrieves and removes the value in one operation (flash behavior)
	successMsg := session.Manager.PopString(sessionCtx, "flash_success")
	errorMsg := session.Manager.PopString(sessionCtx, "flash_error")
	warningMsg := session.Manager.PopString(sessionCtx, "flash_warning")
	infoMsg := session.Manager.PopString(sessionCtx, "flash_info")

	var messages []FlashMessage
	if successMsg != "" {
		messages = append(messages, FlashMessage{Level: "success", Message: successMsg})
	}
	if errorMsg != "" {
		messages = append(messages, FlashMessage{Level: "error", Message: errorMsg})
	}
	if warningMsg != "" {
		messages = append(messages, FlashMessage{Level: "warning", Message: warningMsg})
	}
	if infoMsg != "" {
		messages = append(messages, FlashMessage{Level: "info", Message: infoMsg})
	}

	return messages
}

// SetFlash sets a flash message in the session.
// Flash messages persist across a single redirect and are automatically removed after being read.
//
// Common pattern: Set flash before redirect
//
// Example:
//
//	shared.SetFlash(ctx, "success", "User created successfully")
//	http.Redirect(ctx.Res, ctx.Req, "/users", http.StatusSeeOther)
//
// The flash message will appear as a toast notification on the next page load.
func SetFlash(ctx *handler.Context, level, message string) {
	sessionCtx := ctx.Req.Context()
	key := "flash_" + level
	session.Manager.Put(sessionCtx, key, message)
}
