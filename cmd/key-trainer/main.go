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

func main() {
	word := flag.Bool("w", false, "open in word mode")
	flag.Parse()

	// this little bit of unintuative magic disables input buffering
	// so that the reader can read a single byte immediatly without
	// the need for pressing enter
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").
		Run()

	if *word {
		singleWord()
	} else {
		singleCharacter()
	}
}

func getRandomRune() byte {
	characterLen := len(runes)
	return runes[rand.Intn(characterLen)]
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
			} else if string(input) == "\x7f" {
				fmt.Print("\n\033[1A\033[K")
				buffer = ""
			}
		}
	}
}

func singleCharacter() {
	reader := bufio.NewReader(os.Stdin)
	rand.Seed(time.Now().Unix())

	for {
		char := getRandomRune()
		fmt.Printf("\n%s\n", string(char))

		for {
			input, _ := reader.ReadByte()
			fmt.Println("")
			if char == input {
				break
			}
			fmt.Println(string(char))
		}
	}
}
