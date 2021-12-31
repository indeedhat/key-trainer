package internal

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// loadFromFile will return the word list form the given key
func loadFromFile(filename string, conf RunnerConfig) []string {
	words, err := ioutil.ReadFile(fmt.Sprint("wordlists/", filename, ".txt"))
	if err != nil {
		panic(err)
	}

	wordList := strings.Split(string(words), "\n")
	var filteredWords []string

	for _, word := range wordList {
		if word == "" {
			continue
		}

		if conf.Contains == "" || conf.ContainsAny == "" { 
			filteredWords = append(filteredWords, word)
		} else if strings.Contains(word, conf.Contains) || strings.ContainsAny(word, conf.ContainsAny ) {
			filteredWords = append(filteredWords, word)
        }
	}

	return filteredWords
}

// findWordLists available to the application
func FindWordLists() (wordLists []string) {
	stats, err := ioutil.ReadDir("wordlists/")

	if err != nil {
		return
	}

	for _, stat := range stats {
		if stat.IsDir() || !strings.HasSuffix(stat.Name(), ".txt") {
			continue
		}

		wordLists = append(wordLists, stat.Name()[:len(stat.Name())-4])
	}

	return
}
