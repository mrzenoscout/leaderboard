package main

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/joho/godotenv"
	"github.com/mrzenoscout/leaderboard/internal/core/drivers/psql"
	"github.com/mrzenoscout/leaderboard/internal/player/model"
	"github.com/mrzenoscout/leaderboard/internal/player/store"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("load .env file")
	}
}

func main() {
	ctx := context.Background()

	db, err := psql.Connect(ctx)
	if err != nil {
		log.Fatalf("connect to db: %s", err)
	}

	for i := 0; i < 10000; i++ {
		playersScore := model.PlayersScore{
			Score:     gofakeit.Number(0, 10000000),
			UpdatedAt: gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
		}

		playersScore.Player, err = store.InsertPlayer(ctx, db, &model.Player{Name: gofakeit.Name()})
		if err != nil {
			if psql.IsErrorCode(err, "23505") {
				continue
			}

			log.Fatalf("insert player: %s", err)
		}

		_, err := store.UpsertPlayersScore(ctx, db, &playersScore)
		if err != nil {
			log.Fatalf("insert players score: %s", err)
		}
	}
}
