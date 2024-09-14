package flashcards_test

import (
	"testing"

	"github.com/sanjayJ369/LangApp/exporter"
	"github.com/sanjayJ369/LangApp/flashcard"
	"github.com/sanjayJ369/LangApp/learner"
	"github.com/sanjayJ369/LangApp/lemmatizer"
	"github.com/sanjayJ369/LangApp/meaning"
	"github.com/sanjayJ369/LangApp/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlashcardsUsage(t *testing.T) {
	t.Parallel()

	t.Run("User creates flashcards", func(t *testing.T) {
		t.Parallel()
		settings := flashcard.Settings{
			Learner:    learner.New(testhelper.GetTempFileLoc()),
			Meaning:    meaning.New(meaning.Settings{FileLoc: "./assets/kaikki.org-dictionary-English.jsonl"}),
			Exporter:   exporter.New(),
			Lemmatizer: lemmatizer.New(),
		}
		cards, err := flashcard.New(settings)
		require.NoError(t, err, "checking settings")

		// When User wants to create flashcards.
		someText := "Abaiser fabaceous"

		// Then they can do it.
		res := cards.CreateFlashCards("learner", someText)
		assert.NotNil(t, res)
	})
}
