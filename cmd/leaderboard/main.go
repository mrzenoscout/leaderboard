package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/mrzenoscout/leaderboard/pkg/db"
	"github.com/mrzenoscout/leaderboard/pkg/service"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("load .env file")
	}
}

func main() {
	router := gin.New()

	gormDB, err := db.Connect()
	if err != nil {
		log.Fatalf("connect to db: %s", err)
	}

	if err = db.AutoMigrate(gormDB); err != nil {
		log.Fatal("auto migrate db")
	}

	service.NewBaseHandler(gormDB).LeaderBoardRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("router run: %s", err)
	}
}
