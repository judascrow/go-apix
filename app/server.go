package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/judascrow/go-api-crud/app/infrastructure"
	"github.com/judascrow/go-api-crud/app/models"
	"github.com/judascrow/go-api-crud/app/routes"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print(".env file not found")
	}
}

func migrate(db *gorm.DB) {

	db.AutoMigrate(&models.User{}) // UserSubscription must be migrated first otherwise the join table create has not the shape we are expecting

}

func Run() {

	var err error
	err = godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	database := infrastructure.InitDb()
	defer database.Close()

	migrate(database)

	gin.ForceConsoleColor()
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()

	routes.ApplyRoutes(app)

	port := os.Getenv("SERVER_PORT")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: app,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Print("API is ready to listen and serve on PORT : " + port)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")

}
