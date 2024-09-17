package adminHandlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/middlewares"
	"goHtmlBlog/models"
	"strconv"
)

func AdminRolesPage(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "5"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if limit == 0 && page == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	offset := (page - 1) * limit

	db := database.DB
	var roles []models.Role
	roleLen, _ := db.NewSelect().Model(&roles).Relation("Permissions").Order("name ASC").Limit(limit).Offset(offset).ScanAndCount(ctx)
	fmt.Println(roleLen)
	user, err := middlewares.SelectAuthenticatedUser(c, db)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	paginate := models.Pagination{}.Paginate(page, roleLen, limit)

	viewData := models.ViewData{
		Data: roles,
		Meta: models.Meta{
			PageTitle:  "Roles",
			Pagination: paginate,
			AuthedUser: user,
			BaseURL:    BaseURL,
		},
	}
	return c.Render("pages/admin/role/admin-roles", viewData)
}

func AdminRolePage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	db := database.DB
	authedUser, _ := middlewares.SelectAuthenticatedUser(c, db)
	var role models.Role
	db.NewSelect().Model(&role).Where("id = ?", id).Relation("Permissions").Scan(ctx)
	fmt.Println(role, id)
	var permissions []models.Permission
	db.NewSelect().Model(&permissions).Order("name ASC").Scan(ctx)

	viewData := models.ViewData{
		Data: struct {
			Role        models.Role
			Permissions []models.Permission
		}{
			Role:        role,
			Permissions: permissions,
		},
		Meta: models.Meta{
			PageTitle:  "Roles",
			AuthedUser: authedUser,
			BaseURL:    BaseURL,
		},
	}
	return c.Render("pages/admin/role/admin-role", viewData)
}

func AdminRoleFormPage(c *fiber.Ctx) error {
	db := database.DB
	authedUser, _ := middlewares.SelectAuthenticatedUser(c, db)

	var permissions []models.Permission
	db.NewSelect().Model(&permissions).Order("name ASC").Scan(ctx)

	viewData := models.ViewData{
		Data: permissions,
		Meta: models.Meta{
			PageTitle:  "Roles",
			AuthedUser: authedUser,
			BaseURL:    BaseURL,
		},
	}

	return c.Render("pages/admin/role/role-form", viewData)
}
