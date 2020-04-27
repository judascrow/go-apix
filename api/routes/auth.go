package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/judascrow/go-api-starter/api/infrastructure"
	"github.com/judascrow/go-api-starter/api/middlewares/jwt"
	"github.com/judascrow/go-api-starter/api/models"
	"github.com/judascrow/go-api-starter/api/utils/messages"

	"golang.org/x/crypto/bcrypt"

	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var casbinEnforcer *casbin.Enforcer
var db *gorm.DB
var err error

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// User struct alias
type User = models.User

var identityKey = "slug"
var identityUsername = "username"
var identityRoles = "roles"

func AuthMiddlewareJWT() *jwt.GinJWTMiddleware {
	db := infrastructure.GetDB()

	Dbdriver := os.Getenv("DB_DRIVER")
	DbName := os.Getenv("DB_NAME")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbHost := os.Getenv("DB_HOST")
	DbPort := os.Getenv("DB_PORT")

	pg_conn_info := DbUser + ":" + DbPassword + "@tcp(" + DbHost + ":" + DbPort + ")/" + DbName + "?charset=utf8&parseTime=True&loc=Local"
	casbin_adapter := gormadapter.NewAdapter(Dbdriver, pg_conn_info, true)
	e := casbin.NewEnforcer("./auth.conf", casbin_adapter)
	casbinEnforcer = e
	e.LoadPolicy()

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				var roles []uint
				for _, r := range v.UserRoles {
					roles = append(roles, r.RoleID)
				}
				return jwt.MapClaims{
					identityKey:      v.Slug,
					identityUsername: v.Username,
					identityRoles:    roles,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			var roles = claims["roles"].([]interface{})
			userRoles := make([]models.UserRole, len(roles))
			for i := 0; i < len(roles); i++ {
				userRoles[i].RoleID = uint(roles[i].(float64))
			}
			return &User{
				Slug:      claims["slug"].(string),
				Username:  claims["username"].(string),
				UserRoles: userRoles,
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.BindJSON(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			var user User
			if err := db.Preload("UserRoles").Where("username = ? AND status = 'A' ", username).First(&user).Error; err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			if checkHash(password, user.Password) {
				return &User{
					Slug:      user.Slug,
					Username:  user.Username,
					UserRoles: user.UserRoles,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			//fmt.Println(len(data.(*User).UserRoles))
			if v, ok := data.(*User); ok && v.Username != "" && len(v.UserRoles) > 0 {
				v0 := ""
				for _, role := range v.UserRoles {
					v0 = strconv.Itoa(int(role.RoleID))
					fmt.Println(v0)
					return casbinEnforcer.Enforce(v0, c.Request.URL.String(), c.Request.Method)
				}
				return casbinEnforcer.Enforce(v0, c.Request.URL.String(), c.Request.Method)

			}

			return false
		},
		LoginResponse: func(c *gin.Context, code int, token string, t time.Time, claims map[string]interface{}) {
			var user models.User
			slug := claims[identityKey]
			if slug != "" {
				db := infrastructure.GetDB()
				db.Preload("Roles").Where("slug = ?", slug).First(&user)
			}
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"token":  token,
				// "expire":  t.Format(time.RFC3339),
				"success": true,
				"message": messages.Logged,
				"data":    user.Serialize(),
			})
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"success": false,
				"status":  code,
				"message": message,
				"data":    map[string]interface{}{},
			})
		},

		TokenLookup: "header: Authorization, query: token, cookie: jwt",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}

func checkHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
