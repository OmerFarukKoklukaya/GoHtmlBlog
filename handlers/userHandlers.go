package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/middlewares"
	"goHtmlBlog/models"
	"strconv"
)

func LoginPage(c *fiber.Ctx) error {
	viewData := models.ViewData{
		Meta: models.Meta{
			PageTitle: "Login",
			BaseURL:   BaseURL,
		},
	}
	return c.Render("pages/login", viewData)
}
func RegisterPage(c *fiber.Ctx) error {
	err := c.Query("error")

	viewData := models.ViewData{
		Data: struct {
			Err string
		}{Err: err},
		Meta: models.Meta{
			PageTitle: "Register",
			BaseURL:   BaseURL,
		},
	}

	return c.Render("pages/register", viewData)

}
func ProfilePage(c *fiber.Ctx) error {
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
	var user models.User
	var author models.User
	var authorized = false
	user, err = middlewares.SelectAuthenticatedUser(c, db)
	authorID := c.Params("id")
	fmt.Println("AuthorName: ", authorID)
	if len(authorID) == 0 {
		author = user
	} else {
		err = db.NewSelect().Model(&author).Where("id = ?", authorID).Relation("Blogs").Scan(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}

	if middlewares.IsAuthorized(c, "", author.ID) {
		authorized = true
	}

	var blogs []models.Blog
	blogLen, _ := db.NewSelect().Model(&blogs).Relation("User").Order("id DESC").Where("user_id = ?", author.ID).Limit(limit).Offset(offset).ScanAndCount(ctx)
	fmt.Println("PROFILE PAGE: BLOGS:", blogs)

	for _, blog := range blogs {
		if blog.User == nil {
			blog.User = &models.User{Name: "this author deleted"}
		}

	}

	paginate := models.Pagination{}.Paginate(page, blogLen, limit)

	viewData := models.ViewData{
		Data: struct {
			User  models.User
			Blogs []models.Blog
		}{
			User:  author,
			Blogs: blogs,
		},
		Meta: models.Meta{
			PageTitle:    author.Name,
			AuthedUser:   user,
			IsAuthorized: authorized,
			Pagination:   paginate,
			BaseURL:      BaseURL,
		},
	}

	return c.Render("pages/profile-page", viewData)
}
