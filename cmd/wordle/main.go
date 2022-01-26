package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"

	"github.com/dannyrandall/wordle/wordle"
)

var realWord = "danny"

func main() {
	dict, err := readDictionary("../../words_alpha.txt")
	if err != nil {
		log.Fatalf("unable to read dictionary: %s", err)
	}

	game := wordle.NewWordle(realWord)
	solver, err := wordle.NewSolver(game, dict)
	if err != nil {
		log.Fatalf("unable to create solver: %s", err)
	}

	wordle, guesses := solver.Solve()
	log.Printf("Yay! Got the final word in %v guesses: %v", len(guesses), wordle)
}

func readDictionary(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer f.Close()

	var words []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) == len(realWord) && IsLetters(word) {
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("something went wrong with the scanner: %w", err)
	}

	return words, nil
}

func IsLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}
