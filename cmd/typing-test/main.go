package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/indeedhat/key-trainer/internal"
)

const DefaultWordlist = "default"

func main() {
	wordList, config := parseInput()

	// this little bit of unintuative magic disables input buffering
	// so that the reader can read a single byte immediatly without
	// the need for pressing enter
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").
		Run()

	rand.Seed(time.Now().Unix())

	runner := internal.NewRunner(config)

	handleInterupt(runner)
	runner.Run(wordList)
}

// parseInput from the command line to generate config
func parseInput() (string, internal.RunnerConfig) {
	wordList := DefaultWordlist
	config := internal.RunnerConfig{}

	wordLists := internal.FindWordLists()
	wordListsString := ""
	for _, key := range wordLists {
		wordListsString += fmt.Sprintf("  %s\n", key)
	}

	flag.Usage = func() {
		fmt.Printf(
			`Typing Test
A simple tool to help me not suck

Usage:
  ./typing-test [options] [wordlist]

WORDLISTS:
%s
OPTIONS:
`,
			wordListsString,
		)

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

// handleInterupt from the system (ctrl+c)
func handleInterupt(run *internal.Runner) {
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		run.DisplayStatusScreen()
		os.Exit(1)
	}()
}
