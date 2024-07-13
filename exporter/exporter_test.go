package exporter_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sanjayJ369/LangApp/exporter"
	"github.com/sanjayJ369/LangApp/flashcard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExporter(t *testing.T) {
	t.Parallel()

	/* Exporter exports Learner cards. */

	t.Run("Exporter exports Anki cards", func(t *testing.T) {
		t.Parallel()

		// Given Learner has some flashcards.
		cards := []flashcard.Card{
			{
				Word:    "word",
				Meaning: "meaning",
			},
			{
				Word:    "flash",
				Meaning: "card",
			},
		}

		// When Learner Exports to Anki.
		e := exporter.New()
		got := e.Export(cards)
		gotLines := strings.Split(string(got), "\n")

		// Then Cards are Exported.
		require.Greater(t, len(gotLines), 2)
		require.Equal(t, len(gotLines), len(cards)+2)
		assert.Contains(t, gotLines[0], "#separator:tab")
		assert.Contains(t, gotLines[1], "#html:false")

		for i, line := range gotLines[2:] {
			assert.Contains(t, line, cards[i].Word+"\t"+cards[i].Meaning)
		}

		fmt.Printf("%s", string(got))
	})
}
