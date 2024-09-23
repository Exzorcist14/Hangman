package frames

// StageFramesMap - словарь, хранящий этапы игры и соответствующий им слайс кадров.
type StageFramesMap map[string][]Frame

// NewStageFramesMap возвращает инициализированный StageFramesMap, задавая этапам победы и поражения указанный размер кадра.
func NewStageFramesMap(size int) StageFramesMap {
	sfp := make(StageFramesMap)
	sfp["victory"] = make([]Frame, size)
	sfp["defeat"] = make([]Frame, size)

	return sfp
}
