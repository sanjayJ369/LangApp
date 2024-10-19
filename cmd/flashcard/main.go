package main

import (
	"log"

	"github.com/sanjayJ369/LangApp/exporter"
	"github.com/sanjayJ369/LangApp/flashcard"
	"github.com/sanjayJ369/LangApp/learner"
	"github.com/sanjayJ369/LangApp/lemmatizer"
	"github.com/sanjayJ369/LangApp/meaning"
)

func main() {
	meaningSetting := meaning.Settings{}

	settings := flashcard.Settings{
		Learner:    learner.New("../../assets/learner.db"),
		Meaning:    meaning.New(meaningSetting),
		Exporter:   exporter.New(),
		Lemmatizer: lemmatizer.New(),
	}
	_, err := flashcard.New(settings)
	if err != nil {
		log.Fatalf("creating flashcards: %s", err)
	}
}
