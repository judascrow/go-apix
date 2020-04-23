package routes

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/judascrow/go-api-starter/api/controllers"
	"github.com/judascrow/gomiddlewares"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func InitRouter() *gin.Engine {

	r := gin.New()

	// Prometheus
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
	p.ReqCntURLLabelMappingFn = MappingFn

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

	// swagger
	apiv1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1.GET("/healthcheck", Healthcheck)

	authMiddleware := AuthMiddlewareJWT()
	apiv1.POST("/login", authMiddleware.LoginHandler)

	users := apiv1.Group("/users")
	users.Use(authMiddleware.MiddlewareFunc())
	{
		users.GET("", controllers.GetAllUsers)
		users.GET("/:slug", controllers.GetUserBySlug)
		users.POST("", controllers.CreateUser)
	}

	return r
}

func Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "API is Online",
	})

}

func MappingFn(c *gin.Context) string {
	url := c.Request.URL.Path
	for _, p := range c.Params {
		if p.Key == "id" {
			url = strings.Replace(url, p.Value, ":id", 1)
			break
		}
	}
	return url
}
