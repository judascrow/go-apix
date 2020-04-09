package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username  string `gorm:"not null;unique"`
	Password  string `gorm:"not null"`
	FirstName string `gorm:"varchar(255);not null"`
	LastName  string `gorm:"varchar(255);not null"`
	Email     string `gorm:"unique_index"`
	Bio       string `gorm:"size:1024"`
	Image     *string
}
