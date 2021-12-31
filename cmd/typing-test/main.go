package main

import (
	"context"
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

	ctx, cancel := context.WithCancel(context.Background())
	runner := internal.NewRunner(wordList, config)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		cancel()

		<-time.After(time.Second)
		os.Exit(1)
	}()

	runner.Run(ctx)
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

	flag.UintVar(&config.TotalWords, "wc", 0, "Total words to be displayed before the run is complete")
	flag.StringVar(&config.Contains, "c", "", "Only use words that contain the given substring")
	flag.StringVar(&config.ContainsAny, "ca", "", "Only use words that contain any of the given substring")
	flag.UintVar(&config.TimeLimit, "t", 0, "Set a time limit for the test")

	flag.Parse()

	if flag.NArg() != 0 {
		wordList = flag.Arg(0)
	}

	return wordList, config
}
