package exporter

import (
	"strings"

	"github.com/sanjayJ369/LangApp/flashcard"
)

type Exporter struct{}

func New() Exporter {
	return Exporter{}
}

func (e Exporter) Export(cards []flashcard.Card) []byte {
	var builder strings.Builder

	// add headers
	builder.WriteString("#separator:tab\n")
	builder.WriteString("#html:false\n")

	// add content
	for i, card := range cards {
		builder.WriteString(card.Word + "\t" + card.Meaning)
		if i != len(cards)-1 {
			builder.WriteString("\n")
		}
	}

	return []byte(builder.String())
}
