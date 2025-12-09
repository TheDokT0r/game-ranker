package users

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJwtSingedToken(user User) string {
	claims := JwtClaims{
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 1, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, present := os.LookupEnv("SECRET")
	if !present {
		log.Fatal("Make sure you've set all of your enviorment variables")
	}

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Panic(err)
	}

	return signed
}

func getSecret() string {
	secret, present := os.LookupEnv("SECRET")
	if !present {
		log.Fatal("Make sure you've set all of your enviorment variables")
	}

	return secret
}
