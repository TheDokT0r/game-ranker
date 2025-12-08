package users

import (
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         string
	Username   string
	Email      string
	HashedPass string
}

type RegisterRequest struct {
	Username string
	Password string
	Email    string
}

func RegisterAccount(c *gin.Context) {
	var body RegisterRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	pBytes, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	if err != nil {
		log.Fatal(err)
	}

	hashedPassword := string(pBytes)

	user := User{
		Username:   body.Username,
		HashedPass: hashedPassword,
		Email:      body.Email,
	}
}
