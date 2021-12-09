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
	specialCharacter := flag.Bool("s", false, "open in function special character mode")
	contains := flag.String("c", "", "open in word mode and only user words containing the given substring")
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
		singleWordMode()
	} else if *specialCharacter {
		specialCharacterMode()
	} else {
		singleCharacterMode()
	}
}

// getRandomRune picks a character at random from the geven charset
func getRandomRune(charset string) byte {
	characterLen := len(charset)
	return charset[rand.Intn(characterLen)]
}

// handleInput from the users keyboard
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

// singleWordMode prints and waits for a single word
func singleWordMode() {
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

// containsMode is the same a word mode but only includes words with the given substring
func containsMode(substring string) {
	words, err := ioutil.ReadFile("./words.txt")
	if err != nil {
		panic(err)
	}

	rawWordList := strings.Split(string(words), "\n")
	wordList := []string{}
	for _, word := range rawWordList {
		if strings.Contains(word, substring) {
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

// singleCharacterMode displays and waits for a single character
func singleCharacterMode() {
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

// specialCharacterMode picks a single special character and waits for input
func specialCharacterMode() {
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
