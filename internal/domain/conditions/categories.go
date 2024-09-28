package conditions

import "github.com/es-debug/backend-academy-2024-go-template/internal/domain/words"

// Categories - множество категорий.
type Categories map[string]struct{}

// New инициализирует Categories по переданному словарю игры.
func New(ws words.Words) Categories {
	categories := make(map[string]struct{})

	for ct := range ws {
		categories[ct] = struct{}{}
	}

	return categories
}
