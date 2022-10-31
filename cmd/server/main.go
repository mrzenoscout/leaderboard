package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mrzenoscout/leaderboard/internal/core/drivers/psql"
	"github.com/mrzenoscout/leaderboard/internal/tansport/http"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("load .env file")
	}
}

func main() {
	ctx := context.Background()
	router := gin.Default()

	db, err := psql.Connect(ctx)
	if err != nil {
		log.Fatalf("connect to db: %s", err)
	}

	if err := psql.MigratePostgres(ctx, "file://migrations"); err != nil {
		log.Fatalf("migrate postgres: %s", err)
	}

	http.NewBaseHandler(db).LeaderBoardRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("router run: %s", err)
	}
}
