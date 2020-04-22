package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/judascrow/gomiddlewares"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/judascrow/go-api-starter/api/infrastructure"
	"github.com/judascrow/go-api-starter/api/models"
	"github.com/judascrow/go-api-starter/api/routes"
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

	gin.SetMode(os.Getenv("SERVER_RUN_MODE"))

	r := routes.InitRouter()
	defer gomiddlewares.CloseLogFile()

	port := os.Getenv("SERVER_PORT")
	readTimeoutInt, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	writeTimeoutInt, _ := strconv.Atoi(os.Getenv("SERVER_WRITE_TIMEOUT"))
	readTimeout := time.Duration(readTimeoutInt) * time.Second
	writeTimeout := time.Duration(writeTimeoutInt) * time.Second
	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
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
