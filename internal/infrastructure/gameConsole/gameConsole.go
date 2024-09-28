package gameConsole

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/conditions"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/conditions/categories"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/frames/frame"
)

const (
	letterInputMessage       = "Введите букву (? для подсказки): "
	border                   = "----------------------------------------------------------------------------------------"
	hintForm                 = "Подсказка: %s"
	categoryForm             = "Категория: %s"
	difficultyForm           = "Уровень сложности: %s"
	attemptsForm             = "Доступно попыток: %v"
	categoryInputMessage     = "Выберите категорию (пропустите для случайного выбора):"
	invalidCategoryMessage   = "Категории не существует. Пожалуйста, выберите одну из представленных категорий"
	difficultyInputMessage   = "Выберите уровень сложности (пропустите для случайного выбора):"
	invalidDifficultyMessage = "Уровня сложности не существует. Пожалуйста, выберите один из представленных уровней сложности"
	lettersUsedMessage       = "Использованные буквы:"
)

// GameConsole реализует игровую консоль, с которой взаимодействует пользователь.
type GameConsole struct {
	reader bufio.Reader
	writer bufio.Writer
}

func New() *GameConsole {
	return &GameConsole{
		reader: *bufio.NewReader(os.Stdin),
		writer: *bufio.NewWriter(os.Stdout),
	}
}

// ChooseConditions возвращает категорию и уровня сложности.
func (gc *GameConsole) ChooseConditions(
	cts categories.Categories,
	dfs conditions.Difficulties,
	randomSelectionCommand string,
) (category, difficulty string, err error) {
	category, err = gc.chooseCategory(cts, randomSelectionCommand)
	if err != nil {
		return "", "", fmt.Errorf("can`t choose category: %w", err)
	}

	difficulty, err = gc.chooseDifficulty(dfs, randomSelectionCommand)
	if err != nil {
		return "", "", fmt.Errorf("can`t choose difficulty: %w", err)
	}

	return category, difficulty, nil
}

// Enter принимает ввод буквы без учёта регистра.
func (gc *GameConsole) Enter() (rune, error) {
	var lineRemainder string

	for {
		gc.print(letterInputMessage, 0)

		r, _, err := gc.reader.ReadRune()
		if err != nil {
			return ' ', fmt.Errorf("can`t read rune: %w", err)
		}

		if unicode.IsLetter(r) || r == '?' {
			lineRemainder, err = gc.reader.ReadString('\n') // Очищаем буфер
			if err != nil {
				return ' ', fmt.Errorf("can`t read rest of line: %w", err)
			}

			if lineRemainder == "\n" {
				return unicode.ToLower(r), nil
			}
		}
	}
}

// DisplayHint выводит передаваемую подсказку.
func (gc *GameConsole) DisplayHint(hint string) {
	gc.printf(1, hintForm, hint)
}

// DisplaySessionStatus выводит статус сессии.
func (gc *GameConsole) DisplaySessionStatus(
	category, difficulty string,
	fr frame.Frame,
	displayedWord []rune,
	attempts int,
	lettersUsed map[rune]struct{},
) {
	gc.write(border, 2)
	gc.writef(1, categoryForm, category)
	gc.writef(2, difficultyForm, difficulty)
	gc.writeFrame(fr, 2)
	gc.writeLettersUsed(lettersUsed, 1)
	gc.writef(1, attemptsForm, attempts)
	gc.writeDisplayedWord(displayedWord, 2)
	gc.flush()
}

// PlayAnimation проигрывает анимацию.
func (gc *GameConsole) PlayAnimation(frs []frame.Frame, msDelay int) {
	for _, fr := range frs {
		gc.writeFrame(fr, 3)
		gc.flush()
		time.Sleep(time.Duration(msDelay) * time.Millisecond)
	}
}

// chooseCategory отображает категории и возвращает выбор.
func (gc *GameConsole) chooseCategory(cts categories.Categories, randomSelectionCommand string) (string, error) {
	gc.write(categoryInputMessage, 0)

	for ct := range cts {
		gc.write(" "+ct, 0)
	}

	gc.write("", 1)
	gc.flush()

	condition, err := gc.enterCategory(cts, randomSelectionCommand)
	if err != nil {
		return "", fmt.Errorf("can`t enter category: %w", err)
	}

	return condition, err
}

