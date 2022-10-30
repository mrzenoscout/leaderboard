package http

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrzenoscout/leaderboard/internal/player/model"
	"github.com/mrzenoscout/leaderboard/internal/player/store"
)

func (b *BaseHandler) savePlayerScore(c *gin.Context) {
	ctx := c.Request.Context()

	var request struct {
		Name  string `json:"name"`
		Score int    `json:"score"`
	}

	if err := c.BindJSON(&request); err != nil {
		err = fmt.Errorf("invalid request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Name == "" {
		err := errors.New("invalid request body: player name must be provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := store.GetPlayerByName(ctx, b.db, request.Name)
	if err != nil {
		log.Print(fmt.Errorf("get player by name: %w", err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if player.ID == 0 {
		player, err = store.InsertPlayer(ctx, b.db, &model.Player{Name: request.Name})
		if err != nil {
			log.Print(fmt.Errorf("insert player: %w", err))
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	if _, err := store.UpsertPlayersScore(ctx, b.db,
		&model.PlayersScore{
			Player:    player,
			Score:     request.Score,
			UpdatedAt: time.Now(),
		},
	); err != nil {
		log.Print(fmt.Errorf("upsert players score: %w", err))
		c.Status(http.StatusInternalServerError)
		return
	}

	playersScore, err := store.GetPlayersScoreByPlayerName(ctx, b.db, player.Name)
	if err != nil {
		log.Print(fmt.Errorf("get players score by player name: %w", err))
		c.Status(http.StatusInternalServerError)
		return
	}

	type Response struct {
		Rank int `json:"rank"`
	}

	c.JSON(http.StatusOK, Response{
		Rank: playersScore.Player.Rank,
	})
}

func (b *BaseHandler) listPlayersScores(c *gin.Context) {
	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	name := c.Query("name")
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))

	// if limit is not set defaults to 10 results per page
	if limit == 0 {
		limit = 10
	}

	// if page number is not provided default to 0
	if page == 0 {
		page = 1
	}

	// if month is provided but year is not set defaults to current year
	if month != 0 && year == 0 {
		year = int(time.Now().Year())
	}

	playersScores, err := store.GetPlayersScores(ctx, b.db, store.GetPlayersScoresOpts{
		Offset: (page - 1) * limit,
		// retrieve additional record to check if page number can be increased
		Limit: limit + 1,
		Month: month,
		Year:  year,
	})
	if err != nil {
		log.Print(fmt.Errorf("get players score: %w", err))
		c.Status(http.StatusInternalServerError)
		return
	}

	// check if requested page number is valid
	if len(playersScores) == 0 && page > 1 {
		err = errors.New("requested page doesn't hold any records")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else if len(playersScores) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	// calculate next page number
	var nextPage int
	if len(playersScores) == limit+1 {
		playersScores = playersScores[:len(playersScores)-1]
		nextPage = page + 1
	}

	type Results struct {
		Name  string `json:"name"`
		Score int    `json:"score"`
		Rank  int    `json:"rank"`
	}

	type Response struct {
		Results  []Results `json:"results"`
		AroundMe []Results `json:"around_me,omitempty"`
		NextPage int       `json:"next_page"`
	}

	response := Response{
		Results:  make([]Results, 0, len(playersScores)-1),
		NextPage: nextPage,
	}

	for _, playersScore := range playersScores {
		response.Results = append(response.Results, Results{
			Name:  playersScore.Player.Name,
			Score: playersScore.Score,
			Rank:  playersScore.Player.Rank,
		})
	}

	// if name is provided retrieve players around this player
	if name != "" && nextPage != 0 {
		playersScore, err := store.GetPlayersScoreByPlayerName(ctx, b.db, name)
		if err != nil {
			log.Print(fmt.Errorf("get players score by player name: %w", err))
			c.Status(http.StatusInternalServerError)
			return
		}

		// if named player's rank is higher thank last player's in the retrieved list
		// retrieve around player scores
		if playersScore.Player.Rank > playersScores[len(playersScores)-1].Player.Rank {
			var fromRank int
			if playersScore.Player.Rank > 2 {
				fromRank = playersScore.Player.Rank - 2
			}

			ps, err := store.GetPlayersScores(ctx, b.db, store.GetPlayersScoresOpts{
				FromRank: fromRank,
				ToRank:   playersScore.Player.Rank + 2,
			})
			if err != nil {
				log.Print(fmt.Errorf("get players scores: %w", err))
				c.Status(http.StatusInternalServerError)
				return
			}

			for _, playersScore := range ps {
				response.AroundMe = append(response.AroundMe, Results{
					Name:  playersScore.Player.Name,
					Score: playersScore.Score,
					Rank:  playersScore.Player.Rank,
				})
			}
		}
	}

	c.JSON(http.StatusOK, response)
}
