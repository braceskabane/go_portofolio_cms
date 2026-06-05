package admin

import (
	"portfolio-cms/internal/config"
	"portfolio-cms/internal/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// LoginPage renders the login form
func LoginPage(c *fiber.Ctx) error {
	// Already logged in?
	if c.Cookies(middleware.AdminCookieName()) != "" {
		return c.Redirect("/admin")
	}
	flash := c.Query("error")
	flashMsg := ""
	if flash == "invalid" {
		flashMsg = "ERR:Invalid username or password"
	}
	page := loginHTML(flashMsg)
	return c.Type("html").SendString(page)
}

// LoginHandler validates credentials and sets a session cookie
func LoginHandler(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.FormValue("email")
		password := c.FormValue("password")

		// Compare against admin credentials from config
		if email != cfg.Admin.Username {
			return c.Redirect("/admin/login?error=invalid")
		}
		// Support both plain text (dev) and bcrypt (prod) admin passwords
		err := bcrypt.CompareHashAndPassword([]byte(cfg.Admin.Password), []byte(password))
		if err != nil && cfg.Admin.Password != password {
			return c.Redirect("/admin/login?error=invalid")
		}

		// Generate a JWT token for the session cookie
		token, err := middleware.GenerateAccessToken("admin", email, "admin")
		if err != nil {
			return c.Redirect("/admin/login?error=invalid")
		}

		c.Cookie(&fiber.Cookie{
			Name:     middleware.AdminCookieName(),
			Value:    token,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			HTTPOnly: true,
			SameSite: "Lax",
		})

		return c.Redirect("/admin")
	}
}

// LogoutHandler clears the session cookie
func LogoutHandler(c *fiber.Ctx) error {
	c.ClearCookie(middleware.AdminCookieName())
	return c.Redirect("/admin/login")
}

func loginHTML(flash string) string {
    flashHTML := ""
    if flash != "" {
        msg := flash
        if len(msg) > 4 && msg[:4] == "ERR:" {
            msg = msg[4:]
        }
        flashHTML = `<div style="background:#fef2f2; border-left:3px solid #dc2626; color:#991b1b; padding:10px 14px; border-radius:4px; font-size:13px; margin-bottom:16px;">` + msg + `</div>`
    }

    return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
  <title>Login — Portfolio CMS</title>
  <style>
    :root {
      --ch-950: #0f0e0f;
      --ch-900: #161516;
      --ch-800: #1e1d1e;
      --ch-700: #252425;
      --ch-600: #30353b;
      --ch-500: #4d4c4d;
      --ch-400: #6b6a6b;
      --ch-300: #9a9899;
      --ch-200: #c5c3c4;
      --ch-100: #e8e6e7;
      --ch-50: #f5f4f5;
    }
    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
    }
    body {
      font-family: 'Inter', system-ui, -apple-system, sans-serif;
      background: #f0eff0;  /* slightly tinted background */
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
      color: var(--ch-800);
    }
    .login-wrapper {
      width: 100%;
      max-width: 400px;
      padding: 20px;
    }
    .login-header {
      text-align: center;
      margin-bottom: 24px;
    }
    .login-header h1 {
      font-size: 24px;
      font-weight: 600;
      color: var(--ch-900);
    }
    .login-header p {
      font-size: 13px;
      color: var(--ch-500);
      margin-top: 4px;
    }
    .login-card {
      background: #ffffff;
      border-radius: 12px;
      box-shadow: 0 4px 24px rgba(0,0,0,0.05);
      padding: 32px 28px;
    }
    .login-card h2 {
      font-size: 18px;
      font-weight: 600;
      color: var(--ch-900);
      margin-bottom: 20px;
    }
    .form-group {
      margin-bottom: 16px;
    }
    .form-group label {
      display: block;
      font-size: 12.5px;
      font-weight: 500;
      color: var(--ch-700);
      margin-bottom: 6px;
    }
    .form-group input {
      width: 100%;
      padding: 10px 14px;
      font-size: 13px;
      border: 1px solid #d1cfd0;
      border-radius: 6px;
      color: var(--ch-800);
      background: #ffffff;
      transition: border-color 0.15s, box-shadow 0.15s;
      outline: none;
    }
    .form-group input:focus {
      border-color: var(--ch-500);
      box-shadow: 0 0 0 3px rgba(77, 76, 77, 0.15);
    }
    .btn-submit {
      width: 100%;
      padding: 11px 20px;
      background: var(--ch-800);
      color: var(--ch-50);
      border: none;
      border-radius: 6px;
      font-size: 14px;
      font-weight: 500;
      cursor: pointer;
      transition: background 0.15s;
      margin-top: 8px;
    }
    .btn-submit:hover {
      background: var(--ch-600);
    }
    .footer-text {
      text-align: center;
      margin-top: 20px;
      font-size: 11px;
      color: var(--ch-400);
    }
  </style>
</head>
<body>
  <div class="login-wrapper">
    <div class="login-header">
      <h1>Portfolio CMS</h1>
      <p>Admin Panel</p>
    </div>
    <div class="login-card">
      <h2>Sign In</h2>
      ` + flashHTML + `
      <form method="POST" action="/admin/login">
        <div class="form-group">
          <label for="email">Email / Username</label>
          <input type="text" name="email" id="email" required autofocus placeholder="Enter your email">
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <input type="password" name="password" id="password" required placeholder="Enter your password">
        </div>
        <button type="submit" class="btn-submit">Sign In</button>
      </form>
    </div>
    <p class="footer-text">Portfolio CMS v1.0 · Secure access</p>
  </div>
</body>
</html>`
}
