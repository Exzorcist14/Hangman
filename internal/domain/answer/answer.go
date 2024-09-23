package answer

import (
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/words"
)

// Answer хранит слово с подсказкой к нему, его категорию, уровень сложности и таблицу вида "символ - позиции".
type Answer struct {
	words.WordData
	Category   string
	Difficulty string
	table      map[rune][]int
}

// NewAnswer создаёт новый Answer с приведённым в нижний регистр словом.
func NewAnswer(wd words.WordData, category, difficulty string) Answer {
	wd.Word = strings.ToLower(wd.Word)

	return Answer{
		WordData:   wd,
		Category:   category,
		Difficulty: difficulty,
		table:      initTable(wd.Word),
	}
}

// GetLetterPositions возвращает номера позиций, на которых встречается указанная буква.
func (a *Answer) GetLetterPositions(letter rune) []int {
	return a.table[letter]
}

// initTable инициализирует таблицу, сопоставляя буквам их позиции в слове.
func initTable(word string) map[rune][]int {
	table := make(map[rune][]int)

	for i, r := range []rune(word) {
		table[r] = append(table[r], i)
	}

	return table
}
