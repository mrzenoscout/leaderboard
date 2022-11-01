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
		err := errors.New("invalid request body: player's name must be provided")
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
		log.Print(fmt.Errorf("get player's score by player name: %w", err))
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
	allTime, _ := strconv.ParseBool(c.Query("all-time"))

	// if limit is not set defaults to 10 results per page
	if limit == 0 {
		limit = 10
	}

	// if page number is not provided default to 0
	if page == 0 {
		page = 1
	}

	limit = limit + 1

	if allTime {
		// if month is provided but year is not set default to current year
		if month != 0 && year == 0 {
			year = int(time.Now().Year())
		}

		if month > 12 {
			month = 0
		}
	} else {
		month = int(time.Now().Month())
		year = int(time.Now().Year())
	}

	offset := (page - 1) * limit
	playersScores, err := store.GetPlayersScores(ctx, b.db, store.GetPlayersScoresOpts{
		Offset: offset,
		Limit:  limit,
		Month:  month,
		Year:   year,
	})
	if err != nil {
		log.Print(fmt.Errorf("get players scores: %w", err))
		c.Status(http.StatusInternalServerError)
		return
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
		Results: make([]Results, 0, len(playersScores)),
	}

	// check if requested page number is valid
	if len(playersScores) == 0 && page > 1 {
		c.Status(http.StatusNotFound)
		return
	} else if len(playersScores) == 0 {
		c.JSON(http.StatusOK, response)
		return
	}

	// calculate next page number
	if len(playersScores) == limit {
		playersScores = playersScores[:len(playersScores)-1]
		response.NextPage = page + 1
	}

	for _, playersScore := range playersScores {
		response.Results = append(response.Results, Results{
			Name:  playersScore.Player.Name,
			Score: playersScore.Score,
			Rank:  playersScore.Player.Rank,
		})
	}

	// if name is provided retrieve players around this player
	if name != "" && response.NextPage != 0 {
		playersScore, err := store.GetPlayersScoreByPlayerName(ctx, b.db, name)
		if err != nil {
			if err.Error() == "no rows in result set" {
				err = fmt.Errorf("player named '%s' not found", name)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			log.Print(fmt.Errorf("get players score by player's name: %w", err))
			c.Status(http.StatusInternalServerError)
			return
		}

		// if named player's rank is higher than last player's in the retrieved list
		// retrieve around player scores
		if playersScore.Player.Rank > playersScores[len(playersScores)-1].Player.Rank {
			var fromRank int
			if playersScore.Player.Rank > 2 {
				fromRank = playersScore.Player.Rank - 2
			}

			toRank := playersScore.Player.Rank + 2
			ps, err := store.GetPlayersScores(ctx, b.db, store.GetPlayersScoresOpts{
				FromRank: fromRank,
				ToRank:   toRank,
			})
			if err != nil {
				log.Print(fmt.Errorf("get scores around player: %w", err))
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
