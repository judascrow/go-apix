package controllers

import (
	"errors"
	"net/http"

	"github.com/judascrow/go-api-starter/api/middlewares/jwt"
	"github.com/judascrow/go-api-starter/api/models"

	"github.com/judascrow/go-api-starter/api/services"
	"github.com/judascrow/go-api-starter/api/utils/messages"
	"github.com/judascrow/go-api-starter/api/utils/responses"

	"github.com/gin-gonic/gin"
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
	err := c.BindJSON(&user)
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
	responses.JSON(c, http.StatusCreated, "user", user.Serialize(), "User "+messages.Created)
}

// @Summary แก้ไขผู้ใช้งาน
// @Description แก้ไขผู้ใช้งาน
// @Tags ผู้ใช้งาน
// @Accept  json
// @Produce  json
// @Param slug path string true "slug ผู้ใช้งาน"
// @Param user body models.SwagUserBody true "แก้ไขผู้ใช้งาน"
// @Success 201 {object} models.SwagUpdateUserResponse
// @Failure 400 {object} models.SwagError400
// @Failure 404 {object} models.SwagError404
// @Failure 500 {object} models.SwagError500
// @Security ApiKeyAuth
// @Router /users/{slug} [put]
func UpdateUser(c *gin.Context) {

	slug := c.Param("slug")

	var userData *models.User
	err := c.BindJSON(&userData)
	if err != nil {
		responses.ERROR(c, http.StatusBadRequest, messages.ErrorsResponse(err))
		return
	}

	var user models.User
	userData.Username = ""
	userData.Password = ""

	if user, err = services.UpdateUser(slug, userData); err != nil {
		responses.ERROR(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSON(c, http.StatusOK, "user", user.Serialize(), "User "+messages.Updated)

}

// @Summary ลบผู้ใช้งาน
// @Description ลบผู้ใช้งาน
// @Tags ผู้ใช้งาน
// @Accept  json
// @Produce  json
// @Param slug path string true "slug ผู้ใช้งาน"
// @Success 201 {object} models.SwagDeleteBase
// @Failure 400 {object} models.SwagError400
// @Failure 404 {object} models.SwagError404
// @Failure 500 {object} models.SwagError500
// @Security ApiKeyAuth
// @Router /users/{slug} [delete]
func DeleteUser(c *gin.Context) {
	slug := c.Param("slug")
	err := services.DeleteUser(&models.User{Slug: slug})
	if err != nil {
		responses.ERROR(c, http.StatusBadRequest, err.Error())
		return
	}
	responses.JSONNODATA(c, http.StatusOK, "User "+messages.Deleted)
}

// @Summary เปลี่ยนรหัสผ่าน
// @Description เปลี่ยนรหัสผ่าน
// @Tags ผู้ใช้งาน
// @Accept  json
// @Produce  json
// @Param slug path string true "slug ผู้ใช้งาน"
// @Param user body models.ChangePassword true "เปลี่ยนรหัสผ่าน"
// @Success 201 {object} models.SwagChangePasswordResponse
// @Failure 400 {object} models.SwagError400
// @Failure 404 {object} models.SwagError404
// @Failure 500 {object} models.SwagError500
// @Security ApiKeyAuth
// @Router /users/{slug}/password [put]
func ChangePassword(c *gin.Context) {

	slug := c.Param("slug")

	var requestBody models.ChangePassword
	err := c.BindJSON(&requestBody)
	if err != nil {
		responses.ERROR(c, http.StatusBadRequest, messages.ErrorsResponse(err))
		return
	}

	user, err := services.FindOneUserBySlug(slug)
	if err != nil {
		responses.ERROR(c, http.StatusNotFound, messages.NotFound)
		return
	}

	if !ClaimsOwner(c, slug) {
		responses.ERROR(c, http.StatusForbidden, messages.NotPermission)
		return
	}

	password := user.Password
	if requestBody.CurrentPassword == "" && requestBody.NewPassword != "" {
		err = errors.New("Please Provide current_password")
		responses.ERROR(c, http.StatusBadRequest, err.Error())
		return
	}
	if requestBody.CurrentPassword != "" && requestBody.NewPassword == "" {
		err = errors.New("Please Provide new_password")
		responses.ERROR(c, http.StatusBadRequest, err.Error())
		return
	}
	if requestBody.CurrentPassword != "" && requestBody.NewPassword != "" {
		//Also check if the new password
		if len(requestBody.NewPassword) < 6 {
			err = errors.New("Password should be atleast 6 characters")
			responses.ERROR(c, http.StatusBadRequest, err.Error())
			return
		}
		//if they do, check that the former password is correct
		err = verifyPassword(user.Password, requestBody.CurrentPassword)
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			err = errors.New("The password not correct")
			responses.ERROR(c, http.StatusBadRequest, err.Error())
			return
		}

		// Generate password
		bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			responses.ERROR(c, http.StatusBadRequest, err.Error())
			return
		}
		password = string(bcryptPassword)
	}

	userData := models.User{
		Password: password,
	}

	if user, err = services.UpdateUser(slug, userData); err != nil {
		responses.ERROR(c, http.StatusBadRequest, err.Error())
		return
	}
	responses.JSONNODATA(c, http.StatusOK, messages.ChangePasswordSuccess)

}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func ClaimsOwner(c *gin.Context, slug string) bool {

	claims := jwt.ExtractClaims(c)

	var roles = claims["roles"].([]interface{})
	for i := 0; i < len(roles); i++ {
		if uint(roles[i].(float64)) == 1 {
			return true
		}
	}

	if slug == claims["slug"].(string) || ClaimsIsAdmin(claims) {
		return true
	}
	return false
}

func ClaimsIsAdmin(claims jwt.MapClaims) bool {

	var roles = claims["roles"].([]interface{})
	for i := 0; i < len(roles); i++ {
		if uint(roles[i].(float64)) == 1 {
			return true
		}
	}

	return false
}
