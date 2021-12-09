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
var layerRunes = "1234567890!@#$%^&*()-=[]_+{};:'/\"|><.,~`"

func main() {
	word := flag.Bool("w", false, "open in word mode")
	layers := flag.Bool("l", false, "open in function layer mode")
	contains := flag.String("c", "", "open in word mode and only user words containing the given character")
	flag.Parse()

	// this little bit of unintuative magic disables input buffering
	// so that the reader can read a single byte immediatly without
	// the need for pressing enter
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").
		Run()

		// and this one will stop the characters you type from appearing in
		// the terminal
	exec.Command("stty", "-F", "/dev/tty", "-echo", "min", "1").
		Run()

	if *contains != "" {
		containsMode(*contains)
	} else if *word {
		singleWord()
	} else if *layers {
		layerMode()
	} else {
		singleCharacter()
	}
}

func getRandomRune(charset string) byte {
	characterLen := len(charset)
	return charset[rand.Intn(characterLen)]
}

func handleInput(reader *bufio.Reader, buffer *string) {
	input, _ := reader.ReadByte()
	char := string(input)

	if char == "\x7f" {
		fmt.Print("\n\033[1A\033[K")
		*buffer = ""
	} else if char != " " && char != "\n" {
		fmt.Print(string(input))
		*buffer += char
	}
}

func singleWord() {
	words, err := ioutil.ReadFile("./words.txt")
	if err != nil {
		panic(err)
	}

	wordList := strings.Split(string(words), "\n")
	wordCount := len(wordList)

	reader := bufio.NewReader(os.Stdin)
	rand.Seed(time.Now().Unix())

	buffer := ""
	for {
		subject := wordList[rand.Intn(wordCount)]
		fmt.Printf("\n%s\n", string(subject))

		for {
			input, _ := reader.ReadByte()
			buffer += string(input)

			if buffer == subject {
				buffer = ""
				break
			} else if string(input) == "\x7f" || string(input) == "\n" {
				fmt.Print("\n\033[1A\033[K")
				buffer = ""
			} else {
				fmt.Print(string(input))
			}
		}
	}
}

func containsMode(letter string) {
	words, err := ioutil.ReadFile("./words.txt")
	if err != nil {
		panic(err)
	}

	rawWordList := strings.Split(string(words), "\n")
	wordList := []string{}
	for _, word := range rawWordList {
		if strings.Contains(word, letter) {
			wordList = append(wordList, word)
		}
	}

	wordCount := len(wordList)

	reader := bufio.NewReader(os.Stdin)
	rand.Seed(time.Now().Unix())

	buffer := ""
	for {
		subject := wordList[rand.Intn(wordCount)]
		fmt.Printf("\n%s\n", string(subject))

		for {
			handleInput(reader, &buffer)

			if buffer == subject {
				buffer = ""
				break
			}
		}
	}
}

func singleCharacter() {
	reader := bufio.NewReader(os.Stdin)
	rand.Seed(time.Now().Unix())

	for {
		char := getRandomRune(runes)
		fmt.Printf("\n%s\n", string(char))

		for {
			input, _ := reader.ReadByte()
			fmt.Println(string(input))
			if char == input {
				break
			}
			fmt.Println(string(char))
		}
	}
}

func layerMode() {
	reader := bufio.NewReader(os.Stdin)
	rand.Seed(time.Now().Unix())

	for {
		char := getRandomRune(layerRunes)
		fmt.Printf("\n%s\n", string(char))

		for {
			input, _ := reader.ReadByte()
			fmt.Println(string(input))
			if char == input {
				break
			}
			fmt.Println(string(char))
		}
	}
}
