package controllers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/middlewares"
	"goHtmlBlog/models"
	"strconv"
)

var ctx = context.Background()

func PostBlog(c *fiber.Ctx) error {
	db := database.DB
	fmt.Println(string(c.Body()))
	newBlog := models.Blog{}
	userId := c.FormValue("user_id")
	newBlog.UserID, _ = strconv.Atoi(userId)
	newBlog.Title = c.FormValue("title")
	newBlog.Body = c.FormValue("body")
	newBlog.Summary = c.FormValue("summary")
	if newBlog.Title == "" || newBlog.Body == "" || newBlog.Summary == "" {
		return c.Status(fiber.StatusBadRequest).JSON("Title nor body cannot be empty")
	}
	/*if len(newBlog.Body) < 50 {
		newBlog.Summary = newBlog.Body
	} else {
		newBlog.Summary = newBlog.Body[:50] + "..."
	}*/

	user := models.User{}
	user, _ = middlewares.SelectAuthenticatedUser(c, db)
	if !middlewares.IsAuthorized(c, "", newBlog.ID) {
		newBlog.UserID = user.ID
	}

	fmt.Println("NEW BLOG:", newBlog)
	_, err := db.NewInsert().Model(&newBlog).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)

	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": newBlog})
}
func PutBlog(c *fiber.Ctx) error {
	db := database.DB
	var err error
	var oldBlog models.Blog
	var blog models.Blog
	var authedUser models.User

	blog.ID, _ = c.ParamsInt("id")
	blog.Title = c.FormValue("title")
	blog.Body = c.FormValue("body")
	userID := c.FormValue("user_id")
	blog.UserID, _ = strconv.Atoi(userID)
	blog.Summary = c.FormValue("summary")

	if blog.Title == "" || blog.Body == "" || blog.Summary == "" {
		return c.Status(fiber.StatusBadRequest).JSON("Title nor body cannot be empty")
	}

	/*
		if len(blog.Body) < 50 {
			blog.Summary = blog.Body
		} else {
			blog.Summary = blog.Body[:50] + "..."
		}
	*/
	db.NewSelect().Model(&oldBlog).Where("id = ?", blog.ID).Scan(ctx)
	if oldBlog.ID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON("bad request")
	}
	authedUser, _ = middlewares.SelectAuthenticatedUser(c, db)
	if !middlewares.IsAuthorized(c, "", blog.ID) || blog.UserID <= 0 {
		blog.UserID = authedUser.ID
	}

	fmt.Println(blog.UserID)

	_, err = db.NewUpdate().Model(&blog).Column("title", "body", "summary", "user_id").Where("id = ?", blog.ID).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return c.JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": blog})
}
func DeleteBlog(c *fiber.Ctx) error {
	db := database.DB
	var blogId, _ = c.ParamsInt("id")
	var blog models.Blog
	var authedUser models.User
	if blogId <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	authedUser, _ = middlewares.SelectAuthenticatedUser(c, db)

	db.NewSelect().Model(&blog).Where("id = ?", blogId).Scan(ctx)

	if !middlewares.IsAuthorized(c, "", blog.UserID) && blog.UserID != authedUser.ID {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	db.NewDelete().Model(&models.Blog{}).Where("id = ?", blogId).Exec(ctx)
	return c.Status(fiber.StatusOK).JSON("Blog successfully deleted")
}
