package adminHandlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/middlewares"
	"goHtmlBlog/models"
)

func AdminPermissionsPage(c *fiber.Ctx) error {
	db := database.DB
	authedUser, _ := middlewares.SelectAuthenticatedUser(c, db)
	var permissions []models.Permission
	err := db.NewSelect().Model(&permissions).Order("name ASC").Scan(ctx)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	viewData := models.ViewData{
		Data: permissions,
		Meta: models.Meta{
			PageTitle:  "Permissions",
			AuthedUser: authedUser,
			BaseURL:    BaseURL,
		},
	}

	return c.Render("pages/admin/permission/admin-permissions", viewData)
}

func AdminPermissionPage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	db := database.DB
	authedUser, _ := middlewares.SelectAuthenticatedUser(c, db)
	var permission models.Permission
	db.NewSelect().Model(&permission).Where("id = ?", id).Scan(ctx)
	viewData := models.ViewData{
		Data: permission,
		Meta: models.Meta{
			PageTitle:  "Permissions",
			AuthedUser: authedUser,
			BaseURL:    BaseURL,
		},
	}

	return c.Render("pages/admin/permission/admin-permission", viewData)
}

func AdminPermissionsFormPage(c *fiber.Ctx) error {
	db := database.DB
	authedUser, _ := middlewares.SelectAuthenticatedUser(c, db)
	viewData := models.ViewData{
		Meta: models.Meta{
			PageTitle:  "Permissions",
			AuthedUser: authedUser,
			BaseURL:    BaseURL,
		},
	}
	return c.Render("pages/admin/permission/permission-form", viewData)
}
