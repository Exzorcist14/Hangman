package game

import (
	"fmt"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/config"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/frames"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/session"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/words"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/console"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/loader"
)

// Game хранит конфиг, словарь, кадры и сессию.
type Game struct {
	config         config.Config
	words          words.Words
	stageFramesMap frames.StageFramesMap
	session        session.Session
}

// New возвращает инициализированную структуру Game.
func New() (Game, error) {
	g := Game{}

	err := g.loadGameData()
	if err != nil {
		return g, fmt.Errorf("can`t load game data: %w", err)
	}

	return g, nil
}

// Run запускает игру.
func (g *Game) Run() error {
	g.session = session.New(console.New())

	err := g.session.Play(g.words, g.config.Difficulties, g.config.RandomSelectionCommand, g.config.MsFrameDelay, g.stageFramesMap)
	if err != nil {
		return fmt.Errorf("can`t play session: %w", err)
	}

	return nil
}

// loadGameData инициализирует данные об игре, загружая их из файлов.
func (g *Game) loadGameData() error {
	g.config = config.New()

	err := loader.LoadDataFromFile("./internal/infrastructure/files/config.json", &g.config)
	if err != nil {
		return fmt.Errorf("can`t load config from file: %w", err)
	}

	g.words = make(words.Words)

	err = loader.LoadDataFromFile("./internal/infrastructure/files/words.json", &g.words)
	if err != nil {
		return fmt.Errorf("can`t load words from file: %w", err)
	}

	g.stageFramesMap = frames.New(g.config.FramesInAnimation)

	err = loader.LoadDataFromFile("./internal/infrastructure/files/frames.json", &g.stageFramesMap)
	if err != nil {
		return fmt.Errorf("can`t load stageFramesMap from file: %w", err)
	}

	return nil
}
