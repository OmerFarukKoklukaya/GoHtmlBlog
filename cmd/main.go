package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/router"
)

func main() {
	database.ConnectDatabase()

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	router.Router(app)
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}

}
