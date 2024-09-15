package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func JSController(c *fiber.Ctx) error {
	pageName := c.Params("pageName")
	return c.SendFile("./views/scripts/"+pageName+"/scripts.js", true)
}

func CSSController(c *fiber.Ctx) error {
	pageName := c.Params("pageName")
	return c.SendFile("./views/assets/"+pageName+"/styles.css", true)
}
