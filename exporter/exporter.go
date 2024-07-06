package exporter

import "github.com/sanjayJ369/LangApp/flashcard"

type Exporter struct{}

func New() Exporter {
	return Exporter{}
}
func (e Exporter) Export(card []flashcard.Card) []byte {
	return []byte("hello")
}
