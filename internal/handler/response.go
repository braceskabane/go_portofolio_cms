package handler

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func OK(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Success"
	if len(message) > 0 { msg = message[0] }
	return c.Status(fiber.StatusOK).JSON(Response{Success: true, Message: msg, Data: data})
}

func OKWithMeta(c *fiber.Ctx, data interface{}, meta interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{Success: true, Data: data, Meta: meta})
}

func Created(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Created successfully"
	if len(message) > 0 { msg = message[0] }
	return c.Status(fiber.StatusCreated).JSON(Response{Success: true, Message: msg, Data: data})
}

func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{Success: false, Message: message})
}

func Unauthorized(c *fiber.Ctx, message ...string) error {
	msg := "Unauthorized"
	if len(message) > 0 { msg = message[0] }
	return c.Status(fiber.StatusUnauthorized).JSON(Response{Success: false, Message: msg})
}

func NotFound(c *fiber.Ctx, message ...string) error {
	msg := "Resource not found"
	if len(message) > 0 { msg = message[0] }
	return c.Status(fiber.StatusNotFound).JSON(Response{Success: false, Message: msg})
}

func InternalError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response{Success: false, Message: "Internal server error: " + message})
}
