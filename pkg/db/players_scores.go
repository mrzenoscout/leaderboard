package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/mrzenoscout/leaderboard/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func StorePlayersScore(db *gorm.DB, name string, score int) (*model.PlayersScore, error) {
	player := model.Player{Name: name}

	err := db.Where("name = ?", player.Name).First(&player).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("select player: %w", err)
	}

	if player.ID == 0 {
		if err := db.Create(&player).Error; err != nil {
			return nil, fmt.Errorf("create player: %w", err)
		}
	}

	playersScore := model.PlayersScore{
		Player:    player,
		Score:     score,
		UpdatedAt: time.Now(),
	}

	if err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "player_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"score":      gorm.Expr(fmt.Sprintf("GREATEST(players_scores.score, %d)", playersScore.Score)),
			"updated_at": playersScore.UpdatedAt,
		}),
	}).Create(&playersScore).Error; err != nil {
		return nil, fmt.Errorf("upsert players score: %w", err)
	}

	err = db.Where("player_id = ?", player.ID).First(&playersScore).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("select players score: %w", err)
	}

	return &playersScore, nil
}

func ListPlayersScores(db *gorm.DB) ([]model.PlayersScore, error) {
	var playersScores []model.PlayersScore

	err := db.Joins("Player").Find(&playersScores).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("select players scores: %w", err)
	}

	return playersScores, nil
}
