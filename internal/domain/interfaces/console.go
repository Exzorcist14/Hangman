package interfaces

import (
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/conditions"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/frames"
)

// Console описывает интерфейс консоли.
type Console interface {
	ChooseConditions(
		categories conditions.Categories,
		dfs conditions.Difficulties,
		randomSelectionCommand string,
	) (category, difficulty string, err error)
	Enter() (letter rune, err error)
	DisplayHint(hint string)
	DisplaySessionStatus(
		category, difficulty string,
		fr frames.Frame,
		displayedWord []rune,
		attempts int,
		lettersUsed map[rune]struct{},
	)
	PlayAnimation(frs []frames.Frame, msDelay int)
}
