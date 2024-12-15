package parser_test

import (
	"os"
	"testing"

	"github.com/sanjayJ369/LangApp/parser"
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

	f, err := os.Open("./testfiles/vocabulary.jsonl")
	require.NoError(tb, err, "opening vocabulary file")

	tb.Cleanup(func() {
		require.NoError(tb, f.Close(), "closing vocabulary file")
	})

	var h fakeDBHandler
	tb.Cleanup(func() {
		require.NoError(tb, h.Close(), "closing db handler")
	})

	return parser.Settings{
		Content:   f,
		DBhandler: &h,
	}
}

type fakeDBHandler struct {
	meanings [][]string
}

func (f *fakeDBHandler) Insert(key, val string) error {
	f.meanings = append(f.meanings, []string{key, val})

	return nil
}

func (f *fakeDBHandler) Get(key string) (string, error) {
	for _, v := range f.meanings {
		if v[0] == key {
			return v[1], nil
		}
	}

	return "", nil
}

func (f *fakeDBHandler) Close() error {
	return nil
}
