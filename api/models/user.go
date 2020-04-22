package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username  string `json:"username" form:"username" gorm:"not null;unique_index" binding:"required"`
	Password  string `json:"password" form:"password" gorm:"not null" binding:"required"`
	FirstName string `json:"firstName" form:"firstName" gorm:"not null" binding:"required"`
	LastName  string `json:"lastName" form:"lastName" gorm:"not null" binding:"required"`
	Email     string `json:"email" form:"email" gorm:"unique_index"`
	Avatar    string `json:"avatar" form:"avatar" `
}

func (u User) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"id":        u.ID,
		"username":  u.Username,
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"email":     u.Email,
		"avatar":    u.Avatar,
	}
}
