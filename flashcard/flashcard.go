package flashcard

import (
	"fmt"
	"strings"
)

type Responce struct {
	Learner string
	Cards   []Card
}

type Card struct {
	Word    string
	Meaning string
}

type Flashcards struct {
	learner    Learner
	meaning    Meaning
	exporter   Exporter
	lemmatizer Lemmatizer
}

type Learner interface {
	Flashcards(learnerID string) *Responce
	AddFlashcards(learnerID string, flashcards *Responce)
}

type Meaning interface {
	GetMeaning(string) string
}

type Exporter interface {
	Export([]Card) []byte
}

type Lemmatizer interface {
	Lemmatize(string) string
}

type Settings struct {
	Learner    Learner
	Meaning    Meaning
	Exporter   Exporter
	Lemmatizer Lemmatizer
}

func New(settings Settings) (*Flashcards, error) {
	err := check(settings)
	if err != nil {
		return nil, fmt.Errorf("checking settings: %w", err)
	}

	return &Flashcards{
		learner:    settings.Learner,
		meaning:    settings.Meaning,
		exporter:   settings.Exporter,
		lemmatizer: settings.Lemmatizer,
	}, nil
}

func (f Flashcards) CreateFlashCards(learnerID string, text string) Responce {
	flashcards := f.learner.Flashcards(learnerID)
	if flashcards == nil {
		flashcards = &Responce{
			Learner: learnerID,
		}
		f.learner.AddFlashcards(learnerID, flashcards)
	}

	seen := make(map[string]bool)
	for _, card := range flashcards.Cards {
		seen[card.Word] = true
	}

	words := strings.Split(text, " ")

	for _, word := range words {
		lemmtizedWord := f.lemmatizer.Lemmatize(word)
		if !seen[lemmtizedWord] {
			flashcards.Cards = append(flashcards.Cards, Card{
				Word:    lemmtizedWord,
				Meaning: f.meaning.GetMeaning(lemmtizedWord),
			})
			seen[lemmtizedWord] = true
		}
	}

	return *flashcards
}

func (f Flashcards) Export(learner string) []byte {
	learnerCards := f.learner.Flashcards(learner).Cards
	return f.exporter.Export(learnerCards)
}

func check(setting Settings) error {
	if setting.Learner == nil {
		return fmt.Errorf("learner is not defined")
	}
	if setting.Meaning == nil {
		return fmt.Errorf("menaing is not defined")
	}
	if setting.Exporter == nil {
		return fmt.Errorf("exporter is not defined")
	}
	if setting.Lemmatizer == nil {
		return fmt.Errorf("lemmatizer is not defined")
	}
	return nil
}
