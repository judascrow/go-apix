package services

import (
	"github.com/judascrow/go-api-starter/api/infrastructure"
	"github.com/judascrow/go-api-starter/api/models"
)

func FindAllUsers(pageSizeStr, pageStr string) ([]models.User, PageMeta, error) {
	db := infrastructure.GetDB()
	pageMeta := getPageMeta(pageSizeStr, pageStr)

	var users []models.User
	var count int

	db.Model(&models.User{}).Count(&count)
	err := db.Offset((pageMeta.Page - 1) * pageMeta.PageSize).Limit(pageMeta.PageSize).Find(&users).Error

	pageMeta.Count = count

	return users, pageMeta, err
}

func FindOneUser(id uint) (models.User, error) {
	db := infrastructure.GetDB()
	var user models.User
	err := db.First(&user, id).Error
	return user, err
}
