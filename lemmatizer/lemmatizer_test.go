package lemmatizer_test

import (
	"testing"

	"github.com/sanjayJ369/LangApp/lemmatizer"
	"github.com/stretchr/testify/assert"
)

func TestLemmatizer(t *testing.T) {
	t.Parallel()

	t.Run("Lemmatizerse Word", func(t *testing.T) {
		t.Parallel()

		type testCase struct {
			Word string `field:"word"`
			Root string `field:"root"`
		}

		testCases := map[string]testCase{
			"Abducting_abduct": {"Abducting", "abduct"},
			"Racing_race":      {"Racing", "race"},
		}

		for name, testCase := range testCases {
			testCase := testCase

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				_ = testCase // TODO: Use and remove.
				// When <word> is lemmatized.
				got := lemmatizer.Lemmatize(testCase.Word)

				// Then it becomes it's <root> word.
				assert.Equal(t, testCase.Root, got)
			})
		}
	})
}
