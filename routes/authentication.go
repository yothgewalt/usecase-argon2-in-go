package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matthewhartstonge/argon2"
	"gorm.io/gorm"
)

var argon = argon2.DefaultConfig()

type Users struct {
	gorm.Model
	EmailAddress string `gorm:"not null;unique"`
	Username     string `gorm:"not null;unique"`
	Password     string `gorm:"not null"`
	Firstname    string `gorm:"not null;unique"`
	Lastname     string `gorm:"not null"`
}

type RegisterHandlerRequest struct {
	EmailAddress string `json:"email_address" db:"email_address"`
	Username     string `json:"username" db:"username"`
	Password     string `json:"password" db:"password"`
	Firstname    string `json:"firstname" db:"firstnname"`
	Lastname     string `json:"lastname" db:"lastname"`
}

type LoginHandlerRequest struct {
	EmailAddress string `json:"email_address" db:"email_address"`
	Password     string `json:"password" db:"password"`
}

func (r routes) addAuthentication(rg *gin.RouterGroup) {
	authetication := rg.Group("/auth")
	authetication.POST("/register", registerHandler)
	authetication.POST("/login", loginHandler)
}

func registerHandler(c *gin.Context) {
	bodyparser := RegisterHandlerRequest{}

	if err := c.BindJSON(&bodyparser); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tx := DB.Table("users").Select("email_address").Where("email_address = ?", bodyparser.EmailAddress).First(&bodyparser)
	if tx.RowsAffected == 1 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "We found that someone else is already using this email.",
		})
		return
	}

	tx = DB.Table("users").Select("username").Where("username = ?", bodyparser.Username).First(&bodyparser)
	if tx.RowsAffected == 1 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "We found that someone else is already using this username.",
		})
		return
	}

	tx = DB.Table("users").Select("firstname").Where("firstname = ?", bodyparser.Firstname).First(&bodyparser)
	if tx.RowsAffected == 1 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "We found that someone else is already using this firstname.",
		})
		return
	}

	encoded, err := argon.HashEncoded([]byte(bodyparser.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	payload := Users{
		EmailAddress: bodyparser.EmailAddress,
		Username:     bodyparser.Username,
		Password:     string(encoded),
		Firstname:    bodyparser.Firstname,
		Lastname:     bodyparser.Lastname,
	}

	tx = DB.Table("users").Create(&payload)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, tx.Error)
		return
	}

	c.JSON(http.StatusAccepted, &bodyparser)
}

func loginHandler(c *gin.Context) {
	bodyparser := LoginHandlerRequest{}
	if err := c.BindJSON(&bodyparser); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	payload := LoginHandlerRequest{
		EmailAddress: bodyparser.EmailAddress,
		Password:     bodyparser.Password,
	}

	tx := DB.Table("users").Select("email_address", "password").Where("email_address = ?", bodyparser.EmailAddress).First(&bodyparser)
	if tx.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "We couldn't find your account.",
		})
		return
	}

	ok, err := argon2.VerifyEncoded([]byte(payload.Password), []byte(bodyparser.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "We are unable to compare and decode at this time.",
			"error":   err,
		})
		return
	}

	if !ok {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Invalid password, please try again.",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "You have successfully logged in.",
	})
}
