package random

import (
	"crypto/rand"
	"math/big"
)

// RandInt возвращает случайное число из полуинтервала [0, maxValue). Если max < 0, возвращает 0.
func RandInt(maxValue int) int {
	result, err := rand.Int(rand.Reader, big.NewInt(int64(maxValue)))
	if err != nil {
		return 0
	}

	return int(result.Int64())
}
