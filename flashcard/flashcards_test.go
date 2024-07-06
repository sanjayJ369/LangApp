package flashcard_test

import (
	"testing"

	"github.com/sanjayJ369/LangApp/exporter"
	"github.com/sanjayJ369/LangApp/flashcard"
	"github.com/sanjayJ369/LangApp/learner"
	"github.com/sanjayJ369/LangApp/meaning"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlashcards(t *testing.T) {
	t.Parallel()

	t.Run("Learner creates flashcards", func(t *testing.T) {
		t.Parallel()
		f, _ := flashcard.New(validSettings())
		// When the Learner passes some text.
		container := f.CreateFlashCards("learner", "some text")

		// Then the Learner receives flashcards from it.
		require.Len(t, container.Cards, 2)

		// And the flashcards contain words from the text.
		assert.Equal(t, "some", container.Cards[0].Word)
		assert.Equal(t, "text", container.Cards[1].Word)

		// And the Learner receives the flashcards they owns.
		assert.Equal(t, "learner", container.Learner)
	})

	t.Run("Learner gets flashcards without duplicates", func(t *testing.T) {
		t.Parallel()
		f, _ := flashcard.New(validSettings())
		// When the Learner passes some text with repeated words.
		container := f.CreateFlashCards("learner", "some text text")

		// Then the Learner receives flashcards without duplicates.
		require.Len(t, container.Cards, 2)

		assert.Equal(t, "some", container.Cards[0].Word)
		assert.Equal(t, "text", container.Cards[1].Word)
	})

	t.Run("Learner does not get new flashcards from the same text", func(t *testing.T) {
		t.Parallel()
		f, _ := flashcard.New(validSettings())
		const (
			learnerID = "learner"
			someText  = "some text text"
		)

		// When the Learner passes some text.
		cards1 := f.CreateFlashCards(learnerID, someText)

		// Then the Learner receives flashcards from it.
		require.NotEmpty(t, cards1)

		// When the Learner passes the same text again.
		cards2 := f.CreateFlashCards(learnerID, someText)

		// Then the Learner does not receive new flashcards.
		assert.Equal(t, cards1, cards2)
	})

	t.Run("Multiple Learner can create flashcards", func(t *testing.T) {
		t.Parallel()
		f, _ := flashcard.New(validSettings())
		// When Learner Sanjay creates flashcards.
		sanjayFlashcards := f.CreateFlashCards("Sanjay", "sanjay")

		// When Learner Dima creates flashcards.
		dimaFlashcards := f.CreateFlashCards("Dima", "dima")

		// And Dima does not see Sanjay flashcards.
		assert.NotContains(t, dimaFlashcards.Cards, flashcard.Card{Word: "sanjay"})

		// And Sanjay does not see Dima flashcards.
		assert.NotContains(t, sanjayFlashcards.Cards, flashcard.Card{Word: "dima"})
	})

	t.Run("Flashcards contain word along with it's meaning", func(t *testing.T) {
		t.Parallel()
		f, _ := flashcard.New(validSettings())
		// When the Learner passes some text.
		learnerFlashcards := f.CreateFlashCards("Learner", "test sentence")

		// Then they receive flashcards.
		assert.NotEmpty(t, learnerFlashcards.Cards)

		// And each flashcards has meaning of the word.
		for _, card := range learnerFlashcards.Cards {
			assert.NotEmpty(t, card.Meaning, card.Word+": is empty")
		}
	})

	t.Run("Learner can export flashcards", func(t *testing.T) {
		t.Parallel()
		f, _ := flashcard.New(validSettings())

		// When the Learner creates flashcards.
		f.CreateFlashCards("Learner", "test sentence")
		// Then they can export them.
		assert.NotEmpty(t, f.Export("Learner"))
	})

}

func validSettings() flashcard.Settings {
	return flashcard.Settings{
		Learner:  learner.New(),
		Meaning:  meaning.New(),
		Exporter: exporter.New(),
	}
}
