package lemmatizer

import (
	"github.com/aaaton/golem"
	"github.com/sanjayJ369/LangApp/dicts/en"
)

func Lemmatize(word string) string {
	lemmatizer, err := golem.New(en.New())
	if err != nil {
		panic(err)
	}
	return lemmatizer.Lemma(word)
}
