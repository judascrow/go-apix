package services

import (
	"os"
	"strconv"

	"github.com/judascrow/go-apix/api/infrastructure"
)

func CreateOne(data interface{}) error {
	db := infrastructure.GetDB()
	err := db.Create(data).Error
	return err
}

type PageMeta struct {
	TotalItemsCount   int `json:"totalItemsCount"`
	CurrentItemsCount int `json:"currentItemsCount"`
	PageSize          int `json:"pageSize"`
	Page              int `json:"page"`
}

func getPageMeta(pageSizeStr, pageStr string) (pageMeta PageMeta) {

	pageMeta = PageMeta{}

	size, err := strconv.Atoi(os.Getenv("APP_API_PAGE_SIZE"))
	if err != nil {
		size = 10
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = size
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	pageMeta.PageSize = pageSize
	pageMeta.Page = page

	return pageMeta
}
