package flashcard

import (
	"strings"
)

type Flashcards struct {
	Learner string
	Cards   []Flashcard
}

type Flashcard struct {
	Word string
}

type Learner interface {
	Flashcards(learnerID string) *Flashcards
	AddFlashcards(learnerID string, flashcards *Flashcards)
}

func CreateFlashCards(learner Learner, learnerID string, text string) Flashcards {
	flashcards := learner.Flashcards(learnerID)
	if flashcards == nil {
		flashcards = &Flashcards{
			Learner: learnerID,
		}
		learner.AddFlashcards(learnerID, flashcards)
	}

	seen := make(map[string]bool)
	for _, card := range flashcards.Cards {
		seen[card.Word] = true
	}

	words := strings.Split(text, " ")

	for _, word := range words {
		if !seen[word] {
			flashcards.Cards = append(flashcards.Cards, Flashcard{
				Word: word,
			})
			seen[word] = true
		}
	}

	return *flashcards
}
