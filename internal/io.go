package internal

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func loadFromFile(filename, contains string) []string {
	words, err := ioutil.ReadFile(fmt.Sprint("wordlists/", filename, ".txt"))
	if err != nil {
		panic(err)
	}

	wordList := strings.Split(string(words), "\n")
	if len(contains) == 0 {
		return wordList
	}

	var filteredWords []string
	for _, word := range wordList {
		if strings.Contains(word, contains) {
			filteredWords = append(filteredWords, word)
		}
	}

	return filteredWords
}
