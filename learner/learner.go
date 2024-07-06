// This is a naive data source inplementation for the Learner object.
package learner

import "github.com/sanjayJ369/LangApp/flashcard"

type Learner map[string]*flashcard.Responce

func New() Learner {
	var l Learner = make(map[string]*flashcard.Responce)

	return l
}

func (l Learner) Flashcards(learnerID string) *flashcard.Responce {
	return l[learnerID]
}

func (l Learner) AddFlashcards(learnerID string, flashcards *flashcard.Responce) {
	l[learnerID] = flashcards
}
