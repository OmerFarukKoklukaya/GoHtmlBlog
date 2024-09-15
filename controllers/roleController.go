package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"goHtmlBlog/database"
	"goHtmlBlog/models"
)

func PostRole(c *fiber.Ctx) error {
	var roleAndPermissions = struct {
		Name        string `json:"name"`
		Permissions []int  `json:"permissions"`
	}{}
	err := json.Unmarshal(c.Body(), &roleAndPermissions)
	if err != nil {
		fmt.Println("err")
	}
	if len(roleAndPermissions.Name) == 0 {
		return c.JSON("roleAndPermissions name cannot be empty")
	}
	var role = models.Role{Name: roleAndPermissions.Name}
	db := database.DB
	db.NewInsert().Model(&role).Exec(ctx)
	var rtp models.RoleToPermission
	for _, permissionID := range roleAndPermissions.Permissions {
		rtp = models.RoleToPermission{RoleID: role.ID, PermissionID: permissionID}
		db.NewInsert().Model(&rtp).Exec(ctx)
	}
	return c.JSON("roles")
}

func PutRole(c *fiber.Ctx) error {
	var role models.Role
	var err error

	role.ID, _ = c.ParamsInt("id")
	if role.ID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON("corrupted id")
	}
	var roleAndPermissions = struct {
		Name        string `json:"name"`
		Permissions []int  `json:"permissions"`
	}{}
	err = json.Unmarshal(c.Body(), &roleAndPermissions)
	if err != nil {
		fmt.Println("err")
	}
	if len(roleAndPermissions.Name) == 0 {
		return c.JSON("role name cannot be empty")
	}
	role.Name = roleAndPermissions.Name
	db := database.DB
	roleID, err := c.ParamsInt("id")
	if roleID == 0 || err != nil {
		return c.JSON("bad request")
	}

	var rtp models.RoleToPermission
	db.NewDelete().Model(&rtp).Where("role_id = ?", roleID).Exec(ctx)
	_, err = db.NewUpdate().Model(&role).Column("name").Where("id = ?", roleID).Exec(ctx)
	if err != nil {
		return c.JSON("update problem")
	}

	for _, permissionID := range roleAndPermissions.Permissions {
		rtp = models.RoleToPermission{RoleID: role.ID, PermissionID: permissionID}
		db.NewInsert().Model(&rtp).Exec(ctx)
	}

	return c.JSON("roles")
}

func DeleteRole(c *fiber.Ctx) error {
	roleID, err := c.ParamsInt("id")
	if roleID == 0 || err != nil {
		return c.JSON("bad request")
	}
	db := database.DB
	_, err = db.NewDelete().Model(&models.Role{}).Where("id = ?", roleID).Exec(ctx)
	if err != nil {
		return c.JSON("delete problem")
	}
	return c.JSON("roles")
}
