package meaning_test

import (
	"errors"
	"testing"

	"github.com/sanjayJ369/LangApp/meaning"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			"abaiser_Ivory_black;_animal_charcoal.,":               {"abaiser", "Ivory black; animal charcoal.,"},
			"fabaceous_Having_the_nature_of_a_bean;_like_a_bean.,": {"fabaceous", "Having the nature of a bean; like a bean.,"},
		}

		for name, testCase := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				var meaningGetter fakeMeaningGetter
				meaningGetter.add(testCase.Word, testCase.Meaning)

				meaningH, err := meaning.New(meaning.Settings{
					GetMeaning: &meaningGetter,
				})
				require.NoError(t, err, "creating meaginig handler")

				// When user request the meaning of a <word>.
				got := meaningH.GetMeaning(testCase.Word)

				// Then they receive it's <meaning>.
				assert.Equal(t, testCase.Meaning, got)
			})
		}
	})
}

type fakeMeaningGetter struct {
	meanings [][]string
}

func (f *fakeMeaningGetter) add(key, val string) {
	f.meanings = append(f.meanings, []string{key, val})
}

func (f *fakeMeaningGetter) Get(key string) (string, error) {
	for _, v := range f.meanings {
		if v[0] == key {
			return v[1], nil
		}
	}

	return "", errors.New("no meaning")
}
