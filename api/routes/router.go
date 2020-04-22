package routes

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/judascrow/gomiddlewares"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func InitRouter() *gin.Engine {

	r := gin.New()

	// Prometheus
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path
		for _, p := range c.Params {
			if p.Key == "id" {
				url = strings.Replace(url, p.Value, ":id", 1)
				break
			}
		}
		return url
	}

	r.HandleMethodNotAllowed = true
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"success": false, "error": "Method Not Allowed"})
		return
	})

	// middlewares
	r.Use(gomiddlewares.GoLogger(), gomiddlewares.GoCors())
	if os.Getenv("APP_ENV") == "dev" {
		gin.ForceConsoleColor()
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())

	// Routes
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
