package wordle

import (
	"fmt"
	"strings"
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
	var guess Guess

	switch {
	case len(w.word) != len(word):
		return guess, fmt.Errorf("guessed word length does not match the wordle")
	case w.finished:
		return guess, fmt.Errorf("already guessed the wordle")
	}

	for i := range word {
		lr := LetterResult{
			Letter: string(word[i]),
		}

		switch {
		case w.word[i] == word[i]:
			lr.Result = InWordAndCorrectSpot
		case strings.Contains(w.word, string(word[i])):
			lr.Result = InWord
		default:
			lr.Result = NotInWord
		}

		guess = append(guess, lr)
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
