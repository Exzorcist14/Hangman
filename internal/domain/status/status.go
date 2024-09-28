package status

// Status хранит отображаемое слово и количество отгаданных букв.
type Status struct {
	DisplayedWord  []rune
	guessedLetters int
}

// New инициализирует текущее состояние ответа на основе загаданного слова.
func New(word string) Status {
	as := Status{guessedLetters: 0}

	as.DisplayedWord = make([]rune, len([]rune(word)))
	for i := range as.DisplayedWord {
		as.DisplayedWord[i] = '_'
	}

	return as
}

// ShowLetters проявляет в отображаемом слове заданную букву на указанных позициях.
func (as *Status) ShowLetters(letter rune, positions []int) {
	for _, p := range positions {
		as.DisplayedWord[p] = letter
		as.guessedLetters++
	}
}

// IsGuessed возвращает true, если слово отгадано, иначе false.
func (as *Status) IsGuessed() bool {
	return as.guessedLetters == len(as.DisplayedWord)
}
