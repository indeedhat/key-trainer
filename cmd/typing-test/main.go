package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os/exec"
	"time"

	"github.com/indeedhat/key-trainer/internal"
)

const DefaultWordlist = "default"

func main() {
	wordList, config := parseInput()

	fmt.Print(wordList)
	time.Sleep(time.Second)

	// this little bit of unintuative magic disables input buffering
	// so that the reader can read a single byte immediatly without
	// the need for pressing enter
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").
		Run()

	rand.Seed(time.Now().Unix())

	runner := internal.NewRunner(config)
	runner.Run(wordList)
}

// parseInput from the command line to generate config
func parseInput() (string, internal.RunnerConfig) {
	wordList := DefaultWordlist
	config := internal.RunnerConfig{}

	flag.Usage = func() {
		fmt.Print("Typing Test\n\n")
		fmt.Println("A simple tool to help me not suck")

		flag.PrintDefaults()
	}

	flag.UintVar(&config.TotalWords, "t", 0, "Total words to be displayed before the run is complete")
	flag.StringVar(&config.Contains, "c", "", "Only use words that contain the given substring")
	flag.Parse()

	if flag.NArg() != 0 {
		wordList = flag.Arg(0)
	}

	return wordList, config
}
