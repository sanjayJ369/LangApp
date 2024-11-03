package parser_test

import (
	"fmt"
	"testing"

	"github.com/sanjayJ369/LangApp/database"
	"github.com/sanjayJ369/LangApp/parser"
	"github.com/sanjayJ369/LangApp/testhelper"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Parallel()
	settings, clean := validSettings(t)
	t.Cleanup(clean)
	err := parser.New(settings).Parse()
	require.NoError(t, err, "parsing")
	meaning, err := settings.DBhandler.Get("abaiser")
	require.NoError(t, err)
	t.Log(meaning)
}

func TestParallelParse(t *testing.T) {
	t.Parallel()
	settings, clean := validSettings(t)
	t.Cleanup(clean)
	err := parser.New(settings).ParallelParse(1)
	require.NoError(t, err, "parsing")
	meaning, err := settings.DBhandler.Get("abaiser")
	require.NoError(t, err)
	t.Log(meaning)
}

func validSettings(tb testing.TB) (parser.Settings, func()) {
	tb.Helper()
	name := testhelper.GetTempFileLoc()
	handler, err := database.NewBadger(name)
	require.NoError(tb, err, "creating handler")
	return parser.Settings{
			FileLoc:   "../assets/temp.json",
			DBhandler: handler,
		}, func() {
			fmt.Println("closing db(parser)")
			require.NoError(tb, handler.Close(), "closing badger")
		}
}
