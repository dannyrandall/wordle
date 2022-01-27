package wordle

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

type Solver struct {
	game         Wordle
	dict         []string
	guessedWords map[string]bool
}

func NewSolver(game Wordle, dict []string) (Solver, error) {
	if len(game.Guesses()) > 0 {
		return Solver{}, fmt.Errorf("wordle has already been played")
	}

	// TODO make sure wordle is in dict?
	// filter out words that don't match wordle length from dict
	var d []string
	for _, word := range dict {
		if len(word) == game.Len() {
			d = append(d, word)
		}
	}

	s := Solver{
		game:         game,
		dict:         d,
		guessedWords: make(map[string]bool),
	}

	s.sortDict()
	return s, nil
}

func (s *Solver) Solve() (string, []Guess) {
	for {
		guess := s.dict[0]
		log.Printf("Remaining Words: %v. Total Guesses: %v. Guess: %q", len(s.dict), len(s.game.Guesses()), guess)

		res, err := s.game.Guess(guess)
		if err != nil {
			panic(err) // TODO handle
		}
		if res.Correct() {
			return res.Word(), s.game.Guesses()
		}

		lettersInWordButNotSpots := make(map[string][]int)

		for i, lr := range res {
			switch lr.Result {
			case NotInWord:
				s.filterNotInWord(lr.Letter)
			case InWordButNotSpot:
				lettersInWordButNotSpots[lr.Letter] = append(lettersInWordButNotSpots[lr.Letter], i)
			case InWordAndCorrectSpot:
				s.filterInWordAndCorrectSpot(lr.Letter, i)
			}
		}

		for letter, spots := range lettersInWordButNotSpots {
			s.filterInWordButNotSpots(letter, spots)
		}

		s.sortDict()
	}
}

// filterNotInWord removes all words from dict that contain letter.
func (s *Solver) filterNotInWord(letter string) {
	var d []string
	for _, word := range s.dict {
		if !strings.Contains(word, letter) {
			d = append(d, word)
		}
	}

	s.dict = d
}

// filterInWordAndCorrectSpot removes all words from dict that don't contain letter at spot.
func (s *Solver) filterInWordAndCorrectSpot(letter string, spot int) {
	var d []string
	for _, word := range s.dict {
		if string(word[spot]) == letter {
			d = append(d, word)
		}
	}

	s.dict = d
}

// filterInWordAndCorrectSpot removes all words from dict that have letter in any of the given spots and don't have at least len(spots) of letter in them.
func (s *Solver) filterInWordButNotSpots(letter string, spots []int) {
	var d []string
	for _, word := range s.dict {
		if !letterInSpots(word, letter, spots) && strings.Count(word, letter) >= len(spots) {
			d = append(d, word)
		}
	}

	s.dict = d
}

// letterInSpots returns true if word has letter at any of the given spots
func letterInSpots(word string, letter string, spots []int) bool {
	for _, spot := range spots {
		if string(word[spot]) == letter {
			return true
		}
	}

	return false
}

func (s *Solver) sortDict() {
	wordScores := make(map[string]int)
	for _, word := range s.dict {
		wordScores[word] = s.wordScore(word, s.dict)
	}

	sort.Slice(s.dict, func(i, j int) bool {
		return wordScores[s.dict[i]] > wordScores[s.dict[j]]
	})
}

// a word is more valuable the more letters it shares with other words.
// v0: words that have common letters
// *v1: number of letters overall in common accross all words
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
