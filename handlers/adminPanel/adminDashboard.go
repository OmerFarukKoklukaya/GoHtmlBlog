package adminHandlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/middlewares"
	"goHtmlBlog/models"
)

const BaseURL = "localhost:3000"

var ctx = context.Background()

func AdminDashboard(c *fiber.Ctx) error {
	db := database.DB
	user, _ := middlewares.SelectAuthenticatedUser(c, db)
	if middlewares.IsAuthorized(c, "", 0) {
		c.Status(fiber.StatusUnauthorized)
		return c.Redirect("/blogs")
	}
	viewData := models.ViewData{
		Data: nil,
		Meta: models.Meta{
			PageTitle:  "Dashboard",
			AuthedUser: user,
			BaseURL:    BaseURL,
		},
	}
	return c.Render("pages/admin/admin-dashboard", viewData)
}
