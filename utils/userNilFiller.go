package utils

import (
	"errors"
	"goHtmlBlog/models"
)

func FillNilUser(user *models.User) error {
	if user == nil {
		user.Name = "this user is deleted"
		user.RoleID = 0
		user.Role.Name = "empty role"
	}
	return errors.New("user is not nil")
}
