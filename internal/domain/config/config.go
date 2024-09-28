package config

import "github.com/es-debug/backend-academy-2024-go-template/internal/domain/conditions"

type Config struct {
	Difficulties           conditions.Difficulties
	RandomSelectionCommand string
	FramesInAnimation      int
	MsFrameDelay           int
}

// New возвращает инициализированный Config с предустановленными настройками по-умолчанию
func New() Config {
	return Config{
		Difficulties: conditions.Difficulties{
			"лёгкая":  7,
			"средняя": 5,
			"трудная": 3,
		},
		RandomSelectionCommand: "",
		FramesInAnimation:      4,
		MsFrameDelay:           1250,
	}
}
