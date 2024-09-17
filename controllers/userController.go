package controllers

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/middlewares"
	"goHtmlBlog/models"
	"goHtmlBlog/utils"
	"strconv"
)

func Register(c *fiber.Ctx) error {
	db := database.DB
	userName := c.FormValue("name")
	if len(userName) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Name cannot be empty")
	}
	password := []byte(c.FormValue("password"))
	if len(password) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Password cannot be empty")
	}
	passwordVer := []byte(c.FormValue("password_verification"))
	if len(passwordVer) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Password verification cannot be empty")
	}
	fmt.Println(password, passwordVer)
	if !bytes.Equal(password, passwordVer) {
		return c.Status(fiber.StatusBadRequest).JSON("Password verification failed")
	}

	var user models.User
	db.NewSelect().Model(&user).Where("name = ?", userName).Relation("Blogs").Scan(ctx)
	if len(user.Name) != 0 {
		return c.Status(fiber.StatusBadRequest).JSON("This name already taken")
	}

	var role models.Role
	db.NewSelect().Model(&role).Where("name = ?", "basic").Scan(ctx)

	user = models.User{Name: userName, RoleID: role.ID}
	user.ChangePassword(password)
	db.NewInsert().Model(&user).Exec(ctx)

	return c.SendStatus(fiber.StatusOK)
}
func Login(c *fiber.Ctx) error {
	db := database.DB
	user := models.User{}
	userName := c.FormValue("name")
	if len(userName) == 0 {
		fmt.Println("Login page problem: name cannot be empty")
		return c.Status(fiber.StatusBadRequest).JSON("Name cannot be empty")
	}
	password := []byte(c.FormValue("password"))
	if len(password) == 0 {
		fmt.Println("Login page problem: password cannot be empty")
		return c.Status(fiber.StatusBadRequest).JSON("Password cannot be empty")
	}

	err := db.NewSelect().Model(&user).Where("name = ?", userName).Scan(ctx)
	if !user.CheckPassword(password) {
		fmt.Println("Login err:", err)
		return c.Status(fiber.StatusBadRequest).JSON("Wrong Name or Password")
	}

	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: utils.GenerateToken(user.ID),
	})

	return c.SendStatus(fiber.StatusOK)
}
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: "",
	})
	c.Method("GET")
	return c.Redirect("http://localhost:3000/blogs")
}
func InsertUser(c *fiber.Ctx) error {
	db := database.DB
	var user models.User
	user.Name = c.FormValue("name")
	password := []byte(c.FormValue("password"))
	user.RoleID, _ = strconv.Atoi(c.FormValue("role_id"))
	if user.Name == "" || len(password) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Name or Password cannot be empty")
	}
	user.ChangePassword(password)
	if user.RoleID <= 0 {
		user.RoleID = 2
	}

	authedUser, err := middlewares.SelectAuthenticatedUser(c, db)
	if authedUser.ID == 0 || err != nil || !middlewares.IsAuthorized(c, "", 1) {
		return c.SendStatus(fiber.StatusUnauthorized)
		fmt.Println(err)
	}

	db.NewInsert().Model(&user).Exec(ctx)

	return c.SendStatus(200)
}
func UpdateUser(c *fiber.Ctx) error {
	db := database.DB
	var user, oldUser models.User
	user.ID, _ = c.ParamsInt("id")

	user.Name = c.FormValue("name")
	user.RoleID, _ = strconv.Atoi(c.FormValue("role_id"))
	var authedUser models.User
	authedUser, err := middlewares.SelectAuthenticatedUser(c, db)
	if authedUser.ID <= 0 || err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
		fmt.Println(err)
	}
	if user.ID <= 0 {
		user.ID = authedUser.ID
		oldUser = authedUser
	} else {
		db.NewSelect().Model(&oldUser).Where("id = ?", user.ID).Relation("Role").Exec(ctx)
	}
	if user.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON("name cannot be empty")
	}
	if !middlewares.IsAuthorized(c, "", user.ID) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if !middlewares.IsAuthorized(c, "roles", 1) || user.RoleID <= 0 {
		user.RoleID = oldUser.RoleID
	}

	db.NewUpdate().Model(&user).Column("name", "role_id").Where("id = ?", user.ID).Exec(ctx)

	return c.Status(200).JSON(fiber.Map{"user": user})
}
func DeleteProfile(c *fiber.Ctx) error {
	db := database.DB
	authedUser, _ := middlewares.SelectAuthenticatedUser(c, db)
	id, _ := c.ParamsInt("id")
	fmt.Println("DELETEPROFILE: ID:", id)
	if id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON("unacceptable id")
	}
	if id != authedUser.ID {
		middlewares.IsAuthorized(c, "", id)
	}

	db.NewDelete().Model(&models.User{}).Where("id = ?", id).Exec(ctx)
	return c.SendStatus(fiber.StatusOK)
}
func UpdatePassword(c *fiber.Ctx) error {
	oldPassword := []byte(c.FormValue("password"))
	if len(oldPassword) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Password Cannot Be Empty")
	}
	newPassword := []byte(c.FormValue("new_password"))
	if len(newPassword) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("New Password Cannot Be Empty")
	}
	passVerification := []byte(c.FormValue("new_password_verification"))
	if len(passVerification) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("New Password verification Cannot Be Empty")
	}
	if !bytes.Equal(newPassword, passVerification) {
		return c.Status(fiber.StatusBadRequest).JSON("Passwords doesnt match")
	}
	db := database.DB
	var user models.User
	user, err := middlewares.SelectAuthenticatedUser(c, db)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized")
	}
	if !user.CheckPassword(oldPassword) {
		return c.Status(fiber.StatusBadRequest).JSON("Wrong password")
	}

	user.ChangePassword(newPassword)

	db.NewUpdate().Model(&user).Where("id = ?", user.ID).Exec(ctx)

	return c.SendStatus(fiber.StatusOK)
}
