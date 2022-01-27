package wordle

type Guess []LetterResult

type LetterResult struct {
	Letter string
	Result Result
}

type Result string

const (
	NotInWord            Result = "NotInWord"
	InWordButNotSpot     Result = "InWordButNotSpot"
	InWordAndCorrectSpot Result = "InWordAndCorrectSpot"
)

func (g Guess) Word() string {
	var word string
	for _, r := range g {
		word += r.Letter
	}
	return word
}

func (g Guess) Correct() bool {
	if len(g) == 0 {
		return false
	}

	for _, r := range g {
		if r.Result != InWordAndCorrectSpot {
			return false
		}
	}

	return true
}
