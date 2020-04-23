package controllers

import (
	"net/http"

	"github.com/judascrow/go-api-starter/api/models"

	"github.com/judascrow/go-api-starter/api/services"
	"github.com/judascrow/go-api-starter/api/utils/messages"
	"github.com/judascrow/go-api-starter/api/utils/responses"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
)

// @Summary รายการผู้ใช้งาน
// @Description รายการผู้ใช้งาน
// @Tags ผู้ใช้งาน
// @Accept  json
// @Produce  json
// @Success 200 {object} models.SwagGetAllUsersResponse
// @Failure 400 {object} models.SwagError400
// @Failure 404 {object} models.SwagError404
// @Failure 500 {object} models.SwagError500
// @Router /users [get]
func GetAllUsers(c *gin.Context) {
	// Query Pages
	pageSizeStr := c.Query("pageSize")
	pageStr := c.Query("page")

	// Find Users
	users, pageMeta, err := services.FindAllUsers(pageSizeStr, pageStr)
	if err != nil {
		responses.ERROR(c, http.StatusNotFound, messages.NotFound)
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

// @Summary ข้อมูลผู้ใช้งาน
// @Description ข้อมูลผู้ใช้งาน
// @Tags ผู้ใช้งาน
// @Accept  json
// @Produce  json
// @Param slug path string true "slug ผู้ใช้งาน"
// @Success 200 {object} models.SwagGetUserResponse
// @Failure 400 {object} models.SwagError400
// @Failure 404 {object} models.SwagError404
// @Failure 500 {object} models.SwagError500
// @Router /users/{slug} [get]
func GetUserBySlug(c *gin.Context) {
	// Get Slug from URI
	slug := c.Param("slug")

	// Find User
	user, err := services.FindOneUserBySlug(slug)
	if err != nil {
		responses.ERROR(c, http.StatusNotFound, messages.NotFound)
		return
	}

	// Response
	responses.JSON(c, http.StatusOK, "user", user.Serialize(), messages.DataFound)
}

// @Summary เพิ่มผู้ใช้งาน
// @Description เพิ่มผู้ใช้งาน
// @Tags ผู้ใช้งาน
// @Accept  json
// @Produce  json
// @Param user body models.SwagUserBodyIncludePassword true "เพิ่มผู้ใช้งาน"
// @Success 201 {object} models.SwagCreateUserResponse
// @Failure 400 {object} models.SwagError400
// @Failure 404 {object} models.SwagError404
// @Failure 500 {object} models.SwagError500
// @Security ApiKeyAuth
// @Router /users [post]
func CreateUser(c *gin.Context) {

	// Define struct user variable
	var user models.User

	// Map jsonBody to struct
	err := c.ShouldBindWith(&user, binding.JSON)
	if err != nil {
		responses.ERROR(c, http.StatusBadRequest, messages.ErrorsResponse(err))
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
	responses.JSON(c, http.StatusCreated, "user", user.Serialize(), "user "+messages.Created)
}
