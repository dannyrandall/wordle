package wordle

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
)

type wordleTest struct {
	wordle  string
	guesses []wordleTestGuess
}

type wordleTestGuess struct {
	word   string
	result Guess
}

var wordleTests = []wordleTest{
	{
		wordle: "danny",
		guesses: []wordleTestGuess{
			{
				word:   "nnnnn",
				result: Guess{{"n", NotInWord}, {"n", NotInWord}, {"n", InWordAndCorrectSpot}, {"n", InWordAndCorrectSpot}, {"n", NotInWord}},
			},
		},
	},
	// tests from https://nerdschalk.com/wordle-same-letter-twice-rules-explained-how-does-it-work/
	{
		wordle: "abbey",
		guesses: []wordleTestGuess{
			{
				word:   "algae",
				result: Guess{{"a", InWordAndCorrectSpot}, {"l", NotInWord}, {"g", NotInWord}, {"a", NotInWord}, {"e", InWordButNotSpot}},
			},
			{
				word:   "keeps",
				result: Guess{{"k", NotInWord}, {"e", InWordButNotSpot}, {"e", NotInWord}, {"p", NotInWord}, {"s", NotInWord}},
			},
			{
				word:   "orbit",
				result: Guess{{"o", NotInWord}, {"r", NotInWord}, {"b", InWordAndCorrectSpot}, {"i", NotInWord}, {"t", NotInWord}},
			},
			{
				word:   "orbbt",
				result: Guess{{"o", NotInWord}, {"r", NotInWord}, {"b", InWordAndCorrectSpot}, {"b", InWordButNotSpot}, {"t", NotInWord}},
			},
			{
				word:   "abate",
				result: Guess{{"a", InWordAndCorrectSpot}, {"b", InWordAndCorrectSpot}, {"a", NotInWord}, {"t", NotInWord}, {"e", InWordButNotSpot}},
			},
			{
				word:   "abbte",
				result: Guess{{"a", InWordAndCorrectSpot}, {"b", InWordAndCorrectSpot}, {"b", InWordAndCorrectSpot}, {"t", NotInWord}, {"e", InWordButNotSpot}},
			},
			{
				word:   "abbey",
				result: Guess{{"a", InWordAndCorrectSpot}, {"b", InWordAndCorrectSpot}, {"b", InWordAndCorrectSpot}, {"e", InWordAndCorrectSpot}, {"y", InWordAndCorrectSpot}},
			},
		},
	},
	{
		wordle: "abbey",
		guesses: []wordleTestGuess{
			{
				word:   "opens",
				result: Guess{{"o", NotInWord}, {"p", NotInWord}, {"e", InWordButNotSpot}, {"n", NotInWord}, {"s", NotInWord}},
			},
			{
				word:   "babes",
				result: Guess{{"b", InWordButNotSpot}, {"a", InWordButNotSpot}, {"b", InWordAndCorrectSpot}, {"e", InWordAndCorrectSpot}, {"s", NotInWord}},
			},
			{
				word:   "kebab",
				result: Guess{{"k", NotInWord}, {"e", InWordButNotSpot}, {"b", InWordAndCorrectSpot}, {"a", InWordButNotSpot}, {"b", InWordButNotSpot}},
			},
			{
				word:   "abyss",
				result: Guess{{"a", InWordAndCorrectSpot}, {"b", InWordAndCorrectSpot}, {"y", InWordButNotSpot}, {"s", NotInWord}, {"s", NotInWord}},
			},
			{
				word:   "abbey",
				result: Guess{{"a", InWordAndCorrectSpot}, {"b", InWordAndCorrectSpot}, {"b", InWordAndCorrectSpot}, {"e", InWordAndCorrectSpot}, {"y", InWordAndCorrectSpot}},
			},
		},
	},
}

func TestWordle(t *testing.T) {
	for i, tt := range wordleTests {
		t.Run(fmt.Sprintf("%v/%q", i, tt.wordle), func(t *testing.T) {
			is := is.New(t)

			game := NewWordle(tt.wordle)
			for _, guess := range tt.guesses {
				res, err := game.Guess(guess.word)
				is.NoErr(err)
				is.Equal(guess.result, res)
			}
		})
	}
}
