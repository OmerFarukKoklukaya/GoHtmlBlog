package adminHandlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/middlewares"
	"goHtmlBlog/models"
	"strconv"
)

func AdminUsersPage(c *fiber.Ctx) error {
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
	var users []models.User
	userLen, _ := db.NewSelect().Model(&users).Relation("Role").Order("name ASC").Limit(limit).Offset(offset).ScanAndCount(ctx)
	fmt.Println(userLen)
	user, err := middlewares.SelectAuthenticatedUser(c, db)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	paginate := models.Pagination{}.Paginate(page, userLen, limit)

	viewData := models.ViewData{
		Data: users,
		Meta: models.Meta{
			PageTitle:  "Users",
			Pagination: paginate,
			AuthedUser: user,
			BaseURL:    BaseURL,
		},
	}
	return c.Render("pages/admin/user/admin-users", viewData)
}

func AdminUserPage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	db := database.DB
	authedUser, _ := middlewares.SelectAuthenticatedUser(c, db)
	var user models.User
	db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	var roles []models.Role
	db.NewSelect().Model(&roles).Scan(ctx)

	viewData := models.ViewData{
		Data: struct {
			User  models.User
			Roles []models.Role
		}{
			User:  user,
			Roles: roles,
		},
		Meta: models.Meta{
			PageTitle:  "Users",
			AuthedUser: authedUser,
			BaseURL:    BaseURL,
		},
	}
	return c.Render("pages/admin/user/admin-user", viewData)
}

func AdminUserFormPage(c *fiber.Ctx) error {
	db := database.DB
	authedUser, _ := middlewares.SelectAuthenticatedUser(c, db)

	var roles []models.Role
	db.NewSelect().Model(&roles).Scan(ctx)

	viewData := models.ViewData{
		Data: roles,
		Meta: models.Meta{
			PageTitle:  "Users",
			AuthedUser: authedUser,
			BaseURL:    BaseURL,
		},
	}

	return c.Render("pages/admin/user/user-form", viewData)
}
