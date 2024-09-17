package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/middlewares"
	"goHtmlBlog/models"
	"goHtmlBlog/utils"
	"os"
	"strconv"
	"strings"
)

func GetImage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if id == 0 || err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendFile("./images/blog-"+strconv.Itoa(id)+".jpeg", true)
}

func PostImage(c *fiber.Ctx) error {
	db := database.DB

	form, err := c.MultipartForm()
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var blog models.Blog
	var authedUser models.User
	authedUser, err = middlewares.SelectAuthenticatedUser(c, db)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	db.NewSelect().Model(&blog).Where("id = ?", id).Scan(ctx)
	if blog.ID == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if blog.UserID != authedUser.ID && !middlewares.IsAuthorized(c, "", blog.ID) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if err != nil {
		fmt.Println("POST IMAGE: Img read problem:", err)
		return c.JSON("Img read problem")
	}

	if form.File["image"] == nil {
		fmt.Println("POST IMAGE: Incoming image problem: file is empty")
		return c.Status(fiber.StatusBadRequest).JSON("Img cannot be founded")
	}
	file := form.File["image"][0]
	Header := strings.Split(file.Header["Content-Type"][0], "/")
	fileType := Header[0]
	if fileType != "image" {
		return c.SendString("unacceptable file type")
	}
	file.Filename = "blog-" + strconv.Itoa(blog.ID)
	err = utils.ImgScaleAndSave(file)
	if err != nil {
		fmt.Println("Img problem:", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Img problem")
	}

	blog.Image = "/api/blogs/images/" + strconv.Itoa(blog.ID)
	_, err = db.NewUpdate().Model(&blog).Column("image").Where("id = ?", blog.ID).Exec(ctx)
	if err != nil {
		fmt.Println("POST IMAGE: update blog err:", err)
	}

	return c.SendStatus(200)
}

func DeleteImage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var blog models.Blog
	var authedUser models.User
	db := database.DB
	authedUser, err = middlewares.SelectAuthenticatedUser(c, db)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	db.NewSelect().Model(&blog).Where("id = ?", id).Scan(ctx)
	if blog.ID == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if blog.UserID != authedUser.ID && !middlewares.IsAuthorized(c, "", blog.ID) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if id == 0 || err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	blog.Image = ""

	db.NewUpdate().Model(&blog).Column("image").Where("id = ?", id).Scan(ctx)
	os.Remove("./images/blog-" + strconv.Itoa(blog.ID) + ".jpeg")
	return c.JSON("Image Deleted")
}
