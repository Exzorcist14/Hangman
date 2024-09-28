package session

import (
	"fmt"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/answer"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/conditions"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/conditions/categories"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/frames/frame"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/frames/stageFramesMap"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/random"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/status"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/words"
	"github.com/es-debug/backend-academy-2024-go-template/internal/view/storyBoard"
)

type Session struct {
	console     console
	answer      answer.Answer
	status      status.Status
	maxAttmeps  int
	storyBoard  stageFramesMap.StageFramesMap
	lettersUsed map[rune]struct{}
}

// Console описывает интерфейс консоли.
type console interface {
	ChooseConditions(
		categories categories.Categories,
		dfs conditions.Difficulties,
		randomSelectionCommand string,
	) (category, difficulty string, err error)
	Enter() (letter rune, err error)
	DisplayHint(hint string)
	DisplaySessionStatus(
		category, difficulty string,
		fr frame.Frame,
		displayedWord []rune,
		attempts int,
		lettersUsed map[rune]struct{},
	)
	PlayAnimation(frs []frame.Frame, msDelay int)
}

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
	sfm stageFramesMap.StageFramesMap,
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
		s.console.PlayAnimation(s.storyBoard["victory"], msFrameDelay)
	} else {
		s.console.PlayAnimation(s.storyBoard["defeat"], msFrameDelay)
	}

	return nil
}

// configure конфигурирует игровую сессию на основе выбора пользователем категории и уровня сложности.
func (s *Session) configure(
	ws words.Words,
	dfs conditions.Difficulties,
	randomSelectionCommand string,
	sfm stageFramesMap.StageFramesMap,
) error {
	categories := categories.New(ws)

	category, difficulty, err := s.console.ChooseConditions(categories, dfs, randomSelectionCommand)
	if err != nil {
		return fmt.Errorf("can`t choose conditions: %w", err)
	}

	if category == randomSelectionCommand {
		category = getRandomCategory(categories)
	}

	if difficulty == randomSelectionCommand {
		difficulty = getRandomDifficulty(dfs)
	}

	wordData := ws.GetRandomWordData(category, difficulty)

	s.maxAttmeps = dfs[difficulty]
	s.answer = answer.New(wordData, category, difficulty)
	s.status = status.New(wordData.Word)
	s.storyBoard = storyBoard.CreateStoryBoard(sfm, s.maxAttmeps)

	return nil
}

// Play запускает проигрывание раунда.
func (s *Session) playRound(attempts int) (int, error) {
	frame := s.storyBoard["process"][s.maxAttmeps-attempts]

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
func getRandomCategory(cts categories.Categories) string {
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