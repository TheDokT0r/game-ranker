package users

import (
	"game-ranker/users-manager/internal"
	"game-ranker/users-manager/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User = internal.User

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type JwtClaims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`

	jwt.RegisteredClaims
}

func RegisterAccount(c *gin.Context) {
	var body RegisterRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(body.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Password"})
		return
	}

	hashedPassword, err := hashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := User{
		Username:   body.Username,
		HashedPass: hashedPassword,
		Email:      body.Email,
		ID:         uuid.NewString(),
		Role:       "user",
	}

	err = database.AddUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	signedToken := CreateJwtSingedToken(user)

	c.JSON(http.StatusOK, gin.H{
		"jwt": signedToken,
	})
}

type LoginInfo struct {
	Email    string
	Password string
}

func Login(c *gin.Context) {
	var body LoginInfo
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := database.GetUser(body.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid information"})
		return
	}

	err = checkPasswordHash(user.HashedPass, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid information"})
		return
	}

	signedToken := CreateJwtSingedToken(*user)

	c.JSON(http.StatusOK, gin.H{"jwt": signedToken})
}

func hashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}

func checkPasswordHash(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
