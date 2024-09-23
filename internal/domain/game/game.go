package game

import (
	"fmt"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/conditions"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/frames"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/interfaces"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/words"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/loader"
)

type Game struct {
	console                interfaces.Console
	words                  words.Words
	difficulties           conditions.Difficulties
	frames                 frames.StageFramesMap
	randomSelectionCommand string
	session                Session
}

func NewGame(console interfaces.Console) Game {
	game := Game{console: console}
	return game
}

// Run запускает игру.
func (g *Game) Run() error {
	err := g.loadGameData()
	if err != nil {
		return fmt.Errorf("can`t loader game data: %w", err)
	}

	g.session = NewSession(g.console)

	err = g.session.Play(g.words, g.difficulties, g.frames, g.randomSelectionCommand)
	if err != nil {
		return fmt.Errorf("can`t play game: %w", err)
	}

	return nil
}

// loadGameData инициализирует данные об игре, загружая их из файлов.
func (g *Game) loadGameData() error {
	g.words = words.NewWords()

	err := loader.LoadDataFromFile("./internal/infrastructure/files/words.json", &g.words)
	if err != nil {
		return fmt.Errorf("can`t load words from file: %w", err)
	}

	g.difficulties = conditions.NewDifficulties()

	err = loader.LoadDataFromFile("./internal/infrastructure/files/difficulties.json", &g.difficulties)
	if err != nil {
		return fmt.Errorf("can`t load conditions from file: %w", err)
	}

	g.frames = frames.NewStageFramesMap(4)

	err = loader.LoadDataFromFile("./internal/infrastructure/files/frames.json", &g.frames)
	if err != nil {
		return fmt.Errorf("can`t load frames from file: %w", err)
	}

	err = loader.LoadDataFromFile("./internal/infrastructure/files/randomSelectionCommand.json", &g.randomSelectionCommand)
	if err != nil {
		return fmt.Errorf("can`t load random selection command: %w", err)
	}

	return nil
}
