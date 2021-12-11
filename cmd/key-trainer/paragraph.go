package main

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

type paragraphLine struct {
	wordLen int
	charLen int
	words   []string
}

type paragraphMode struct {
	// input
	reader *bufio.Reader

	// render
	termWidth int
	// place tracking
	buffer          string
	lines           []*paragraphLine
	completeInLine  int
	errorInLastWord bool

	// wordlist stuff
	wordList    []string
	wordListLen int
	wordCount   int

	print struct {
		error *color.Color
		done  *color.Color
	}
}

// run the mode
func (pm *paragraphMode) run() {
	pm.wordList = loadFromFile()
	pm.wordListLen = len(pm.wordList)
	pm.termWidth = defaultTermWidth
	pm.reader = bufio.NewReader(os.Stdin)

	pm.print.error = color.New(color.BgRed, color.FgWhite)
	pm.print.done = color.New(color.BgGreen, color.FgBlack)

	for {
		pm.advanceLines()
		pm.resize()
		pm.render()
		pm.handleInput()
	}
}

// render the terminal output
func (pm *paragraphMode) render() {
	goterm.Clear()
	goterm.MoveCursor(1, 1)
	goterm.Flush()

	fmt.Println("")

	for i, line := range pm.lines {
		fmt.Print("  ")
		if i != 0 {
			fmt.Println(strings.Join(line.words, " "))
			continue
		}

		if pm.completeInLine > 0 {
			for _, word := range line.words[:pm.completeInLine] {
				pm.print.done.Print(word)
				fmt.Print(" ")
			}
		}

		currentWord := line.words[pm.completeInLine]
		bufferLen := len(pm.buffer)

		if bufferLen == 0 {
			fmt.Print(currentWord)
		} else if pm.errorInLastWord {
			pm.print.error.Print(currentWord)
		} else {
			pm.print.done.Print(string(currentWord[:bufferLen]))
			fmt.Print(string(currentWord[bufferLen:]))
		}

		fmt.Print(" ")
		fmt.Print(strings.Join(line.words[pm.completeInLine+1:], " ") + "\n")
	}

	fmt.Print("\n  ", pm.buffer)
}

// resize the window
// this aso remaps the whole word mapping
func (pm *paragraphMode) resize() {
	width := goterm.Width()
	if width == 0 {
		width = defaultTermWidth
	}

	if width-4 == pm.termWidth {
		return
	}

	pm.termWidth = width - 4

	var words []string
	for _, line := range pm.lines {
		words = append(words, line.words...)
	}

	pm.lines = []*paragraphLine{}
	cursor := pm.lastLine()

	for i := 0; i < pm.wordCount; i++ {
		newWord := words[i]
		newWordLen := len(newWord)

		if cursor.charLen+newWordLen+1 > pm.termWidth {
			cursor = &paragraphLine{}
			pm.lines = append(pm.lines, cursor)
		}

		cursor.words = append(cursor.words, newWord)
		cursor.wordLen++
		cursor.charLen += newWordLen + 1
	}

	pm.advanceLines()
}

// advanceLines will remove any complete lines and add new ones
// to keep up to the total word count
func (pm *paragraphMode) advanceLines() {
	if pm.completeInLine != 0 && pm.completeInLine == pm.lines[0].wordLen {
		pm.lines = pm.lines[1:]
		pm.wordCount -= pm.completeInLine
		pm.completeInLine = 0
	}

	cursor := pm.lastLine()
	for ; pm.wordCount < paragraphSize; pm.wordCount++ {
		newWord := pm.wordList[rand.Intn(pm.wordListLen)]
		newWordLen := len(newWord)

		if cursor.charLen+newWordLen+1 > pm.termWidth {
			cursor = &paragraphLine{}
			pm.lines = append(pm.lines, cursor)
		}

		cursor.words = append(cursor.words, newWord)
		cursor.wordLen++
		cursor.charLen += newWordLen + 1
	}
}

func (pm *paragraphMode) lastLine() *paragraphLine {
	var cursor *paragraphLine

	if len(pm.lines) == 0 {
		cursor = &paragraphLine{}
		pm.lines = append(pm.lines, cursor)
	} else {
		cursor = pm.lines[len(pm.lines)-1]
	}

	return cursor
}

func (pm *paragraphMode) handleInput() {
	input, _ := pm.reader.ReadByte()
	char := string(input)
	subject := pm.lines[0].words[pm.completeInLine]

	if char == "\n" {
		return
	}

	// backspace
	switch char {
	case "\x7f":
		if len(pm.buffer) > 0 {
			pm.buffer = pm.buffer[:len(pm.buffer)-1]
		}

	case " ":
		if pm.buffer == subject {
			pm.buffer = ""
			pm.completeInLine++
		} else {
			pm.buffer += char
		}

	default:
		pm.buffer += char
	}

	pm.errorInLastWord = len(pm.buffer) == 0 ||
		!strings.HasPrefix(subject, pm.buffer)
}
