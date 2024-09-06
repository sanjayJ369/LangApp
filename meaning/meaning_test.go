package meaning_test

import (
	"testing"

	"github.com/sanjayJ369/LangApp/meaning"
	"github.com/stretchr/testify/assert"
)

func TestMeaning(t *testing.T) {
	t.Parallel()

	t.Run("get word meaning", func(t *testing.T) {
		t.Parallel()

		type testCase struct {
			Word    string `field:"word"`
			Meaning string `field:"meaning"`
		}

		testCases := map[string]testCase{
			"abaiser_Ivory_black;_animal_charcoal.":               {"Abaiser", "Ivory black; animal charcoal."},
			"fabaceous_Having_the_nature_of_a_bean;_like_a_bean.": {"fabaceous", "Having the nature of a bean; like a bean."},
		}

		for name, testCase := range testCases {
			testCase := testCase

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				_ = testCase // TODO: Use and remove.
				// When user request the meaning of a <word>.
				got := meaning.New().GetMeaning(testCase.Word)

				// Then they receive it's <meaning>.
				assert.Equal(t, testCase.Meaning, got)
			})
		}
	})
}
