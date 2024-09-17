package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"goHtmlBlog/controllers"
	"goHtmlBlog/handlers"
	adminHandlers "goHtmlBlog/handlers/adminPanel"
	"goHtmlBlog/middlewares"
)

func Router(app *fiber.App) {
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))

	nonAuthedRoutes(app)
	app.Use(middlewares.AuthenticationMiddleware)
	app.Use(middlewares.AuthorizationMiddleware)
	authedRoutes(app)
}

func nonAuthedRoutes(app *fiber.App) {
	api := app.Group("/api")
	blogs := app.Group("/blogs")

	////////
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Post("/logout", controllers.Logout)
	////////
	app.Get("/", func(c *fiber.Ctx) error { return c.Redirect("/blogs") })
	app.Get("/register", handlers.RegisterPage)
	app.Get("/login", handlers.LoginPage)
	app.Get("/user/:id<int>?", handlers.ProfilePage)
	////////

	////////
	app.Get("/assets/:pageName?/styles*", controllers.CSSController)
	app.Get("/scripts/:pageName?/scripts*", controllers.JSController)
	api.Get("/blogs/images/:id", controllers.GetImage)
	////////

	////////
	blogs.Get("/", handlers.BlogsPage)
	blogs.Get("/:id<int>", handlers.BlogPage)

}

func authedRoutes(app *fiber.App) {
	api := app.Group("/api")
	blogs := app.Group("/blogs")
	////////

	////////
	api.Post("/blogs/images/:id", controllers.PostImage)
	api.Delete("/blogs/images/:id", controllers.DeleteImage)
	////////

	////////
	blogs.Get("/editor/:id<int>?", handlers.BlogEditorPage)
	////////

	////////
	api.Post("/blog", controllers.PostBlog)
	////////

	////////
	user := api.Group("/user")
	user.Post("/", controllers.InsertUser)
	user.Put("/:id<id>?", controllers.UpdateUser)
	user.Put("/password/:id<int>", controllers.UpdatePassword)
	user.Delete("/:id<int>", controllers.DeleteProfile)

	////////
	api.Put("/blog/:id<int>", controllers.PutBlog)
	api.Delete("/blog/:id<int>", controllers.DeleteBlog)
	////////
	roles := api.Group("/roles")
	roles.Post("/", controllers.PostRole)
	roles.Put("/:id", controllers.PutRole)
	roles.Delete("/:id", controllers.DeleteRole)
	////////
	permissions := api.Group("/permissions")
	permissions.Post("/", controllers.PostPermission)
	permissions.Put("/:id", controllers.PutPermission)
	permissions.Delete("/:id", controllers.DeletePermission)
	////////

	////////
	admin := app.Group("/admin")
	admin.Get("/", adminHandlers.AdminDashboard)
	////////
	admin.Get("/users", adminHandlers.AdminUsersPage)
	admin.Get("/users/:id<int>", adminHandlers.AdminUserPage)
	admin.Get("/users/form", adminHandlers.AdminUserFormPage)
	////////
	admin.Get("/roles", adminHandlers.AdminRolesPage)
	admin.Get("/roles/:id<int>", adminHandlers.AdminRolePage)
	admin.Get("/roles/form", adminHandlers.AdminRoleFormPage)
	////////
	admin.Get("/permissions", adminHandlers.AdminPermissionsPage)
	admin.Get("/permissions/:id<int>", adminHandlers.AdminPermissionPage)
	admin.Get("/permissions/form", adminHandlers.AdminPermissionsFormPage)
	////////
	admin.Get("/blogs", adminHandlers.AdminBlogsPage)
	admin.Get("/blogs/:id<int>", adminHandlers.AdminBlogPage)
	admin.Get("/blogs/form", adminHandlers.AdminBlogFormPage)
}
