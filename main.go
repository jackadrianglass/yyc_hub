package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
)

type TestForm struct {
	Email   string `form:"email" json:"email" binding:"required"`
	Message string `form:"message" json:"message" binding:"required"`
}

func main() {
	app := fiber.New(fiber.Config{
		Views: html.New("./templates", ".tmpl"),
	})

	app.Get("/", func(ctx fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{})
	})

	app.Post("/test-form", func(ctx fiber.Ctx) error {
		var form TestForm
		if err := ctx.Bind().Body(&form); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString("Bad Request")
		}

		return ctx.Render("test-form-rsp", fiber.Map{
			"email": form.Email,
			"message": form.Message,
		})
	})

	app.Listen(":3000")
}
