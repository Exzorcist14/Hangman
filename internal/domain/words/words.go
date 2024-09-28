package words

import (
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/random"
)

// Words - cловарь, сопоставляющий категории и сложности слайс данных о слове.
type Words map[string]map[string][]WordData

// GetRandomWordData возвращает случайное слово из словаря.
func (ws Words) GetRandomWordData(category, difficulty string) WordData {
	randIndex := random.RandInt(len(ws[category][difficulty]))
	return ws[category][difficulty][randIndex]
}
