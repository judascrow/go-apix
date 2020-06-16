package seeds

import (
	"math/rand"
	"time"

	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	"github.com/judascrow/go-apix/api/infrastructure"
	"github.com/judascrow/go-apix/api/models"
	"golang.org/x/crypto/bcrypt"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var StringNumberRunes = []rune("1234567890")

func randomInt(min, max int) int {

	return rand.Intn(max-min) + min
}

func randomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func randomStringNumber(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = StringNumberRunes[rand.Intn(len(StringNumberRunes))]
	}
	return string(b)
}

func seedAdmin(db *gorm.DB) {
	count := 0
	adminRole := models.Role{Name: "ROLE_ADMIN", Description: "Only for admin"}
	query := db.Model(&models.Role{}).Where("name = ?", "ROLE_ADMIN")
	query.Count(&count)

	if count == 0 {
		db.Create(&adminRole)
	} else {
		query.First(&adminRole)
	}

	adminRoleUsers := 0
	var adminUsers []models.User
	db.Model(&adminRole).Related(&adminUsers, "Users")

	db.Model(&models.User{}).Where("username = ?", "admin").Count(&adminRoleUsers)
	if adminRoleUsers == 0 {

		// query.First(&adminRole) // First would fetch the Role admin because the query status name='ROLE_ADMIN'
		password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		// Approach 1
		user := models.User{FirstName: "AdminFN", LastName: "AdminLN", Email: "admin@golang.com", Username: "admin", Password: string(password)}
		user.Roles = append(user.Roles, adminRole)

		// Do not try to update the adminRole
		db.Set("gorm:association_autoupdate", false).Create(&user)

		if db.Error != nil {
			print(db.Error)
		}
	}
}

func seedUsers(db *gorm.DB) {
	count := 0
	role := models.Role{Name: "ROLE_USER", Description: "Only for standard users"}
	q := db.Model(&models.Role{}).Where("name = ?", "ROLE_USER")
	q.Count(&count)

	if count == 0 {
		db.Create(&role)
	} else {
		q.First(&role)
	}

	var standardUsers []models.User
	db.Model(&role).Related(&standardUsers, "Users")
	usersCount := len(standardUsers)
	usersToSeed := 5
	usersToSeed -= usersCount
	if usersToSeed > 0 {
		for i := 0; i < usersToSeed; i++ {
			password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
			user := models.User{FirstName: fake.FirstName(), LastName: fake.LastName(), Email: fake.EmailAddress(), Username: fake.UserName(),
				Password: string(password)}
			// No need to add the role as we did for seedAdmin, it is added by the BeforeSave hook
			db.Set("gorm:association_autoupdate", false).Create(&user)
		}
	}
}

func seedCasbinRule(db *gorm.DB) {
	var casbinRule [2]models.CasbinRule

	db.Where(&models.CasbinRule{PType: "p", V0: "1", V1: "/api/v1/*"}).Attrs(models.CasbinRule{V2: "(GET)|(POST)|(PUT)|(DELETE)"}).FirstOrCreate(&casbinRule[0])
	db.Where(&models.CasbinRule{PType: "p", V0: "2", V1: "/api/v1/users/*"}).Attrs(models.CasbinRule{V2: "(GET)|(PUT)"}).FirstOrCreate(&casbinRule[1])

}

func Seed() {
	db := infrastructure.GetDB()
	rand.Seed(time.Now().UnixNano())
	seedAdmin(db)
	seedUsers(db)
	seedCasbinRule(db)
}
