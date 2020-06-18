package models

import (
	"strings"

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
	Status    string `json:"status" form:"status" sql:"type:enum('A','I');DEFAULT:'A'"`
	Avatar    string `json:"avatar" form:"avatar" `

	Roles     []Role     `json:"roles" gorm:"many2many:users_roles;"`
	UserRoles []UserRole `json:"users_roles" gorm:"foreignkey:UserId"`
}

type ChangePassword struct {
	CurrentPassword string `json:"current_password" form:"current_password" binding:"required" example:"password"` // รหัสผ่านปัจจุบัน
	NewPassword     string `json:"new_password" form:"new_password" binding:"required" example:"password123"`      // รหัสผ่านใหม่
}

type UploadAvatar struct {
	Avatar string `json:"avatar" form:"avatar" binding:"required" example:"avatar.png"` // รูปภาพ
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

	roles := []map[string]interface{}{}
	for _, role := range u.Roles {
		roles = append(roles, role.Serialize())
	}

	replaceAllFlag := -1

	return map[string]interface{}{
		"id":        u.ID,
		"username":  u.Username,
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"email":     u.Email,
		"slug":      u.Slug,
		"status":    u.Status,
		"avatar":    strings.Replace(u.Avatar, "\\", "/", replaceAllFlag),
		"roles":     roles,
	}
}

func (user *User) GetUserStatusAsString() string {
	switch user.Status {
	case "A":
		return "Active"
	case "I":
		return "Inctive"
	case "C":
		return "Cancel"
	default:
		return "Unknown"
	}
}

func (user *User) IsAdmin() bool {
	for _, role := range user.Roles {
		if role.Name == "ROLE_ADMIN" {
			return true
		}
	}
	return false
}
func (user *User) IsNotAdmin() bool {
	return !user.IsAdmin()
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
	SwagUserBody
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

type SwagUpdateUserResponse struct {
	SwagUpdateBase
	Data struct {
		Users SwagUser `json:"users"`
	} `json:"data"`
}

type SwagChangePasswordResponse struct {
	Success bool        `json:"success" example:"true"`                         // ผลการเรียกใช้งาน
	Status  int         `json:"status" example:"200"`                           // HTTP Status Code
	Message string      `json:"message" example:"Change Password Successfully"` // ข้อความตอบกลับ
	Data    interface{} `json:"data" `                                          // ข้อมูล
}

type SwagUploadAvatarResponse struct {
	Success bool        `json:"success" example:"true"`                         // ผลการเรียกใช้งาน
	Status  int         `json:"status" example:"200"`                           // HTTP Status Code
	Message string      `json:"message" example:"Uploaded Avatar Successfully"` // ข้อความตอบกลับ
	Data    interface{} `json:"data" `                                          // ข้อมูล
}
