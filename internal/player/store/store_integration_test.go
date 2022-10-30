//go:build integration
// +build integration

package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5"
	"github.com/mrzenoscout/leaderboard/internal/core/drivers/psql"
	"github.com/mrzenoscout/leaderboard/internal/player/model"
	"github.com/stretchr/testify/assert"
)

var (
	rootCtx = context.Background()
	db      *pgx.Conn
)

func TestMain(m *testing.M) {
	var err error

	db, err = psql.Connect(rootCtx)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestGetPlayerByName(t *testing.T) {
	player, err := GetPlayerByName(rootCtx, db, gofakeit.BeerStyle())

	assert.NoError(t, err)
	assert.Empty(t, player)
}

func TestGetPlayersScoreByPlayerID(t *testing.T) {
	t.Parallel()

	playerName := gofakeit.Name()

	player, err := GetPlayerByName(rootCtx, db, playerName)
	assert.NoError(t, err)

	if player.ID == 0 {
		player, err = InsertPlayer(rootCtx, db, &model.Player{Name: playerName})
		assert.NoError(t, err)
	}

	assert.Equal(t, playerName, player.Name)

	playersScore, err := UpsertPlayersScore(rootCtx,
		db, &model.PlayersScore{Player: player, Score: gofakeit.Number(0, 100000)})
	assert.NoError(t, err)

	if playersScore.ID == 0 {
		playersScore, err = GetPlayersScoreByPlayerName(rootCtx, db, player.Name)
		assert.NoError(t, err)
	}

	assert.NotEmpty(t, playersScore.ID)
	assert.NotZero(t, playersScore.Player.Rank)
}

func TestGetPlayersScores(t *testing.T) {
	t.Parallel()

	for i := 0; i < 10; i++ {
		createPlayersScore(t)
	}

	scores, err := GetPlayersScores(rootCtx, db, GetPlayersScoresOpts{
		Limit:  3,
		Offset: 3,
		Month:  int(time.Now().Month()),
		Year:   int(time.Now().Year()),
	})

	assert.NoError(t, err)
	assert.Len(t, scores, 3)
}

func TestGetPlayersScoresAroundRank(t *testing.T) {
	t.Parallel()

	for i := 0; i < 10; i++ {
		createPlayersScore(t)
	}

	scores, err := GetPlayersScores(rootCtx, db, GetPlayersScoresOpts{
		FromRank: 1,
		ToRank:   3,
	})

	assert.NoError(t, err)
	assert.Len(t, scores, 3)
}

// helper funcs

func createPlayersScore(t *testing.T) {
	t.Helper()

	playerName := gofakeit.Name()

	player, err := GetPlayerByName(rootCtx, db, playerName)
	assert.NoError(t, err)

	if player.ID == 0 {
		player, err = InsertPlayer(rootCtx, db, &model.Player{Name: playerName})
		assert.NoError(t, err)
	}

	assert.Equal(t, playerName, player.Name)

	playersScore, err := UpsertPlayersScore(rootCtx,
		db, &model.PlayersScore{Player: player, Score: gofakeit.Number(0, 100000)})
	assert.NoError(t, err)

	if playersScore.ID == 0 {
		playersScore, err = GetPlayersScoreByPlayerName(rootCtx, db, player.Name)
		assert.NoError(t, err)
	}

	assert.NotEmpty(t, playersScore.ID)
}
