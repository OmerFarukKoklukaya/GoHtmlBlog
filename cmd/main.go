package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/router"
	"goHtmlBlog/utils"
)

func main() {
	database.ConnectDatabase()
	db := database.DB
	err := utils.AddRoles(db)
	if err != nil {
		panic(err)
	}

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	router.Router(app)
	app.Listen(":3000")

}
