package flashcard_test

import (
	"testing"

	"github.com/sanjayJ369/LangApp/flashcard"
	"github.com/sanjayJ369/LangApp/learner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlashcards(t *testing.T) {
	t.Parallel()

	t.Run("Learner creates flashcards", func(t *testing.T) {
		t.Parallel()

		l := learner.New()

		// When the Learner passes some text.
		container := flashcard.CreateFlashCards(l, "learner", "some text")

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

		l := learner.New()

		// When the Learner passes some text with repeated words.
		container := flashcard.CreateFlashCards(l, "learner", "some text text")

		// Then the Learner receives flashcards without duplicates.
		require.Len(t, container.Cards, 2)

		assert.Equal(t, "some", container.Cards[0].Word)
		assert.Equal(t, "text", container.Cards[1].Word)
	})

	t.Run("Learner does not get new flashcards from the same text", func(t *testing.T) {
		t.Parallel()

		const (
			learnerID = "learner"
			someText  = "some text text"
		)

		l := learner.New()

		// When the Learner passes some text.
		cards1 := flashcard.CreateFlashCards(l, learnerID, someText)

		// Then the Learner receives flashcards from it.
		require.NotEmpty(t, cards1)

		// When the Learner passes the same text again.
		cards2 := flashcard.CreateFlashCards(l, learnerID, someText)

		// Then the Learner does not receive new flashcards.
		assert.Equal(t, cards1, cards2)
	})

	t.Run("Multiple Learner can create flashcards", func(t *testing.T) {
		t.Parallel()

		l := learner.New()

		// When Learner Sanjay creates flashcards.
		sanjayFlashcards := flashcard.CreateFlashCards(l, "Sanjay", "sanjay")

		// When Learner Dima creates flashcards.
		dimaFlashcards := flashcard.CreateFlashCards(l, "Dima", "dima")

		// And Dima does not see Sanjay flashcards.
		assert.NotContains(t, dimaFlashcards.Cards, flashcard.Flashcard{Word: "sanjay"})

		// And Sanjay does not see Dima flashcards.
		assert.NotContains(t, sanjayFlashcards.Cards, flashcard.Flashcard{Word: "dima"})
	})

	t.Run("Flashcards contain word along with it's meaning", func(t *testing.T) {
		t.Parallel()

		// When the Learner passes some text.

		// Then they receive flashcards.

		// And each flashcards has meaning of the word.

	})

	t.Run("Learner can export flashcards to Anki", func(t *testing.T) {
		t.Parallel()

		// When the Learner creates flashcards.

		// Then they can export them to Anki.

	})
}
