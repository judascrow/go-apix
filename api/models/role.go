package models

import "github.com/jinzhu/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string
	Users       []User     `gorm:"many2many:users_roles;"`
	UserRoles   []UserRole `gorm:"foreignkey:RoleID"`
}

type UserRole struct {
	User   User `gorm:"association_foreignkey:UserID"`
	UserID uint
	Role   User `gorm:"association_foreignkey:RoleID"`
	RoleID uint
}

func (UserRole) TableName() string {
	return "users_roles"
}

func (r Role) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"id":          r.ID,
		"name":        r.Name,
		"description": r.Description,
	}
}
