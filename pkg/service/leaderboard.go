package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzenoscout/leaderboard/pkg/db"
	"gorm.io/gorm"
)

type BaseHandler struct {
	db *gorm.DB
}

func NewBaseHandler(db *gorm.DB) *BaseHandler {
	return &BaseHandler{
		db: db,
	}
}

func (b *BaseHandler) LeaderBoardRoutes(router *gin.Engine) {
	router.POST("/leaderboard/score", b.savePlayerScore)
	router.GET("/leaderboard", b.listPlayersScores)
}

func (b *BaseHandler) savePlayerScore(c *gin.Context) {
	var reqBody struct {
		Name  string `json:"name"`
		Score int    `json:"score"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		return
	}

	if reqBody.Name == "" {
		err := errors.New("player name must be provided")
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err)) //nolint:errcheck

		return
	}

	playersScore, err := db.StorePlayersScore(b.db, reqBody.Name, reqBody.Score)
	if err != nil {
		err = fmt.Errorf("store players score: %w", err)
		c.AbortWithError(http.StatusInternalServerError, err) //nolint:errcheck

		return
	}

	c.JSON(http.StatusOK, playersScore.Score)
}

func (b *BaseHandler) listPlayersScores(c *gin.Context) {
	playersScores, err := db.ListPlayersScores(b.db)
	if err != nil {
		err = fmt.Errorf("list players scores: %w", err)
		c.AbortWithError(http.StatusInternalServerError, err) //nolint:errcheck

		return
	}

	type results []struct {
		Name  string `json:"name"`
		Score int    `json:"score"`
		Rank  int    `json:"rank"`
	}

	r := make(results, len(playersScores))

	for i := range playersScores {
		r[i].Name = playersScores[i].Player.Name
		r[i].Score = playersScores[i].Score
		r[i].Rank = i + 1
	}

	c.JSON(http.StatusOK, r)
}
