package words_test

import (
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/words"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/loader"
	"github.com/stretchr/testify/assert"
)

func TestGetRandomWordData(t *testing.T) {
	tt := []struct {
		category   string
		difficulty string
		correspond bool
	}{
		{
			category:   "персонажи",
			difficulty: "лёгкая",
			correspond: true,
		},
		{
			category:   "видеоигры",
			difficulty: "средняя",
			correspond: true,
		},
		{
			category:   "пища",
			difficulty: "трудная",
			correspond: true,
		},
	}

	for _, tc := range tt {
		ws := words.NewWords()

		err := loader.LoadDataFromFile("../../infrastructure/files/words.json", &ws)
		assert.NoError(t, err)

		word := ws.GetRandomWordData(tc.category, tc.difficulty)
		wordsData := ws[tc.category][tc.difficulty]

		ok := false

		for _, wd := range wordsData {
			if wd == word {
				ok = true
				break
			}
		}

		assert.Equal(t, tc.correspond, ok)
	}
}
