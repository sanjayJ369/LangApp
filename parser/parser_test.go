package parser_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/sanjayJ369/LangApp/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Parallel()

	f, err := os.Open("./testfiles/vocabulary.jsonl")
	require.NoError(t, err, "opening vocabulary")
	t.Cleanup(func() {
		require.NoError(t, f.Close(), "closing vocabulary")
	})

	var h fakeDBHandler

	p, err := parser.New(parser.Settings{
		DBhandler: &h,
	})
	require.NoError(t, err, "creating parser")

	require.NoError(t, p.Parse(f), "parsing")

	assert.Equal(t, "Ivory black; animal charcoal.", h.get("abaiser"), "wrong result")
	assert.Equal(
		t,
		strings.Join([]string{
			"The act of destroying a fetus in the womb; feticide.",
			"An agent responsible for an abortion (the destruction of a fetus); abortifacient.",
		}, "\n"),
		h.get("aborticide"),
		"wrong result",
	)
}

func BenchmarkParser(b *testing.B) {
	f, err := os.Open("./testfiles/vocabulary.jsonl")
	require.NoError(b, err, "opening vocabulary")
	b.Cleanup(func() {
		require.NoError(b, f.Close(), "closing vocabulary")
	})

	var h noOpDBHandler

	p, err := parser.New(parser.Settings{
		DBhandler: h,
	})
	require.NoError(b, err, "creating parser")

	content, err := io.ReadAll(f)
	require.NoError(b, err, "reading vocabulary")

	r := bytes.NewReader(content)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		r.Reset(content)

		require.NoError(b, p.Parse(r), "parsing")
	}
}

type fakeDBHandler struct {
	meanings [][]string
}

func (f *fakeDBHandler) Insert(key, val string) error {
	f.meanings = append(f.meanings, []string{key, val})

	return nil
}

func (f *fakeDBHandler) get(key string) string {
	for _, v := range f.meanings {
		if v[0] == key {
			return v[1]
		}
	}

	return ""
}

type noOpDBHandler struct{}

func (n noOpDBHandler) Insert(key, val string) error {
	return nil
}
