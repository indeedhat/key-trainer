package internal

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/buger/goterm"
	"github.com/fatih/color"
)

const (
	paragraphSize    = 100
	defaultTermWidth = 80
)

type runnerLine struct {
	wordLen int
	charLen int
	words   []string
}

type Runner struct {
	config RunnerConfig

	// input
	reader *bufio.Reader

	// render
	termWidth int

	// place tracking
	buffer          string
	lines           []*runnerLine
	completeInLine  int
	errorInLastWord bool
	completedTotal  int

	// wordlist stuff
	wordList      []string
	wordListLen   int
	wordCount     int
	wordPickTotal int

	print struct {
		error *color.Color
		done  *color.Color
	}
}

// NewRunner just returns a reference to new Runner struct
func NewRunner(config ...RunnerConfig) *Runner {
	var conf RunnerConfig

	if len(config) > 0 {
		conf = config[0]
	}

	return &Runner{config: conf}
}

// run the mode
func (run *Runner) Run(wordList string) {
	run.wordList = loadFromFile(wordList, run.config.Contains)
	run.wordListLen = len(run.wordList)
	run.termWidth = defaultTermWidth
	run.reader = bufio.NewReader(os.Stdin)

	run.print.error = color.New(color.BgRed, color.FgWhite)
	run.print.done = color.New(color.BgGreen, color.FgBlack)

	for {
		if run.config.TotalWords != 0 &&
			run.config.TotalWords <= uint(run.completedTotal) {

			fmt.Println("\nWell Done!")
			break
		}

		run.advanceLines()
		run.resize()
		run.render()
		run.handleInput()
	}
}

// render the terminal output
func (run *Runner) render() {
	goterm.Clear()
	goterm.MoveCursor(1, 1)
	goterm.Flush()

	fmt.Println("")

	for i, line := range run.lines {
		fmt.Print("  ")
		if i != 0 {
			fmt.Println(strings.Join(line.words, " "))
			continue
		}

		if run.completeInLine > 0 {
			for _, word := range line.words[:run.completeInLine] {
				run.print.done.Print(word)
				fmt.Print(" ")
			}
		}

		currentWord := line.words[run.completeInLine]
		bufferLen := len(run.buffer)

		if bufferLen == 0 {
			fmt.Print(currentWord)
		} else if run.errorInLastWord {
			run.print.error.Print(currentWord)
		} else {
			run.print.done.Print(string(currentWord[:bufferLen]))
			fmt.Print(string(currentWord[bufferLen:]))
		}

		fmt.Print(" ")
		fmt.Print(strings.Join(line.words[run.completeInLine+1:], " ") + "\n")
	}

	fmt.Print("\n  ", run.buffer)
}

// resize the window
// this aso remaps the whole word mapping
func (run *Runner) resize() {
	width := goterm.Width()
	if width == 0 {
		width = defaultTermWidth
	}

	if width-4 == run.termWidth {
		return
	}

	run.termWidth = width - 4

	var words []string
	for _, line := range run.lines {
		words = append(words, line.words...)
	}

	run.lines = []*runnerLine{}
	cursor := run.lastLine()

	for i := 0; i < run.wordCount; i++ {
		newWord := words[i]
		newWordLen := len(newWord)

		if cursor.charLen+newWordLen+1 > run.termWidth {
			cursor = &runnerLine{}
			run.lines = append(run.lines, cursor)
		}

		cursor.words = append(cursor.words, newWord)
		cursor.wordLen++
		cursor.charLen += newWordLen + 1
	}

	run.advanceLines()
}

// advanceLines will remove any complete lines and add new ones
// to keep up to the total word count
func (run *Runner) advanceLines() {
	if run.completeInLine != 0 && run.completeInLine == run.lines[0].wordLen {
		run.lines = run.lines[1:]
		run.wordCount -= run.completeInLine
		run.completeInLine = 0
	}

	if run.wordPickTotal > int(run.config.TotalWords) {
		return
	}

	cursor := run.lastLine()
	for ; run.wordCount < paragraphSize; run.wordCount++ {
		newWord := run.wordList[rand.Intn(run.wordListLen)]
		newWordLen := len(newWord)

		if cursor.charLen+newWordLen+1 > run.termWidth {
			cursor = &runnerLine{}
			run.lines = append(run.lines, cursor)
		}

		cursor.words = append(cursor.words, newWord)
		cursor.wordLen++
		cursor.charLen += newWordLen + 1
		run.wordPickTotal++

		if run.wordPickTotal > int(run.config.TotalWords) {
			return
		}
	}
}

func (run *Runner) lastLine() *runnerLine {
	var cursor *runnerLine

	if len(run.lines) == 0 {
		cursor = &runnerLine{}
		run.lines = append(run.lines, cursor)
	} else {
		cursor = run.lines[len(run.lines)-1]
	}

	return cursor
}

func (run *Runner) handleInput() {
	input, _ := run.reader.ReadByte()
	char := string(input)
	subject := run.lines[0].words[run.completeInLine]

	if char == "\n" {
		return
	}

	switch char {
	// backspace
	case "\x7f":
		if len(run.buffer) > 0 {
			run.buffer = run.buffer[:len(run.buffer)-1]
		}

	case " ":
		if run.buffer == subject {
			run.buffer = ""
			run.completeInLine++
			run.completedTotal++
		} else {
			run.buffer += char
		}

	default:
		run.buffer += char
	}

	run.errorInLastWord = len(run.buffer) == 0 ||
		!strings.HasPrefix(subject, run.buffer)
}
