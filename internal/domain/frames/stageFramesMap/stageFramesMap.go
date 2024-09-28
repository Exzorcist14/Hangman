package stageFramesMap

import (
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/frames/frame"
)

// StageFramesMap - словарь, хранящий этапы игры и соответствующий им слайс кадров.
type StageFramesMap map[string][]frame.Frame

// NewStageFramesMap возвращает инициализированный StageFramesMap, задавая этапам победы и поражения указанный размер кадра.
func New(size int) StageFramesMap {
	sfp := make(StageFramesMap)
	sfp["victory"] = make([]frame.Frame, size)
	sfp["defeat"] = make([]frame.Frame, size)

	return sfp
}
