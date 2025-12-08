package users

import (
	"game-ranker/users-manager/internal"
	"game-ranker/users-manager/internal/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User = internal.User

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

	if len(body.Password) < 8 {
		c.String(http.StatusBadRequest, "Invalid password")
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
		ID:         uuid.NewString(),
	}

	database.AddUser(user)
}
