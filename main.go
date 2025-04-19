package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"

	_ "golang.org/x/oauth2"
	_ "golang.org/x/oauth2/github"
)

type TestForm struct {
	Email   string `form:"email" json:"email" binding:"required"`
	Message string `form:"message" json:"message" binding:"required"`
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

func main() {
	app := fiber.New(fiber.Config{
		Views: html.New("./templates", ".tmpl"),
	})

	// conf := &oauth2.Config{
	// 	ClientID:     "your-client-id",
	// 	ClientSecret: "your-client-secret",
	// 	RedirectURL:  "your-redirect-url",
	// 	Endpoint:     github.Endpoint,
	// }

	app.Get("/", func(ctx fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{})
	})

	app.Post("/test-form", func(ctx fiber.Ctx) error {
		var form TestForm
		if err := ctx.Bind().Body(&form); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString("Bad Request")
		}

		return ctx.Render("test-form-rsp", fiber.Map{
			"email":   form.Email,
			"message": form.Message,
		})
	})

	app.Listen(":3000")
}
