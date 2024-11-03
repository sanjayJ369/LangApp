package meaning_test

import (
	"testing"

	"github.com/sanjayJ369/LangApp/database"
	"github.com/sanjayJ369/LangApp/meaning"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeaning(t *testing.T) {

	settings, cleanup := validSettings(t)
	t.Cleanup(cleanup)
	meaningGetter := meaning.New(settings)

	t.Run("get word meaning", func(t *testing.T) {

		type testCase struct {
			Word    string `field:"word"`
			Meaning string `field:"meaning"`
		}

		testCases := map[string]testCase{
			"abaiser_Ivory_black;_animal_charcoal.,":               {"abaiser", "Ivory black; animal charcoal.,"},
			"fabaceous_Having_the_nature_of_a_bean;_like_a_bean.,": {"fabaceous", "Having the nature of a bean; like a bean.,"},
		}

		for name, testCase := range testCases {
			testCase := testCase

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				_ = testCase // TODO: Use and remove.
				// When user request the meaning of a <word>.
				got := meaningGetter.GetMeaning(testCase.Word)

				// Then they receive it's <meaning>.
				assert.Equal(t, testCase.Meaning, got)
			})
		}
	})
}

func validSettings(tb testing.TB) (settings meaning.Settings, cleanup func()) {
	tb.Helper()
	handler, err := database.NewSqlite("../assets/meaning.db")
	require.NoError(tb, err, "creating sqlite3")
	return meaning.Settings{
			DBHandler: handler,
		}, func() {
			require.NoError(tb, handler.Close(), "closing sqlite3")
		}
}
