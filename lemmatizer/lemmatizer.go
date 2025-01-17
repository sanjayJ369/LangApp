package lemmatizer

import (
	"strings"
	"unicode"

	"github.com/aaaton/golem"
	"github.com/sanjayJ369/LangApp/dicts/en"
)

type Lemmatizer struct{}

func (l Lemmatizer) Lemmatize(word string) string {
	return Lemmatize(word)
}

func New() *Lemmatizer {
	return &Lemmatizer{}
}

func Lemmatize(word string) string {
	var fmtWord strings.Builder
	for _, r := range word {
		if unicode.IsLetter(r) {
			fmtWord.WriteRune(unicode.ToLower(r))
		}
	}
	word = fmtWord.String()
	lemmatizer, err := golem.New(en.New())
	if err != nil {
		panic(err)
	}
	return lemmatizer.Lemma(word)
}
