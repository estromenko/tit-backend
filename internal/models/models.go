package models

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           uint64 `bun:"id,pk,autoincrement" json:"-"`
	Email        string `bun:"email,unique,notnull" json:"email" validate:"required,email"`
	PasswordHash string `bun:"password_hash,notnull" json:"-"`
	IsSuperUser  bool   `bun:"is_super_user,notnull" json:"isSuperUser"`
	IsActive     bool   `bun:"is_active,notnull"`
	Port         *uint  `bun:"port" json:"port"`

	Password string `bun:"-" json:"password,omitempty" validate:"required,min=8,max=256"`
	Token    string `bun:"-" json:"token,omitempty"`
}
