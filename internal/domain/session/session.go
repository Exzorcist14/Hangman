package session

import (
	"fmt"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/answer"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/conditions"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/frames"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/random"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/status"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/words"
	"github.com/es-debug/backend-academy-2024-go-template/internal/view/storyboard"
)

// Session хранит ответ, текущее состояние ответа, максимальное количество попыток, раскадровку,
// множество использованных букв и использует интерфейс console.
type Session struct {
	console     console
	answer      answer.Answer
	status      status.Status
	maxAttmeps  int
	storyboard  frames.StageFramesMap
	lettersUsed map[rune]struct{}
}

// console описывает интерфейс консоли.
type console interface {
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

// New возвращает инициализированную структуру Session с переданной консолью и пустым множеством использованных букв.
func New(console console) Session {
	return Session{
		console:     console,
		lettersUsed: make(map[rune]struct{}),
	}
}

// Play запускает игровую сессию.
func (s *Session) Play(
	ws words.Words,
	dfs conditions.Difficulties,
	randomSelectionCommand string,
	msFrameDelay int,
	sfm frames.StageFramesMap,
) error {
	err := s.configure(ws, dfs, randomSelectionCommand, sfm)
	if err != nil {
		return fmt.Errorf("can`t configure session: %w", err)
	}

	attempts := s.maxAttmeps

	for !s.status.IsGuessed() && attempts != 0 {
		attempts, err = s.playRound(attempts)
		if err != nil {
			return fmt.Errorf("can`t play round: %w", err)
		}
	}

	if s.status.IsGuessed() {
		s.console.PlayAnimation(s.storyboard["victory"], msFrameDelay)
	} else {
		s.console.PlayAnimation(s.storyboard["defeat"], msFrameDelay)
	}

	return nil
}

// configure конфигурирует игровую сессию на основе выбора пользователем категории и уровня сложности.
func (s *Session) configure(
	ws words.Words,
	dfs conditions.Difficulties,
	randomSelectionCommand string,
	sfm frames.StageFramesMap,
) error {
	cts := conditions.NewCategories(ws)

	category, difficulty, err := s.console.ChooseConditions(cts, dfs, randomSelectionCommand)
	if err != nil {
		return fmt.Errorf("can`t choose conditions: %w", err)
	}

	if category == randomSelectionCommand {
		category = getRandomCategory(cts)
	}

	if difficulty == randomSelectionCommand {
		difficulty = getRandomDifficulty(dfs)
	}

	wordData := ws.GetRandomWordData(category, difficulty)

	s.maxAttmeps = dfs[difficulty]
	s.answer = answer.New(wordData, category, difficulty)
	s.status = status.New(wordData.Word)
	s.storyboard = storyboard.CreateStoryboard(sfm, s.maxAttmeps)

	return nil
}

// Play запускает проигрывание раунда.
func (s *Session) playRound(attempts int) (int, error) {
	frame := s.storyboard["process"][s.maxAttmeps-attempts]

	s.console.DisplaySessionStatus(
		s.answer.Category,
		s.answer.Difficulty,
		frame,
		s.status.DisplayedWord,
		attempts,
		s.lettersUsed,
	)

	var (
		letter rune
		err    error
	)

	for {
		letter, err = s.console.Enter()
		if err != nil {
			return 0, fmt.Errorf("can`t enter letter: %w", err)
		}

		if letter == '?' {
			s.console.DisplayHint(s.answer.Hint)
		} else if _, ok := s.lettersUsed[letter]; !ok {
			break
		}
	}

	positions := s.answer.GetLetterPositions(letter)
	s.status.ShowLetters(letter, positions)

	if len(positions) == 0 {
		attempts--
	}

	s.updateLettersUsed(letter)

	return attempts, nil
}

// updateLettersUsed обновляет список использованных букв, внося переданную букву.
func (s *Session) updateLettersUsed(letter rune) {
	s.lettersUsed[letter] = struct{}{}
}

// getRandomCategory возвращает случайную категорию.
func getRandomCategory(cts conditions.Categories) string {
	return getRandomCondition(cts)
}

// getRandomDifficulty возвращает случайный уровень сложности.
func getRandomDifficulty(dfs conditions.Difficulties) string {
	return getRandomCondition(dfs)
}

// getRandomDifficulty возвращает случайное условие.
func getRandomCondition[T any](conds map[string]T) string {
	var condition string

	randIndex := random.RandInt(len(conds))
	i := 0

	for cond := range conds {
		if i == randIndex {
			condition = cond
			break
		}

		i++
	}

	return condition
}
