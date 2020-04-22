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
}

func (u *User) BeforeSave() (err error) {
	u.Slug = slug.Make(u.Username)
	return
}

func (u User) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"id":        u.ID,
		"username":  u.Username,
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"email":     u.Email,
		"slug":      u.Slug,
		"avatar":    u.Avatar,
	}
}
