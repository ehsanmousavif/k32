package main

import (
	"html/template"
	"os"

	"github.com/gofiber/fiber/v3"
)

type User struct {
	Name       string `form:"name" required:"true"`
	LastName   string `form:"lastName"  required:"true"`
	CardNumber string `form:"cardNumber"`
	Slug       string `form:"slug" required:"true"`
}

var tmpl = template.Must(template.ParseFiles("./template/template.html"))

func main() {
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	app.Post("/submit", func(c fiber.Ctx) error {
		var data User

		if err := c.Bind().Form(&data); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		if data.Slug == "" {
			errorMessage := "Slug is required"
			return c.Status(400).SendString(errorMessage)
		}

		file, err := os.Create(data.Slug + ".html")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer file.Close()

		tmpl.Execute(file, data)
		c.Next()
		c.Set("Content-Type", "text/html")
		return tmpl.Execute(c.Response().BodyWriter(), data)
	})

	app.Listen(":3000")
}
