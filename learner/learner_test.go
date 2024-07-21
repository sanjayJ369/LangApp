package learner_test

import (
	"testing"

	"github.com/sanjayJ369/LangApp/flashcard"
	"github.com/sanjayJ369/LangApp/learner"
	"github.com/sanjayJ369/LangApp/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestLearner(t *testing.T) {
	t.Parallel()

	t.Run("Persistent Data", func(t *testing.T) {
		t.Parallel()
		fileLoc := testhelper.GetTempFileLoc()
		l1 := learner.New(fileLoc)
		l1Responce := &flashcard.Responce{
			Learner: "tester",
			Cards: []flashcard.Card{
				{
					Word:    "flash",
					Meaning: "card",
				},
				{
					Word:    "Persistent",
					Meaning: "Data",
				},
			},
		}

		l1.AddFlashcards("tester", l1Responce)
		// When Learner restarts app.
		l2 := learner.New(fileLoc)

		// Then they can access their cards.
		assert.Equal(t, l1Responce, l2.Flashcards("tester"))
	})
}
