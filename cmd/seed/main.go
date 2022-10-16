package main

import (
	"log"

	"github.com/brianvoe/gofakeit"
	"github.com/joho/godotenv"
	"github.com/mrzenoscout/leaderboard/pkg/db"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("load .env file")
	}
}

func main() {
	gormDB, err := db.Connect()
	if err != nil {
		log.Fatalf("connect to db: %s", err)
	}

	for i := 0; i < 1000; i++ {
		if _, err := db.StorePlayersScore(gormDB, gofakeit.Name(), gofakeit.Number(0, 1000000)); err != nil {
			log.Fatalf("store players score: %s", err)
		}
	}
}
