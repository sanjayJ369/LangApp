// This is a naive data source inplementation for the Learner object.
package learner

import (
	"encoding/json"
	"os"

	"github.com/sanjayJ369/LangApp/flashcard"
)

type Learner struct {
	loc  string
	data map[string]*flashcard.Responce
}

func New(loc string) Learner {
	file, _ := os.Open(loc)
	var l Learner = Learner{
		loc:  loc,
		data: make(map[string]*flashcard.Responce),
	}
	if file != nil {
		json.NewDecoder(file).Decode(&l.data)
	}
	return l
}

func (l Learner) Flashcards(learnerID string) *flashcard.Responce {
	return l.data[learnerID]
}

func (l Learner) AddFlashcards(learnerID string, flashcards *flashcard.Responce) {
	l.data[learnerID] = flashcards
	data, _ := json.Marshal(l.data)
	os.WriteFile(l.loc, data, 0666)
}
