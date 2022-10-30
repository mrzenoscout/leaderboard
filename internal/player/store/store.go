package store

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/mrzenoscout/leaderboard/internal/player/model"
)

func InsertPlayer(ctx context.Context, db *pgx.Conn, player *model.Player) (*model.Player, error) {
	if err := db.QueryRow(ctx, `INSERT INTO players (name) VALUES ($1) RETURNING id`, player.Name).
		Scan(&player.ID); err != nil {
		return nil, fmt.Errorf("insert player: %w", err)
	}

	return player, nil
}

func GetPlayerByName(ctx context.Context, db *pgx.Conn, name string) (*model.Player, error) {
	var player model.Player

	query := `SELECT id, name FROM players WHERE name = $1`

	rows, err := db.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&player.ID, &player.Name); err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}
	}

	return &player, nil
}

func UpsertPlayersScore(ctx context.Context, db *pgx.Conn, playersScore *model.PlayersScore) (*model.PlayersScore, error) {
	query := `INSERT INTO players_scores 
	(score, player_id, updated_at) VALUES ($1, $2, $3)
ON CONFLICT ON CONSTRAINT unique_player_fk 
DO UPDATE SET score = GREATEST($1, players_scores.score)`

	_, err := db.Exec(ctx, query, playersScore.Score, playersScore.Player.ID, playersScore.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("upsert players score: %w", err)
	}

	return playersScore, nil
}

func GetPlayersScoreByPlayerName(ctx context.Context, db *pgx.Conn, playerName string) (*model.PlayersScore, error) {
	var playersScore = &model.PlayersScore{
		Player: &model.Player{},
	}

	query := `SELECT id, score, player_id, name, updated_at, rank FROM (
		SELECT
			ps.id,
			ps.score,
			ps.player_id,
			p.name,
			ps.updated_at,
			row_number() OVER (ORDER BY score DESC) AS rank
		FROM players_scores ps
		LEFT JOIN players p
			ON ps.player_id = p.id
		ORDER BY score DESC
		) a WHERE name = $1;`

	if err := db.QueryRow(ctx, query, playerName).
		Scan(
			&playersScore.ID,
			&playersScore.Score,
			&playersScore.Player.ID,
			&playersScore.Player.Name,
			&playersScore.UpdatedAt,
			&playersScore.Player.Rank,
		); err != nil {
		return nil, fmt.Errorf("get players score by player name: %w", err)
	}

	return playersScore, nil
}

type GetPlayersScoresOpts struct {
	Offset   int // number of items to offset from list
	Limit    int // number of items to return
	Month    int
	Year     int
	FromRank int
	ToRank   int
}

func GetPlayersScores(ctx context.Context, db *pgx.Conn, opts GetPlayersScoresOpts) ([]model.PlayersScore, error) {
	var (
		q    strings.Builder
		args []interface{}
	)

	var skipClauses bool

	if opts.ToRank > 0 {
		q.WriteString(`SELECT id, score, player_id, name, updated_at, rank FROM (`)
		skipClauses = true
	}

	q.WriteString(`SELECT
	ps.id,
	ps.score,
	ps.player_id,
	p.name,
	ps.updated_at,
	row_number() OVER (ORDER BY score DESC) AS rank
FROM players_scores ps
LEFT JOIN players p
	ON ps.player_id = p.id`)

	var whereClause bool

	write := func(s string) {
		if whereClause {
			q.WriteString("\n\tAND " + s)
		} else {
			q.WriteString("\nWHERE " + s)
			whereClause = true
		}
	}

	arg := func(arg interface{}) string {
		args = append(args, arg)
		return "$" + strconv.Itoa(len(args))
	}

	if !skipClauses {
		if opts.Year > 0 {
			write("EXTRACT(YEAR FROM updated_at) = " + arg(opts.Year))
		}

		if opts.Month > 0 {
			write("EXTRACT(MONTH FROM updated_at) = " + arg(opts.Month))
		}
	}

	q.WriteString("\tORDER BY score DESC")

	if opts.ToRank > 0 {
		q.WriteString(`) q WHERE rank BETWEEN ` + arg(opts.FromRank) + ` AND ` + arg(opts.ToRank))
		skipClauses = true
	}

	if !skipClauses {
		if opts.Offset > 0 {
			q.WriteString("\nOFFSET " + arg(opts.Offset))
		}

		if opts.Limit > 0 {
			q.WriteString("\nLIMIT " + arg(opts.Offset))
		}
	}

	rows, err := db.Query(ctx, q.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("select players scores: %w", err)
	}

	defer rows.Close()

	var playersScores []model.PlayersScore

	for rows.Next() {
		playersScore := model.PlayersScore{
			Player: &model.Player{},
		}

		if err := rows.Scan(
			&playersScore.ID,
			&playersScore.Score,
			&playersScore.Player.ID,
			&playersScore.Player.Name,
			&playersScore.UpdatedAt,
			&playersScore.Player.Rank,
		); err != nil {
			return nil, fmt.Errorf("scan player's score: %w", err)
		}

		playersScores = append(playersScores, playersScore)
	}

	return playersScores, nil
}
