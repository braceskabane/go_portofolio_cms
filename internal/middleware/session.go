package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const sessionCookieName = "session_token"

// SessionCookieName returns the cookie name (exported for other packages)
func SessionCookieName() string { return sessionCookieName }

// SetSessionCookie menyimpan token ke cookie — reusable untuk konteks apapun
func SetSessionCookie(c *fiber.Ctx, token string, duration time.Duration) {
	c.Cookie(&fiber.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Expires:  time.Now().Add(duration),
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})
}

// ClearSessionCookie menghapus cookie dengan atribut identik saat di-set
func ClearSessionCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})
}

// SessionRequired adalah middleware generik yang memvalidasi session cookie.
// redirectTo: path tujuan jika session tidak valid, misal "/admin/login" atau "/login"
func SessionRequired(redirectTo string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies(sessionCookieName)
		if token == "" {
			return c.Redirect(redirectTo)
		}
		subject, err := ParseRefreshToken(token)
		if err != nil {
			ClearSessionCookie(c)
			return c.Redirect(redirectTo)
		}
		c.Locals("sessionUser", subject) // subject = email atau userID
		return c.Next()
	}
}