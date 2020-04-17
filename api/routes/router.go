package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/judascrow/go-api-crud/api/middlewares"
)

func InitRouter() *gin.Engine {

	r := gin.New()
	r.HandleMethodNotAllowed = true
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"success": false, "error": "Method Not Allowed"})
		return
	})

	r.Use(middlewares.CustomLoggerZap(), middlewares.CorsMiddleware())
	if os.Getenv("APP_ENV") == "dev" {
		gin.ForceConsoleColor()
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())

	apiv1 := r.Group(os.Getenv("APP_API_BASE_URL"))

	apiv1.GET("/healthcheck", healthcheck)

	return r
}

func healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "API is Online",
	})

}
