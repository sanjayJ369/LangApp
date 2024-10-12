package parser_test

import (
	"testing"

	"github.com/sanjayJ369/LangApp/database"
	"github.com/sanjayJ369/LangApp/parser"
	"github.com/sanjayJ369/LangApp/testhelper"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Parallel()
	settings := validSettings(t)
	err := parser.New(settings).Parse()
	require.NoError(t, err, "parsing")
	meaning, err := settings.DBhandler.Get("abaiser")
	require.NoError(t, err)
	t.Log(meaning)
}

func TestParallelParse(t *testing.T) {
	t.Parallel()
	settings := validSettings(t)
	err := parser.New(settings).ParallelParse(1)
	require.NoError(t, err, "parsing")
	meaning, err := settings.DBhandler.Get("abaiser")
	require.NoError(t, err)
	t.Log(meaning)
}

func validSettings(tb testing.TB) parser.Settings {
	tb.Helper()
	name := testhelper.GetTempFileLoc()
	handler, err := database.NewBadger(name)
	require.NoError(tb, err, "creating handler")
	return parser.Settings{
		FileLoc:   "../assets/temp.json",
		DBhandler: handler,
	}
}
