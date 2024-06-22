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
	seen := make(map[string]bool)

	for _, word := range words {
		if !seen[word] {
			cards = append(cards, Flashcard{
				Word: word,
			})
			seen[word] = true
		}
	}

	return cards
}
