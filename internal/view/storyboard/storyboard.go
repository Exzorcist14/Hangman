package storyboard

import (
	"math"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/frames"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/random"
)

// CreateStoryboard создаёт раскадровку типа StageFramesMap по входному набору кадров и количеству попыток.
func CreateStoryboard(sfp frames.StageFramesMap, attempts int) frames.StageFramesMap {
	storyboard := frames.New(len(sfp["victory"]))
	copy(storyboard["defeat"], sfp["defeat"])   // кадры анимации поражения соответствуют предусмотренному набору кадров
	copy(storyboard["victory"], sfp["victory"]) // кадры анимации победы соответствуют предусмотренному набору кадров

	frameIndexes := generateFrameIndexes(sfp, attempts)
	for _, frameIndex := range frameIndexes {
		frame := make(frames.Frame, len(sfp["process"][frameIndex]))
		copy(frame, sfp["process"][frameIndex])

		storyboard["process"] = append(storyboard["process"], frame)
	}

	return storyboard
}

// generateFrameIndexes генерирует номера кадров из исходного набора, которые будут включены в раскадровку.
func generateFrameIndexes(frs frames.StageFramesMap, attempts int) []int {
	// Выберем кадры для каждой новой попытки
	framesNumber := len(frs["process"]) // количество кадров в изначальном наборе
	segmentsNumber := attempts          // количество кадров, необходимое для соответствия каждой новой попытке
	segmentLength := int(math.Ceil(float64(framesNumber) / float64(segmentsNumber)))

	frameIndexes := make([]int, segmentsNumber) // номера кадров, которые нужно включить в раскадровку
	frameIndexes[0] = 0                         // нулевой кадр всегда включается в раскадровку

	// генерируем в границах сегмента индекс кадра
	for i := 1; i < segmentsNumber-1; i++ {
		index := i*segmentLength + random.RandInt(segmentLength)
		frameIndexes[i] = index
	}

	frameIndexes[segmentsNumber-1] = framesNumber - 1 // Последний кадр всегда включается в раскадровку

	return frameIndexes
}
