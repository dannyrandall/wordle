package wordle

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

type Solver struct {
	wordle       Wordle
	dict         []string
	guessedWords map[string]bool

	potentialLetters []map[rune]bool
	yellowLetters    string
}

func NewSolver(w Wordle, dict []string) (Solver, error) {
	if len(w.Guesses()) > 0 {
		return Solver{}, fmt.Errorf("wordle has already been played")
	}

	// make sure wordle is in dict?
	var correctLenDict []string
	for _, word := range dict {
		if len(word) == w.Len() {
			correctLenDict = append(correctLenDict, word)
		}
	}

	s := Solver{
		wordle:           w,
		dict:             correctLenDict,
		guessedWords:     make(map[string]bool),
		potentialLetters: make([]map[rune]bool, w.Len()),
	}

	for i := range s.potentialLetters {
		s.potentialLetters[i] = map[rune]bool{
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

	s.updateDict()
	return s, nil
}

func (s *Solver) Solve() (string, []Guess) {
	for {
		guess := s.dict[0]
		log.Printf("Remaining Words: %v. Total Guesses: %v. Guess: %q", len(s.dict), len(s.wordle.Guesses()), guess)

		res, err := s.wordle.Guess(guess)
		if err != nil {
			panic(err)
		}
		if res.Correct() {
			return res.Word(), s.wordle.Guesses()
		}

		for i, lr := range res {
			switch lr.Result {
			case NotInWord:
				s.markLetterGrey(rune(lr.Letter[0]))
			case InWord:
				s.markLetterYellow(i, rune(lr.Letter[0]))
			case InWordAndCorrectSpot:
				s.markLetterGreen(i, rune(lr.Letter[0]))
			}
		}

		s.updateDict()

		time.Sleep(1 * time.Second)
	}
}

func (s *Solver) markLetterGrey(r rune) {
	for i := range s.potentialLetters {
		s.potentialLetters[i][r] = false
	}
}

func (s *Solver) markLetterGreen(i int, r rune) {
	s.potentialLetters[i] = map[rune]bool{
		r: true,
	}
}

func (s *Solver) markLetterYellow(i int, r rune) {
	// it can't be in position i
	s.potentialLetters[i][r] = false

	if !strings.ContainsAny(s.yellowLetters, string(r)) {
		s.yellowLetters += string(r)
	}
}

func (s *Solver) updateDict() {
	var validWords []string

	for _, word := range s.dict {
		if s.couldWordWork(word) {
			validWords = append(validWords, word)
		}
	}

	wordScores := make(map[string]int)
	for _, word := range validWords {
		wordScores[word] = s.wordScore(word, validWords)
	}

	sort.Slice(validWords, func(i, j int) bool {
		return wordScores[validWords[i]] > wordScores[validWords[j]]
	})

	s.dict = validWords
}

// a word is more valuable the more letters it shares with other words.
// v0: words that have common letters
// *v1: number of letters overall in common accross all words
// v2: gotta factor in letters that have been guessed already/already know position
func (s *Solver) wordScore(w string, otherWords []string) int {
	score := 0
	for _, word := range otherWords {
		seen := make(map[rune]bool)

		for _, r := range w {
			if seen[r] {
				continue
			}

			seen[r] = true
			if strings.ContainsRune(word, r) {
				score++
			}
		}
	}

	return score
}

// couldWordWork checks that each of the letters in word
// matches with the potential letters for that index
func (s *Solver) couldWordWork(word string) bool {
	if s.guessedWords[word] {
		return false
	}

	for i, r := range word {
		if !s.potentialLetters[i][r] {
			return false
		}
	}

	// make sure it has all the yellow letters
	for _, r := range s.yellowLetters {
		if !strings.ContainsRune(word, r) {
			return false
		}
	}

	return true
}
