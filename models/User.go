package models

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete"`
	ID        int       `bun:"id,pk,autoincrement" json:"id"`
	Name      string    `bun:"name" json:"name"`
	Password  string    `bun:"password" json:"password"`
	Blogs     []*Blog   `bun:"rel:has-many,join:id=user_id" json:"blogs"`
	Role      *Role     `bun:"rel:belongs-to,join:role_id=id" json:"role"`
	RoleID    int       `bun:"role_id" json:"role_id"`
}
