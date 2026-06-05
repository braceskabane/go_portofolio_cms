package middleware

import (
	"github.com/gofiber/fiber/v2"
)

const adminSessionCookie = "admin_session"

// AdminCookieName returns the session cookie name (exported for admin package)
func AdminCookieName() string { return adminSessionCookie }

// AdminSessionRequired checks for a valid admin session cookie
func AdminSessionRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies(adminSessionCookie)
		if token == "" {
			return c.Redirect("/admin/login")
		}
		// Validate the session token (JWT reuse)
		claims, err := ParseAdminToken(token)
		if err != nil {
			c.ClearCookie(adminSessionCookie)
			return c.Redirect("/admin/login")
		}
		c.Locals("adminUser", claims.Email)
		return c.Next()
	}
}
