package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/middlewares"
	"goHtmlBlog/models"
	"strconv"
)

const BaseURL = "http://localhost:3000"

var ctx = context.Background()

func BlogsPage(c *fiber.Ctx) error {
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
	var blogs []models.Blog
	blogLen, _ := db.NewSelect().Model(&blogs).Relation("User").Order("id DESC").Limit(limit).Offset(offset).ScanAndCount(ctx)

	for _, blog := range blogs {
		if blog.User == nil {
			blog.User = &models.User{Name: "this author deleted"}
		}
	}

	var user models.User
	user, _ = middlewares.SelectAuthenticatedUser(c, db)

	paginate := models.Pagination{}.Paginate(page, blogLen, limit)

	viewData := models.ViewData{
		Data: blogs,
		Meta: models.Meta{
			PageTitle:  "Blogs",
			AuthedUser: user,
			Pagination: paginate,
			BaseURL:    BaseURL,
		},
	}
	return c.Render("pages/blog/blog-home", viewData)
}
func BlogPage(c *fiber.Ctx) error {
	blogID, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	db := database.DB

	blog := models.Blog{}
	err = db.NewSelect().Model(&blog).Where(" \"blog\".\"id\" = ?", blogID).Relation("User").Scan(ctx)

	if blog.User == nil {
		blog.User = &models.User{Name: "this author deleted"}
	}

	var authedUser models.User
	var isAuthorized bool

	authedUser, _ = middlewares.SelectAuthenticatedUser(c, db)
	isAuthorized = middlewares.IsAuthorized(c, "", blog.ID)

	viewData := models.ViewData{
		Data: blog,
		Meta: models.Meta{
			PageTitle:    blog.Title,
			AuthedUser:   authedUser,
			IsAuthorized: isAuthorized,
			BaseURL:      BaseURL,
		},
	}
	return c.Render("pages/blog/blog-page", viewData)
}
func BlogEditorPage(c *fiber.Ctx) error {
	db := database.DB
	var user models.User
	var blog = models.Blog{}
	user, _ = middlewares.SelectAuthenticatedUser(c, db)

	blogId := 0
	blogId, _ = c.ParamsInt("id")
	if blogId != 0 {
		db.NewSelect().Model(&blog).Where(" \"blog\".\"id\" = ?", blogId).Scan(ctx)
	}

	viewData := models.ViewData{
		Data: blog,
		Meta: models.Meta{
			PageTitle:  "Editor",
			AuthedUser: user,
			BaseURL:    BaseURL,
		},
	}
	return c.Render("pages/blog/editor-page", viewData)
}
