package main

import (
	"game-ranker/users-manager/internal/database"
	"game-ranker/users-manager/internal/users"

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

	r.Run()
}
