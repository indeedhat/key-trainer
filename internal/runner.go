package internal

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/buger/goterm"
	"github.com/fatih/color"
)

const (
	paragraphSize    = 100
	defaultTermWidth = 80

	cpmToWpm = 0.2
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
	buffer         string
	lines          []*runnerLine
	completeInLine int
	typo           bool

	// wordlist stuff
	wordList      []string
	wordListLen   int
	wordCount     int
	wordPickTotal int

	print struct {
		error *color.Color
		done  *color.Color
	}

	stats struct {
		startTime      *time.Time
		complete       uint
		characterCount uint
		typos          uint
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

	run.DisplayStartScreen(wordList)

	for {
		if run.config.TotalWords != 0 &&
			run.config.TotalWords <= uint(run.stats.complete) {

			run.DisplayStatusScreen()
			break
		}

		run.advanceLines()
		run.resize()
		run.render()
		run.handleInput()
	}
}

// DisplayStartScreen will briefly show settings befor starting the test
func (run *Runner) DisplayStartScreen(wordList string) {
	contains := "N/A"
	if run.config.Contains != "" {
		contains = run.config.Contains
	}

	totalWords := "N/A"
	if run.config.TotalWords != 0 {
		totalWords = fmt.Sprint(run.config.TotalWords)
	}

	for i := 5; i > 0; i-- {
		goterm.Clear()
		goterm.MoveCursor(0, 0)
		goterm.Flush()

		fmt.Printf(
			`
  Typing Test:
    
  tracking will begin on your first keypress

  SETTINGS:
  Word List:    %s
  Contains:     %s
  Total Words:  %s

  Starting in %d
    `,
			wordList,
			contains,
			totalWords,
			i,
		)
		time.Sleep(time.Second)
	}
}

// DisplayStatusScreen detailing the runs statistics
// this is intended to be ran on close
func (run *Runner) DisplayStatusScreen() {
	if run.stats.startTime == nil {
		return
	}

	elapsed := time.Now().Sub(*run.stats.startTime)
	cpm := float64(run.stats.characterCount) / elapsed.Minutes()

	fmt.Printf(
		`
  Typing Test Complete!

  Time Elapseed:    %s
  Wpm (corrected):  %d
  Cpm (corrected):  %d

  Words Complete:   %d
  Wpm:              %d

  Typos:            %d
    `,
		elapsed.String(),
		uint(cpm*cpmToWpm),
		uint(cpm),
		run.stats.complete,
		uint(float64(run.stats.complete)/elapsed.Minutes()),
		run.stats.typos,
	)
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
		} else if run.typo {
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
	if run.completeInLine != 0 &&
		run.completeInLine == run.lines[0].wordLen {

		run.lines = run.lines[1:]
		run.wordCount -= run.completeInLine
		run.completeInLine = 0
	}

	if run.config.TotalWords != 0 &&
		run.wordPickTotal > int(run.config.TotalWords) {

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

		if int(run.config.TotalWords) != 0 &&
			run.wordPickTotal > int(run.config.TotalWords) {

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

	if run.stats.startTime == nil {
		now := time.Now()
		run.stats.startTime = &now
	}

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
			run.stats.complete++
			run.stats.characterCount += uint(len(subject) + 1)
		} else {
			run.buffer += char
		}

	default:
		run.buffer += char
	}

	hasTypo := len(run.buffer) == 0 ||
		!strings.HasPrefix(subject, run.buffer)

	if !run.typo && hasTypo {
		run.stats.typos++
	}

	run.typo = hasTypo
}
