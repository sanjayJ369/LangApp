package generated_test

import (
	"testing"
)

func TestFlashcards(t *testing.T) {
	t.Parallel()

	t.Run("Learner creates flashcards", func(t *testing.T) {
		t.Parallel()

		// When the Learner passes some text.

		// Then the Learner receives flashcards from it.

	})

	t.Run("Learner gets flashcards without duplicates", func(t *testing.T) {
		t.Parallel()

		// When the Learner passes some text with repeated words.

		// Then the Learner receives flashcards without duplicates.

	})

	t.Run("Learner does not get new flashcards from the same text", func(t *testing.T) {
		t.Parallel()

		// When the Learner passes some text.

		// Then the Learner receives flashcards from it.

		// When the Learner passes the same text again.

		// Then the Learner does not receive new flashcards.

	})

	t.Run("Learner can learn flashcards", func(t *testing.T) {
		t.Parallel()

		type testCase struct {
			Guess     string `field:"<guess>"`
			Memorized string `field:"<memorized>"`
		}

		testCases := map[string]testCase{
			"right_yes": {"right", "yes"},
			"wrong_no":  {"wrong", "no"},
		}

		for name, testCase := range testCases {
			testCase := testCase

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				_ = testCase // TODO: Use and remove.
				// When the Learner receives a flashcard.

				// Then the Learner can <guess> the meaning of it.

				// And the flashcard becomes <memorized>.

			})
		}
	})
}
