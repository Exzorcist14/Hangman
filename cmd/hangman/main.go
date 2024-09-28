package main

import (
	"os"

	"github.com/es-debug/backend-academy-2024-go-template/internal/application/game"
)

func main() {
	g := game.New()

	err := g.Run()
	if err != nil {
		os.Exit(1)
	}
}
