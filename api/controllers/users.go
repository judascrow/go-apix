package controllers

import (
	"net/http"

	"github.com/judascrow/go-api-starter/api/models"

	"github.com/judascrow/go-api-starter/api/services"
	"github.com/judascrow/go-api-starter/api/utils/messages"
	"github.com/judascrow/go-api-starter/api/utils/responses"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(c *gin.Context) {
	// Query Pages
	pageSizeStr := c.Query("pageSize")
	pageStr := c.Query("page")

	// Find Users
	users, pageMeta, err := services.FindAllUsers(pageSizeStr, pageStr)
	if err != nil {
		responses.ERROR(c, http.StatusNotFound, err.Error())
	}

	// Serialize
	length := len(users)
	UserSerialized := make([]map[string]interface{}, length, length)
	for i := 0; i < length; i++ {
		UserSerialized[i] = users[i].Serialize()
	}

	// Response
	responses.JSONLIST(c, http.StatusOK, "users", UserSerialized, messages.DataFound, pageMeta)
}

func GetUserBySlug(c *gin.Context) {
	// Get Slug from URI
	slug := c.Param("slug")

	// Find User
	user, err := services.FindOneUserBySlug(slug)
	if err != nil {
		responses.ERROR(c, http.StatusNotFound, err.Error())
		return
	}

	// Response
	responses.JSON(c, http.StatusOK, "user", user.Serialize(), messages.DataFound)
}

func CreateUser(c *gin.Context) {

	// Define struct user variable
	var user models.User

	// Map jsonBody to struct
	err := c.BindJSON(&user)
	if err != nil {
		responses.ERROR(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check username duplicate
	userCond := models.User{Username: user.Username}
	_, err = services.FindOneUser(userCond)
	if err == nil {
		errMessage := "username " + messages.IsAlreadyExists
		responses.ERROR(c, http.StatusBadRequest, errMessage)
		return
	}
	// Check email duplicate
	userCond = models.User{Email: user.Email}
	_, err = services.FindOneUser(userCond)
	if err == nil {
		errMessage := "email " + messages.IsAlreadyExists
		responses.ERROR(c, http.StatusBadRequest, errMessage)
		return
	}

	// Generate password
	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(password)

	// Create user
	if err = services.CreateOne(&user); err != nil {
		responses.ERROR(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	responses.JSON(c, http.StatusCreated, "user", user.Serialize(), "user "+messages.HasBeenCreated)
}
