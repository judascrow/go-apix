package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username  string `json:"username" form:"username" gorm:"type:varchar(50);not null;unique_index" binding:"required"`
	Password  string `json:"password" form:"password" gorm:"not null" binding:"required"`
	FirstName string `json:"firstName" form:"firstName" gorm:"type:varchar(100);not null" binding:"required"`
	LastName  string `json:"lastName" form:"lastName" gorm:"type:varchar(100);not null" binding:"required"`
	Email     string `json:"email" form:"email" gorm:"type:varchar(100);unique_index"`
	Slug      string `json:"slug" form:"slug" uri:"slug"  gorm:"type:varchar(50);unique_index"`
	Avatar    string `json:"avatar" form:"avatar" `

	Roles     []Role     `json:"roles" gorm:"many2many:users_roles;"`
	UserRoles []UserRole `json:"users_roles" gorm:"foreignkey:UserId"`
}

func (u *User) BeforeSave(db *gorm.DB) (err error) {
	u.Slug = slug.Make(u.Username)
	if len(u.Roles) == 0 {
		userRole := Role{}
		db.Model(&Role{}).Where("name = ?", "ROLE_USER").First(&userRole)
		u.Roles = append(u.Roles, userRole)
	}
	return
}

func (u User) Serialize() map[string]interface{} {

	r := []map[string]interface{}{}
	for _, role := range u.Roles {
		r = append(r, role.Serialize())
	}

	return map[string]interface{}{
		"id":        u.ID,
		"username":  u.Username,
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"email":     u.Email,
		"slug":      u.Slug,
		"avatar":    u.Avatar,
		"roles":     r,
	}
}

type SwagUser struct {
	SwagID
	SwagUserBody
}

type SwagUserPassword struct {
	Password string `json:"password" example:"pass1234"` // รหัสผ่าน
}

type SwagUserBody struct {
	Username  string `json:"username" example:"user01"`        // Username
	FirstName string `json:"firstName" example:"john"`         // ชื่อ
	LastName  string `json:"lastName" example:"doe"`           // นามสกุล
	Email     string `json:"email" example:"user01@email.com"` // อีเมล์
	Slug      string `json:"slug" example:"user01"`            // Slug
	Avatar    string `json:"avatar" example:"user01.png"`      // รูป Avatar
}

type SwagUserBodyIncludePassword struct {
	Username  string `json:"username" example:"user01"`        // Username
	FirstName string `json:"firstName" example:"john"`         // ชื่อ
	LastName  string `json:"lastName" example:"doe"`           // นามสกุล
	Email     string `json:"email" example:"user01@email.com"` // อีเมล์
	Slug      string `json:"slug" example:"user01"`            // Slug
	Avatar    string `json:"avatar" example:"user01.png"`      // รูป Avatar
	SwagUserPassword
}

type SwagGetAllUsersResponse struct {
	SwagGetBase
	Data struct {
		PageMeta SwagPageMeta `json:"pageMeta"`
		Users    []SwagUser   `json:"users"`
	} `json:"data"`
}

type SwagGetUserResponse struct {
	SwagGetBase
	Data struct {
		Users SwagUser `json:"users"`
	} `json:"data"`
}

type SwagCreateUserResponse struct {
	SwagCreateBase
	Data struct {
		Users SwagUser `json:"users"`
	} `json:"data"`
}
