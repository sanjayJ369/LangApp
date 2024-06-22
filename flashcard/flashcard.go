package flashcard

import (
	"strings"
)

type Flashcard struct {
	Word string
}

func CreateFlashCards(text string) []Flashcard {
	words := strings.Split(text, " ")

	cards := make([]Flashcard, 0, len(words))

	for _, word := range words {
		cards = append(cards, Flashcard{
			Word: word,
		})
	}

	return cards
}
