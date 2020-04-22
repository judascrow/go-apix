package controllers

import (
	"net/http"
	"strconv"

	"github.com/judascrow/go-api-starter/api/models"

	"github.com/judascrow/go-api-starter/api/services"
	"github.com/judascrow/go-api-starter/api/utils/responses"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ListUsers(c *gin.Context) {
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
	serialized := make([]map[string]interface{}, length, length)
	for i := 0; i < length; i++ {
		serialized[i] = users[i].Serialize()
	}

	// Response
	responses.JSONLIST(c, http.StatusOK, serialized, pageMeta)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		errMessage := "invalid syntax"
		responses.ERROR(c, http.StatusBadRequest, errMessage)
		return
	}

	user, err := services.FindOneUserByID(uint(userID))
	if err != nil {
		responses.ERROR(c, http.StatusNotFound, err.Error())
		return
	}

	responses.JSON(c, http.StatusOK, user.Serialize(), nil)
}

func CreateUser(c *gin.Context) {

	var user models.User

	err := c.BindJSON(&user)
	if err != nil {
		responses.ERROR(c, http.StatusBadRequest, err.Error())
		return
	}

	condUsername := models.User{Username: user.Username}
	_, err = services.FindOneUser(condUsername)
	if err == nil {
		errMessage := "username is already exists"
		responses.ERROR(c, http.StatusNotFound, errMessage)
		return
	}
	condEmail := models.User{Email: user.Email}
	_, err = services.FindOneUser(condEmail)
	if err == nil {
		errMessage := "email is already exists"
		responses.ERROR(c, http.StatusNotFound, errMessage)
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(password)

	if err = services.CreateOne(&user); err != nil {
		responses.ERROR(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSON(c, http.StatusCreated, user.Serialize(), "user created success")
}
