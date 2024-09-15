package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"goHtmlBlog/Handlers"
	"goHtmlBlog/controllers"
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
	files := api.Group("/files")

	////////
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Post("/logout", controllers.Logout)
	////////
	app.Get("/", func(c *fiber.Ctx) error { return c.Redirect("/blogs") })
	app.Get("/register", Handlers.RegisterPage)
	app.Get("/login", Handlers.LoginPage)
	app.Get("/user/:id<int>?", Handlers.ProfilePage)
	////////

	////////
	app.Get("/assets/:pageName?/styles*", controllers.CSSController)
	app.Get("/scripts/:pageName?/scripts*", controllers.JSController)
	files.Get("/images/blog/:id", controllers.GetImage)
	////////

	////////
	blogs.Get("/", Handlers.BlogsPage)
	blogs.Get("/:id<int>", Handlers.BlogPage)

}

func authedRoutes(app *fiber.App) {
	api := app.Group("/api")
	blogs := app.Group("/blogs")
	files := api.Group("/files")
	////////

	////////
	files.Post("/images/blog/:id", controllers.PostImage)
	files.Delete("/images/blog/:id", controllers.DeleteImage)
	////////

	////////
	blogs.Get("/editor/:id<int>?", Handlers.BlogEditorPage)
	////////

	////////
	api.Post("/blog", controllers.PostBlog)
	////////

	////////
	user := api.Group("/user")
	user.Post("/", controllers.InsertUser)
	user.Put("/:id<id>?", controllers.UpdateUser)
	user.Put("/password", controllers.UpdatePassword)
	user.Delete("/:id<int>", controllers.DeleteProfile)

	////////
	api.Put("/blog/:id<int>", controllers.PutBlog)
	api.Delete("/blog/:id<int>", controllers.DeleteBlog)
	////////
	roles := api.Group("/roles")
	roles.Post("/", controllers.PostRole)
	roles.Put("/:id", controllers.PutRole)
	roles.Delete("/id", controllers.DeleteRole)
	////////
	permissions := api.Group("/permissions")
	permissions.Post("/", controllers.PostPermission)
	permissions.Put("/:id", controllers.PutPermission)
	permissions.Delete("/:id", controllers.DeletePermission)
	////////

	////////
	admin := app.Group("/admin")
	admin.Get("/", Handlers.AdminDashboard)
	////////
	admin.Get("/users", Handlers.AdminUsersPage)
	admin.Get("/users/:id<int>", Handlers.AdminUserPage)
	admin.Get("/users/form", Handlers.AdminUserFormPage)
	////////
	admin.Get("/roles", Handlers.AdminRolesPage)
	admin.Get("/roles/:id<int>", Handlers.AdminRolePage)
	admin.Get("/roles/form", Handlers.AdminRoleFormPage)
	////////
	admin.Get("/permissions", Handlers.AdminPermissionsPage)
	admin.Get("/permissions/:id<int>", Handlers.AdminPermissionPage)
	admin.Get("/permissions/form", Handlers.AdminPermissionsFormPage)
	////////
	admin.Get("/blogs", Handlers.AdminBlogsPage)
	admin.Get("/blogs/:id<int>", Handlers.AdminBlogPage)
	admin.Get("/blogs/form", Handlers.AdminBlogFormPage)
}
