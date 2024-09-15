package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/models"
)

func PostPermission(c *fiber.Ctx) error {
	var permission models.Permission
	err := json.Unmarshal(c.Body(), &permission)
	if err != nil || permission.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON("name cannot be empty")
	}

	db := database.DB
	var dummPermission models.Permission
	db.NewSelect().Model(&dummPermission).Where("name = ?", permission.Name).Scan(ctx)
	if dummPermission.Name != "" {
		return c.Status(fiber.StatusBadRequest).JSON("Already have this permission")
	}
	_, err = db.NewInsert().Model(&permission).Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("")
	}
	return c.JSON("permissions")
}

func PutPermission(c *fiber.Ctx) error {
	var permission models.Permission
	err := json.Unmarshal(c.Body(), &permission)
	if err != nil {
		fmt.Println("err")
	}
	if len(permission.Name) == 0 {
		return c.JSON("Permission name cannot be empty")
	}
	db := database.DB
	var oldPermission models.Permission
	oldPermission.ID, err = c.ParamsInt("id")
	if oldPermission.ID == 0 || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("bad request")
	}

	_, err = db.NewUpdate().Model(&permission).ExcludeColumn("id").Where("id = ?", oldPermission.ID).Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("update problem")
	}
	return c.SendStatus(fiber.StatusOK)
}

func DeletePermission(c *fiber.Ctx) error {
	permissionID, err := c.ParamsInt("id")
	if permissionID == 0 || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("bad request")
	}
	db := database.DB
	_, err = db.NewDelete().Model(&models.Permission{}).Where("id = ?", permissionID).Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("delete problem")
	}
	return c.JSON("")
}
