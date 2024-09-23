package conditions

// Difficulties - словарь, хранящий уровни сложности и соответствующее им максимальное количество попыток.
type Difficulties map[string]int

func NewDifficulties() Difficulties {
	return make(Difficulties)
}