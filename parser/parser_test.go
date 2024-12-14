package parser_test

import (
	"os"
	"testing"

	"github.com/sanjayJ369/LangApp/database"
	"github.com/sanjayJ369/LangApp/parser"
	"github.com/sanjayJ369/LangApp/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Parallel()

	settings := validSettings(t)

	p, err := parser.New(settings)
	require.NoError(t, err, "creating parser")

	require.NoError(t, p.Parse(), "parsing")

	meaning, err := settings.DBhandler.Get("abaiser")
	require.NoError(t, err, "getting meaning")

	assert.Equal(t, "Ivory black; animal charcoal.,", meaning)
}

func validSettings(tb testing.TB) parser.Settings {
	tb.Helper()

	name := testhelper.GetTempFileLoc()

	handler, err := database.NewBadger(name)
	require.NoError(tb, err, "creating handler")

	tb.Cleanup(func() {
		require.NoError(tb, handler.Close(), "closing handler")
	})

	f, err := os.Open("./testfiles/vocabulary.jsonl")
	require.NoError(tb, err, "opening vocabulary file")

	tb.Cleanup(func() {
		require.NoError(tb, f.Close(), "closing vocabulary file")
	})

	return parser.Settings{
		Content:   f,
		DBhandler: handler,
	}
}
