package flashcard

import (
	"strings"
)

var learnerCards = make(map[string]Flashcards)

type Flashcards struct {
	Learner string
	Cards   []Flashcard
}

type Flashcard struct {
	Word string
}

func CreateFlashCards(learner string, text string) Flashcards {
	container, ok := learnerCards[learner]
	if !ok {
		learnerCards[learner] = Flashcards{
			Learner: learner,
		}
	}

	seen := make(map[string]bool)
	for _, card := range container.Cards {
		seen[card.Word] = true
	}

	words := strings.Split(text, " ")

	for _, word := range words {
		if !seen[word] {
			container.Cards = append(container.Cards, Flashcard{
				Word: word,
			})
			seen[word] = true
		}
	}

	learnerCards[learner] = container

	return container
}
