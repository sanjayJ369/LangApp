// This is a naive data source inplementation for the Learner object.
package learner

import "github.com/sanjayJ369/LangApp/flashcard"

type Learner map[string]*flashcard.Flashcards

func New() Learner {
	var l Learner = make(map[string]*flashcard.Flashcards)

	return l
}

func (l Learner) Flashcards(learnerID string) *flashcard.Flashcards {
	return l[learnerID]
}

func (l Learner) AddFlashcards(learnerID string, flashcards *flashcard.Flashcards) {
	l[learnerID] = flashcards
}
