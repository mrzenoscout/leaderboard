package model

import "time"

type Player struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"index:idx_name,unique; not null"`
	Rank int    `gorm:"-" json:"rank"`
}

type PlayersScore struct {
	ID        int `gorm:"primaryKey"`
	Score     int `gorm:"not null" json:"score"`
	PlayerID  int `gorm:"index:idx_player,unique; not null"`
	Player    Player
	UpdatedAt time.Time `gorm:"not null"`
}
