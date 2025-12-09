package main

import (
	"game-ranker/users-manager/internal/database"
	"game-ranker/users-manager/internal/users"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.InitDbTable()
	r := gin.Default()

	r.POST("/register", func(ctx *gin.Context) {
		users.RegisterAccount(ctx)
	})

	auth := r.Group("/auth")
	auth.Use(users.AuthMiddleware())

	auth.GET("/me", func(ctx *gin.Context) {
		email := ctx.GetString("user_email")
		username := ctx.GetString("username")
		role := ctx.GetString("role")

		ctx.JSON(http.StatusOK, gin.H{
			"email":    email,
			"username": username,
			"role":     role,
		})
	})

	r.Run()
}
