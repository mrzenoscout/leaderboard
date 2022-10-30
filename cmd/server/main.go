package main

import (
	"log"
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	"github.com/mrzenoscout/leaderboard/internal/core/drivers/psql"
	"github.com/mrzenoscout/leaderboard/internal/player/store"
	"github.com/mrzenoscout/leaderboard/internal/tansport/http"
)

// func init() {
// 	if err := godotenv.Load(".env"); err != nil {
// 		log.Fatal("load .env file")
// 	}
// }

func main() {
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	http.NewBaseHandler(db)

	if err := psql.MigratePostgres(ctx, "file://migrations"); err != nil {
		log.Fatalf("migrate postgres: %s", err)
	}

	http.NewBaseHandler(conn).LeaderBoardRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("router run: %s", err)
	}
}
