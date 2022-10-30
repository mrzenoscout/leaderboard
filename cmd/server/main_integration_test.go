package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	go main()

	code := m.Run()
	os.Exit(code)
}

func TestSavePlayersScore(t *testing.T) {
	type Request struct {
		Name  string `json:"name"`
		Score int    `json:"score"`
	}

	request := Request{
		Name:  gofakeit.Name(),
		Score: gofakeit.Number(0, 100000),
	}

	b, err := json.Marshal(request)

	assert.NoError(t, err)

	w := httptest.NewRecorder()
	resp, err := http.Post("/leaderboard/score", "application/json", bytes.NewBuffer(b))
	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)

	type Response struct {
		Rank int `json:"rank"`
	}

	var res Response

	err = json.NewDecoder(resp.Body).Decode(&res)
	assert.NoError(t, err)

	assert.NotEmpty(t, res)
}
