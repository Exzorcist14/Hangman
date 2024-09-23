package main

import (
	"os"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/game"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/console"
)

func main() {
	g := game.NewGame(console.NewGameConsole())

	err := g.Run()
	if err != nil {
		os.Exit(1)
	}
}
