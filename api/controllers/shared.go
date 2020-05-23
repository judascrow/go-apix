package controllers

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/judascrow/gomiddlewares/jwt"
	"golang.org/x/crypto/bcrypt"
)

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
