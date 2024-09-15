package Handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/middlewares"
	"goHtmlBlog/models"
	"strconv"
)

func AdminDashboard(c *fiber.Ctx) error {
	db := database.DB
	user, _ := middlewares.SelectAuthenticatedUser(c, db)
	if middlewares.IsAuthorized(c, "", 0) {
		c.Status(fiber.StatusUnauthorized)
		return c.Redirect("/")
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

///////////////

///////////////

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

///////////////

///////////////

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

///////////////

///////////////

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

///////////////

///////////////

func AdminBlogsPage(c *fiber.Ctx) error {
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

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	db := database.DB
	var blogs []models.Blog
	blogLen, _ := db.NewSelect().Model(&blogs).Relation("User").Order("updated_at" +
		" DESC").Limit(limit).Offset(offset).ScanAndCount(ctx)

	for id, blog := range blogs {
		if blog.UserID == 0 {
			blogs[id].User = &models.User{Name: "this author deleted"}
		}
	}
	authedUser, _ := middlewares.SelectAuthenticatedUser(c, db)
	pagination := models.Pagination{}.Paginate(page, blogLen, limit)
	viewData := models.ViewData{
		Data: blogs,
		Meta: models.Meta{
			PageTitle:  "Blogs",
			AuthedUser: authedUser,
			BaseURL:    BaseURL,
			Pagination: pagination,
		},
	}
	return c.Render("pages/admin/blog/admin-blogs", viewData)
}

func AdminBlogPage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	db := database.DB
	authedUser, _ := middlewares.SelectAuthenticatedUser(c, db)
	var blog models.Blog
	db.NewSelect().Model(&blog).Where("\"blog\".\"id\" = ?", id).Relation("User").Scan(ctx)

	if blog.User == nil {
		blog.User = &models.User{Name: "this author deleted"}
	}

	viewData := models.ViewData{
		Data: blog,
		Meta: models.Meta{
			PageTitle:  "Blogs",
			AuthedUser: authedUser,
			BaseURL:    BaseURL,
		},
	}
	return c.Render("pages/admin/blog/admin-blog", viewData)
}

func AdminBlogFormPage(c *fiber.Ctx) error {
	db := database.DB
	authedUser, _ := middlewares.SelectAuthenticatedUser(c, db)

	viewData := models.ViewData{
		Meta: models.Meta{
			PageTitle:  "Blogs",
			AuthedUser: authedUser,
			BaseURL:    BaseURL,
		},
	}
	return c.Render("pages/admin/blog/blog-form", viewData)
}
