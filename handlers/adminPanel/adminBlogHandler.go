package adminHandlers

import (
	"github.com/gofiber/fiber/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/middlewares"
	"goHtmlBlog/models"
	"strconv"
)

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
