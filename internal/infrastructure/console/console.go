package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/conditions"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/frames"
)

// GameConsole реализует игровую консоль, с которой взаимодействует пользователь.
type GameConsole struct {
	reader bufio.Reader
	writer bufio.Writer
}

func NewGameConsole() *GameConsole {
	return &GameConsole{
		reader: *bufio.NewReader(os.Stdin),
		writer: *bufio.NewWriter(os.Stdout),
	}
}

// ChooseConditions возвращает категорию и уровня сложности.
func (gc *GameConsole) ChooseConditions(
	cts conditions.Categories,
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
		gc.print("Введите букву (? для подсказки): ", 0)

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
	gc.printf(1, "Подсказка: %s", hint)
}

// DisplaySessionStatus выводит статус сессии.
func (gc *GameConsole) DisplaySessionStatus(
	category, difficulty string,
	fr frames.Frame,
	displayedWord []rune,
	attempts int,
	lettersUsed map[rune]struct{},
) {
	gc.write("----------------------------------------------------------------------------------------", 2)
	gc.writef(1, "Категория: %s", category)
	gc.writef(2, "Уровень сложности: %s", difficulty)
	gc.writeFrame(fr, 2)
	gc.writeLettersUsed(lettersUsed, 1)
	gc.writef(1, "Доступно попыток: %v", attempts)
	gc.writeDisplayedWord(displayedWord, 2)
	gc.flush()
}

// PlayAnimation проигрывает анимацию.
func (gc *GameConsole) PlayAnimation(frs []frames.Frame, msDelay int) {
	for _, fr := range frs {
		gc.writeFrame(fr, 3)
		gc.flush()
		time.Sleep(time.Duration(msDelay) * time.Millisecond)
	}
}

// chooseCategory отображает категории и возвращает выбор.
func (gc *GameConsole) chooseCategory(cts conditions.Categories, randomSelectionCommand string) (string, error) {
	gc.displayConditions(cts)
	return gc.enterCondition(cts, randomSelectionCommand)
}

// chooseCategory() отображает уровни сложности и возвращает выбор.
func (gc *GameConsole) chooseDifficulty(dfs conditions.Difficulties, randomSelectionCommand string) (string, error) {
	gc.displayConditions(dfs)

	condition, err := gc.enterCondition(dfs, randomSelectionCommand)
	if err != nil {
		return "", fmt.Errorf("can`t enter condition: %w", err)
	}

	return condition, err
}

// displayConditions отображает доступные условия.
func (gc *GameConsole) displayConditions(conds any) {
	switch cs := conds.(type) {
	case conditions.Categories:
		gc.write("Выберите категорию (пропустите для случайного выбора):", 0)

		for cond := range cs {
			gc.write(" "+cond, 0)
		}

		gc.write("", 1)
		gc.flush()
	case conditions.Difficulties:
		gc.write("Выберите уровень сложности (пропустите для случайного выбора):", 0)

		for cond := range cs {
			gc.write(" "+cond, 0)
		}

		gc.write("", 1)
		gc.flush()
	}
}

// enterCondition принимает ввод условия и возвращает его.
func (gc *GameConsole) enterCondition(conds any, randomSelectionCommand string) (string, error) {
	var (
		desiredCondition string
		err              error
	)

	switch cds := conds.(type) {
	case conditions.Categories:
		for {
			desiredCondition, err = gc.readLine()
			if err != nil {
				return "", fmt.Errorf("can`t read category: %w", err)
			}

			_, ok := cds[desiredCondition]
			if ok || desiredCondition == randomSelectionCommand {
				break
			}

			gc.print("Категории не существует. Пожалуйста, выберите одну из представленных категорий", 1)
		}
	case conditions.Difficulties:
		for {
			desiredCondition, err = gc.readLine()
			if err != nil {
				return "", fmt.Errorf("can`t read difficulty: %w", err)
			}

			_, ok := cds[desiredCondition]
			if ok || desiredCondition == randomSelectionCommand {
				break
			}

			gc.print("Уровня сложности не существует. Пожалуйста, выберите один из представленных уровней сложности", 1)
		}
	}

	return desiredCondition, nil
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
func (gc *GameConsole) writeFrame(fr frames.Frame, indents int) {
	for _, line := range fr {
		gc.write(line, 1)
	}

	gc.write("", indents)
}

// writeLettersUsed пишет использованные буквы gc.writer.
func (gc *GameConsole) writeLettersUsed(lettersUsed map[rune]struct{}, indents int) {
	gc.write("Использованные буквы:", 0)

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
