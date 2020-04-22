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
		responses.ERROR(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := services.FindOneUser(uint(userID))
	if err != nil {
		responses.ERROR(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSON(c, http.StatusOK, user.Serialize())
}

func CreateUser(c *gin.Context) {

	var user models.User

	err := c.BindJSON(&user)
	if err != nil {
		responses.ERROR(c, http.StatusBadRequest, err.Error())
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(password)

	if err = services.CreateOne(&user); err != nil {
		responses.ERROR(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSON(c, http.StatusCreated, user.Serialize())
}
