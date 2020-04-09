package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApplyRoutes(router *gin.Engine) {

	apiv1 := router.Group("/api/v1")

	apiv1.GET("/healthcheck", healthcheck)

}

func healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "API is Online",
	})

}
