package wordle

import (
	"fmt"
)

type Wordle struct {
	word     string
	guesses  []Guess
	finished bool
}

func NewWordle(word string) Wordle {
	return Wordle{
		word: word,
	}
}

func (w *Wordle) Guess(word string) (Guess, error) {
	switch {
	case len(w.word) != len(word):
		return Guess{}, fmt.Errorf("guessed word length does not match the wordle")
	case w.finished:
		return Guess{}, fmt.Errorf("already guessed the wordle")
	}

	guess := make(Guess, len(w.word))
	for i := range word {
		guess[i].Letter = string(word[i])
		guess[i].Result = NotInWord // default to grey
	}

	used := make([]bool, len(w.word))

	// determine which letters are in correct spot
	for i := range word {
		if w.word[i] == word[i] {
			used[i] = true
			guess[i].Result = InWordAndCorrectSpot
		}
	}

	freq := make(map[string]int)
	for i := range w.word {
		if !used[i] {
			freq[string(w.word[i])]++
		}
	}

	// determine which letters can be marked yellow
	for i := range word {
		if guess[i].Result == InWordAndCorrectSpot {
			continue
		}

		if freq[string(word[i])] > 0 {
			guess[i].Result = InWordButNotSpot
			freq[string(word[i])]--
		}
	}

	w.guesses = append(w.guesses, guess)
	if guess.Correct() {
		w.finished = true
	}

	return guess, nil
}

func (w *Wordle) Guesses() []Guess {
	return w.guesses
}

func (w *Wordle) Finished() bool {
	return w.finished
}

func (w *Wordle) Len() int {
	return len(w.word)
}
