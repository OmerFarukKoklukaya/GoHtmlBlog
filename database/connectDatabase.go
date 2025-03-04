package database

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"goHtmlBlog/models"
)

var DB *bun.DB

func ConnectDatabase() {
	dns := "postgres://username:password@localhost:5432/blog?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dns)))

	DB = bun.NewDB(sqldb, pgdialect.New())
	ctx := context.Background()

	DB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	DB.NewCreateTable().IfNotExists().Model(&models.User{}).Exec(ctx)
	DB.NewCreateTable().IfNotExists().Model(&models.Blog{}).Exec(ctx)
	DB.NewCreateTable().IfNotExists().Model((*models.RoleToPermission)(nil)).Exec(ctx)
	DB.NewCreateTable().IfNotExists().Model(&models.Role{}).Exec(ctx)
	DB.NewCreateTable().IfNotExists().Model(&models.Permission{}).Exec(ctx)

}
