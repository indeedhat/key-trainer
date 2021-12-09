package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

var runes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890!@#$%^&*()-=[]_+{};:'/\"|><.,~`"
var specialRunes = "1234567890!@#$%^&*()-=[]_+{};:'/\"|><.,~`"

type flags struct {
	wordMode             bool
	specialCharacterMode bool
	containsMode         string
}

func main() {
	var singleCharacterMode bool
	var wordList []string

	flags := parseFlags()
	rand.Seed(time.Now().Unix())

	// this little bit of unintuative magic disables input buffering
	// so that the reader can read a single byte immediatly without
	// the need for pressing enter
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").
		Run()

		// and this one will stop the characters you type from appearing in
		// the terminal
	exec.Command("stty", "-F", "/dev/tty", "-echo", "min", "1").
		Run()

	switch {
	case flags.wordMode:
		wordList = loadFromFile()

	case flags.containsMode != "":
		rawWordList := loadFromFile()
		for _, word := range rawWordList {
			if strings.Contains(word, flags.containsMode) {
				wordList = append(wordList, word)
			}
		}

	case flags.specialCharacterMode:
		wordList = splitCharacters(specialRunes)
		singleCharacterMode = true

	default:
		wordList = splitCharacters(runes)
		singleCharacterMode = true
	}

	run(wordList, singleCharacterMode)
}

func parseFlags() (f flags) {
	flag.BoolVar(&f.wordMode, "w", false, "open in word mode")
	flag.BoolVar(&f.specialCharacterMode, "s", false, "open in function special character mode")
	flag.StringVar(&f.containsMode, "c", "", "open in word mode and only user words containing the given substring")
	flag.Parse()

	return
}

// handleInput from the users keyboard
func handleInput(reader *bufio.Reader, buffer *string, singleCharacterMode bool) {
	input, _ := reader.ReadByte()
	char := string(input)

	if singleCharacterMode {
		*buffer = char
		return
	}

	if char == "\x7f" {
		fmt.Print("\n\033[1A\033[K")
		*buffer = ""
	} else if char != " " && char != "\n" {
		fmt.Print(string(input))
		*buffer += char
	}
}

func loadFromFile() []string {
	words, err := ioutil.ReadFile("./words.txt")
	if err != nil {
		panic(err)
	}

	return strings.Split(string(words), "\n")
}

func splitCharacters(characters string) (wordList []string) {
	for i, m := 0, len(characters); i < m; i++ {
		wordList = append(wordList, string(characters[i]))
	}

	return
}

func run(wordList []string, singleCharacterMode bool) {
	wordCount := len(wordList)
	reader := bufio.NewReader(os.Stdin)
	buffer := ""

	for {
		subject := wordList[rand.Intn(wordCount)]
		fmt.Printf("\n%s\n", string(subject))

		for {
			handleInput(reader, &buffer, singleCharacterMode)

			if buffer == subject {
				buffer = ""
				break
			} else if singleCharacterMode {
				fmt.Println(subject)
			}
		}
	}
}
