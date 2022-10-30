package model

import "time"

type Player struct {
	ID   int
	Name string
	Rank int
}

type PlayersScore struct {
	ID        int
	Score     int
	Player    *Player
	UpdatedAt time.Time
}
