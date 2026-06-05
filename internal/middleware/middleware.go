package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// CORS middleware — reads allowed origins from config
func CORS(allowedOrigins string) fiber.Handler {
	origins := allowedOrigins
	if origins == "" {
		origins = "*"
	}

	// Support comma-separated origins
	originList := strings.Split(origins, ",")
	for i := range originList {
		originList[i] = strings.TrimSpace(originList[i])
	}

	return cors.New(cors.Config{
		AllowOrigins:     strings.Join(originList, ","),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Authorization,Accept",
		AllowCredentials: true,
	})
}

// Logger middleware with colored output
func Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	})
}

// Recover middleware — catch panics and return 500
func Recover() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
	})
}
