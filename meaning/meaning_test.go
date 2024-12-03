package meaning_test

import (
	"testing"

	"github.com/sanjayJ369/LangApp/database"
	"github.com/sanjayJ369/LangApp/meaning"
	"github.com/sanjayJ369/LangApp/parser"
	"github.com/sanjayJ369/LangApp/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeaning(t *testing.T) {
	t.Parallel()

	settings, cleanup := validSettings(t)
	t.Cleanup(cleanup)

	meaningGetter := meaning.New(settings)

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

	tmpFile := testhelper.GetTempFileLoc()

	h, err := database.NewSqlite(tmpFile)
	require.NoError(tb, err, "creating sqlite db handler")

	p := parser.New(parser.Settings{
		FileLoc:   "./testfiles/vocabulary.jsonl",
		DBhandler: h,
	})
	require.NoError(tb, p.Parse(), "parsing vocabulary")

	return meaning.Settings{
			DBHandler: h,
		}, func() {
			require.NoError(tb, h.Close(), "closing sqlite db handler")
		}
}
