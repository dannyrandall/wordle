package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

var dictionary []string
var realWord = "hello"

type Wordle struct {
	// PotentialLetters is set of runes for every character in the real word. It represents what the letter at each position could still be.
	PotentialLetters []map[rune]bool
	YellowLetters    string
	GuessedWords     map[string]bool

	// RemainingWords is a list of valid words that could be guessed. It is filtered after every guess.
	RemainingWords []string

	// Guesses is the total number of guesses made
	Guesses int
}

func main() {
	dictionary, err := readDictionary("./words_alpha.txt")
	if err != nil {
		log.Fatalf("unable to read dictionary: %s", err.Error())
	}

	var wordle Wordle
	wordle.Reset(realWord, dictionary)

	for {
		guess := wordle.RemainingWords[0]
		log.Printf("Remaining Words: %v. Total Guesses: %v. Guess: %q", len(wordle.RemainingWords), wordle.Guesses, guess)

		match := wordle.Guess(guess)
		if match {
			log.Printf("Yay! The word was %q. Total guesses: %v", realWord, wordle.Guesses)
			return
		}

		wordle.UpdateRemainingWords()
		time.Sleep(1 * time.Second)
	}
}

func (w *Wordle) Reset(realWord string, dictionary []string) {
	w.Guesses = 0
	w.RemainingWords = dictionary
	w.GuessedWords = make(map[string]bool)
	w.PotentialLetters = make([]map[rune]bool, len(realWord))

	for i := range w.PotentialLetters {
		w.PotentialLetters[i] = map[rune]bool{
			'a': true,
			'b': true,
			'c': true,
			'd': true,
			'e': true,
			'f': true,
			'g': true,
			'h': true,
			'i': true,
			'j': true,
			'k': true,
			'l': true,
			'm': true,
			'n': true,
			'o': true,
			'p': true,
			'q': true,
			'r': true,
			's': true,
			't': true,
			'u': true,
			'v': true,
			'w': true,
			'x': true,
			'y': true,
			'z': true,
		}
	}
}

func (w *Wordle) Guess(guess string) bool {
	w.Guesses++
	w.GuessedWords[guess] = true
	match := true

	for i, r := range guess {
		switch {
		case guess[i] == realWord[i]:
			w.markLetterGreen(i, r)
		case strings.ContainsRune(realWord, r):
			w.markLetterYellow(i, r)
			match = false
		default:
			w.markLetterGrey(r)
			match = false
		}
	}

	return match
}

func (w *Wordle) markLetterGrey(r rune) {
	for i := range w.PotentialLetters {
		w.PotentialLetters[i][r] = false
	}
}

func (w *Wordle) markLetterGreen(i int, r rune) {
	w.PotentialLetters[i] = map[rune]bool{
		r: true,
	}
}

func (w *Wordle) markLetterYellow(i int, r rune) {
	// it can't be in position i
	w.PotentialLetters[i][r] = false

	if !strings.ContainsAny(w.YellowLetters, string(r)) {
		w.YellowLetters += string(r)
	}
}

func (w *Wordle) UpdateRemainingWords() {
	var validWords []string

	for _, word := range w.RemainingWords {
		if w.couldWordWork(word) {
			validWords = append(validWords, word)
		}
	}

	w.RemainingWords = validWords
}

// couldWordWork checks that each of the letters in word
// matches with the potential letters for that index
func (w *Wordle) couldWordWork(word string) bool {
	if w.GuessedWords[word] {
		return false
	}

	for i, r := range word {
		if !w.PotentialLetters[i][r] {
			return false
		}
	}

	// make sure it has all the yellow letters
	for _, r := range w.YellowLetters {
		// log.Printf("w.YellowLetters: %v", w.YellowLetters)
		if !strings.ContainsRune(word, r) {
			return false
		}
	}

	return true
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
