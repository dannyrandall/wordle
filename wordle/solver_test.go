package wordle

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/matryer/is"
)

func TestSolveEveryWord(t *testing.T) {
	dict, err := readDictionary("../words_alpha.txt")
	if err != nil {
		t.Fatalf("unable to read dictionary: %s", err)
	}

	// sort dict once so this doesn't take forever everytime
	dict = sortDict(dict)

	start := time.Now()
	t.Logf("Solving %v words", len(dict))

	for i, word := range dict {
		t.Logf("Solved %v/%v words", i, len(dict))
		t.Run(word, func(t *testing.T) {
			is := is.New(t)

			game := NewWordle(word)
			solver, err := NewSolver(game, dict)
			is.NoErr(err)
			final, _ := solver.Solve()
			is.Equal(word, final)
		})
	}

	t.Logf("Took %v to solve %v words", time.Since(start), len(dict))
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
		if len(word) == 5 && IsLetters(word) {
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

func sortDict(d []string) []string {
	wordScores := make(map[string]int)
	for _, word := range d {
		wordScores[word] = wordScore(word, d)
	}

	sort.Slice(d, func(i, j int) bool {
		return wordScores[d[i]] > wordScores[d[j]]
	})

	return d
}

// a word is more valuable the more letters it shares with other words.
// v0: words that have common letters
// *v1: number of letters overall in common accross all words
func wordScore(w string, otherWords []string) int {
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
