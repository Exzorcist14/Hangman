package game

import (
	"fmt"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/answer"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/conditions"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/frames"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/interfaces"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/random"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/words"
)

type Session struct {
	console      interfaces.Console
	answer       answer.Answer
	answerStatus answer.Status
	maxAttmeps   int
	attempts     int
	storyBoard   frames.StageFramesMap
	lettersUsed  map[rune]struct{}
}

func NewSession(console interfaces.Console) Session {
	return Session{
		console:     console,
		lettersUsed: make(map[rune]struct{}),
	}
}

// Play запускает игровую сессию.
func (s *Session) Play(
	ws words.Words,
	dfs conditions.Difficulties,
	sfm frames.StageFramesMap,
	randomSelectionCommand string,
) error {
	err := s.configure(ws, dfs, sfm, randomSelectionCommand)
	if err != nil {
		return fmt.Errorf("can`t configure session: %w", err)
	}

	for !s.answerStatus.IsGuessed() && s.attempts != 0 {
		err = s.playRound()
		if err != nil {
			return fmt.Errorf("can`t play round: %w", err)
		}
	}

	if s.answerStatus.IsGuessed() {
		s.console.PlayAnimation(s.storyBoard["victory"], 1250)
	} else {
		s.console.PlayAnimation(s.storyBoard["defeat"], 1250)
	}

	return nil
}

// configure конфигурирует игровую сессию на основе выбора пользователем категории и уровня сложности.
func (s *Session) configure(
	ws words.Words,
	dfs conditions.Difficulties,
	sfm frames.StageFramesMap,
	randomSelectionCommand string,
) error {
	categories := conditions.NewCategories(ws)

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
	s.attempts = s.maxAttmeps
	s.answer = answer.NewAnswer(wordData, category, difficulty)
	s.answerStatus = answer.NewStatus(wordData.Word)
	s.storyBoard = frames.CreateStoryBoard(sfm, s.attempts)

	return nil
}

// Play запускает проигрывание раунда.
func (s *Session) playRound() error {
	frame := s.storyBoard["process"][s.maxAttmeps-s.attempts]

	s.console.DisplaySessionStatus(
		s.answer.Category,
		s.answer.Difficulty,
		frame,
		s.answerStatus.DisplayedWord,
		s.attempts,
		s.lettersUsed,
	)

	var (
		letter rune
		err    error
	)

	for {
		letter, err = s.console.Enter()
		if err != nil {
			return fmt.Errorf("can`t enter letter: %w", err)
		}

		if letter == '?' {
			s.console.DisplayHint(s.answer.Hint)
		} else if _, ok := s.lettersUsed[letter]; !ok {
			break
		}
	}

	positions := s.answer.GetLetterPositions(letter)
	s.answerStatus.ShowLetters(letter, positions)

	if len(positions) == 0 {
		s.attempts--
	}

	s.updateLettersUsed(letter)

	return nil
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
func getRandomCondition(conds any) string {
	var condition string

	switch cds := conds.(type) {
	case conditions.Categories:
		randIndex := random.RandInt(len(cds))

		i := 0
		for ct := range cds {
			if i == randIndex {
				condition = ct
				break
			}

			i++
		}
	case conditions.Difficulties:
		randIndex := random.RandInt(len(cds))

		i := 0
		for df := range cds {
			if i == randIndex {
				condition = df
				break
			}

			i++
		}
	}

	return condition
}