// chooseCategory() отображает уровни сложности и возвращает выбор.
func (gc *GameConsole) chooseDifficulty(dfs conditions.Difficulties, randomSelectionCommand string) (string, error) {
	gc.write(difficultyInputMessage, 0)

	for df := range dfs {
		gc.write(" "+df, 0)
	}

	gc.write("", 1)
	gc.flush()

	condition, err := gc.enterDifficulty(dfs, randomSelectionCommand)
	if err != nil {
		return "", fmt.Errorf("can`t enter difficulty: %w", err)
	}

	return condition, err
}

// enterCategory принимает ввод категории и возвращает её.
func (gc *GameConsole) enterCategory(cts categories.Categories, randomSelectionCommand string) (string, error) {
	var (
		category string
		err      error
	)

	for {
		category, err = gc.readLine()
		if err != nil {
			return "", fmt.Errorf("can`t read category: %w", err)
		}

		_, ok := cts[category]
		if ok || category == randomSelectionCommand {
			break
		}

		gc.print(invalidCategoryMessage, 1)
	}

	return category, nil
}

// enterDifficulty принимает ввод сложности и возвращает его.
func (gc *GameConsole) enterDifficulty(dfs conditions.Difficulties, randomSelectionCommand string) (string, error) {
	var (
		difficulty string
		err        error
	)

	for {
		difficulty, err = gc.readLine()
		if err != nil {
			return "", fmt.Errorf("can`t read difficulty: %w", err)
		}

		_, ok := dfs[difficulty]
		if ok || difficulty == randomSelectionCommand {
			break
		}

		gc.print(invalidDifficultyMessage, 1)
	}

	return difficulty, nil
}

// readLine читает строки без учёта регистра.
func (gc *GameConsole) readLine() (string, error) {
	word, err := gc.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("can`t read string: %w", err)
	}

	return strings.ToLower(strings.TrimSpace(word)), nil
}

// write пишет данные в gc.writer.
func (gc *GameConsole) write(data string, indents int) {
	fmt.Fprint(&gc.writer, data)
	gc.writeIndent(indents)
}

// writef форматирует и пишет данные в gc.writer.
func (gc *GameConsole) writef(indents int, format string, a ...any) {
	fmt.Fprintf(&gc.writer, format, a...)
	gc.writeIndent(indents)
}

// writeIndent пишет переданное количество переносов строки.
func (gc *GameConsole) writeIndent(count int) {
	for i := 0; i < count; i++ {
		fmt.Fprint(&gc.writer, "\n")
	}
}

// flush выбрасывает данные из gc.writer.
func (gc *GameConsole) flush() {
	gc.writer.Flush()
}

// print пишет данные и выбрасывает их.
func (gc *GameConsole) print(data string, indents int) {
	gc.write(data, indents)
	gc.flush()
}

// printf форматирует, пишет данные и выбрасывает их.
func (gc *GameConsole) printf(indents int, format string, a ...any) {
	gc.writef(indents, format, a...)
	gc.flush()
}

// writeFrame пишет линии кадра fr в gc.writer.
func (gc *GameConsole) writeFrame(fr frame.Frame, indents int) {
	for _, line := range fr {
		gc.write(line, 1)
	}

	gc.write("", indents)
}

// writeLettersUsed пишет использованные буквы gc.writer.
func (gc *GameConsole) writeLettersUsed(lettersUsed map[rune]struct{}, indents int) {
	gc.write(lettersUsedMessage, 0)

	for letter := range lettersUsed {
		gc.write(" "+string(letter), 0)
	}

	gc.write("", indents)
}

// writeDisplayedWord пишет символы отображаемого слова в gc.writer.
func (gc *GameConsole) writeDisplayedWord(displayedWord []rune, indents int) {
	for _, letter := range displayedWord {
		gc.write(string(letter), 0)
	}

	gc.write("", indents)
}
