package users

import (
	"game-ranker/users-manager/internal"
	"game-ranker/users-manager/internal/database"
	"log"
	"net/http"
	"os"
	"time"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	user := User{
		Username:   body.Username,
		HashedPass: hashedPassword,
		Email:      body.Email,
		ID:         uuid.NewString(),
	}

	err = database.AddUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	claims := JwtClaims{
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 1, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret, present := os.LookupEnv("SECRET")
	if !present {
		log.Fatal("Make sure you've set all of your environment variables")
	}

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"jwt": signed,
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

}

func hashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}

func checkPasswordHash(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
