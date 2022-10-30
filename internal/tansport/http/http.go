package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type BaseHandler struct {
	db *pgx.Conn
}

func NewBaseHandler(db *pgx.Conn) *BaseHandler {
	return &BaseHandler{
		db: db,
	}
}

func (b *BaseHandler) LeaderBoardRoutes(router *gin.Engine) {
	router.POST("/leaderboard/score", b.savePlayerScore)
	router.GET("/leaderboard", b.listPlayersScores)
}
