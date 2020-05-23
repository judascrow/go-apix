package responses

import (
	"fmt"
	"math"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/judascrow/go-apix/api/services"
)

type empty map[string]interface{}

func JSON(c *gin.Context, statusCode int, dataName string, data interface{}, message interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  statusCode,
		"success": true,
		"data": map[string]interface{}{
			dataName: data,
		},
		"message": message,
	})

}

func JSONLIST(c *gin.Context, statusCode int, dataName string, data interface{}, message interface{}, p services.PageMeta) {

	pageSize := p.PageSize
	page := p.Page
	totalItemsCount := p.TotalItemsCount
	currentItemsCount := p.CurrentItemsCount

	pageMeta := map[string]interface{}{}
	pageMeta["offset"] = (page - 1) * pageSize
	pageMeta["requestedPageSize"] = pageSize
	pageMeta["currentPageNumber"] = page
	pageMeta["currentItemsCount"] = currentItemsCount
	pageMeta["totalItemsCount"] = totalItemsCount

	pageMeta["prevPageNumber"] = 1
	totalPagesCount := int(math.Ceil(float64(totalItemsCount) / float64(pageSize)))
	pageMeta["totalPagesCount"] = totalPagesCount

	if page < totalPagesCount {
		pageMeta["hasNextPage"] = true
		pageMeta["nextPageNumber"] = page + 1
	} else {
		pageMeta["hasNextPage"] = false
		pageMeta["nextPageNumber"] = 1
	}
	if page > 1 {
		pageMeta["prevPageNumber"] = page - 1
	} else {
		pageMeta["hasPrevPage"] = false
		pageMeta["prevPageNumber"] = 1
	}

	pageMeta["nextPageUrl"] = fmt.Sprintf("%v?page=%d&pageSize=%d", c.Request.URL.Path, pageMeta["nextPageNumber"], pageMeta["requestedPageSize"])
	pageMeta["prevPageUrl"] = fmt.Sprintf("%s?page=%d&pageSize=%d", c.Request.URL.Path, pageMeta["prevPageNumber"], pageMeta["requestedPageSize"])

	c.JSON(statusCode, gin.H{
		"status":  statusCode,
		"success": true,
		"data": map[string]interface{}{
			dataName:   data,
			"pageMeta": pageMeta,
		},

		"message": message,
	})

}

func JSONNODATA(c *gin.Context, statusCode int, message interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  statusCode,
		"success": true,
		"data":    empty{},
		"message": message,
	})

}

func ERROR(c *gin.Context, statusCode int, messageError interface{}) {

	if reflect.TypeOf(messageError).String() == "string" {
		c.JSON(statusCode, gin.H{
			"status":  statusCode,
			"success": false,
			"data":    empty{},
			"message": messageError,
		})
		return
	}

	c.JSON(statusCode, gin.H{
		"status":  statusCode,
		"success": false,
		"data":    empty{},
		"message": "Validation Failed",
		"errors":  messageError,
	})

}
