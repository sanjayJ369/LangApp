package meaning

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Meaning struct {
	fileloc string
}

type Settings struct {
	FileLoc string
}

func New(s Settings) Meaning {
	return Meaning{
		fileloc: s.FileLoc,
	}
}

func (m Meaning) GetMeaning(word string) string {

	word = strings.ToLower(word)
	fp, err := os.OpenFile(m.fileloc, os.O_RDONLY, 0644)
	if err != nil {
		return ""
	}

	dec := json.NewDecoder(fp)

	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if t == "word" {
			w, _ := dec.Token()
			if fmt.Sprint(w) == word {
				return extractMeaning(dec)
			}
		}
	}
	return ""
}

func extractMeaning(dec *json.Decoder) string {

	for {
		w, _ := dec.Token()
		if w == "glosses" {
			w, _ = dec.Token()
			if fmt.Sprint(w) == "[" {
				return extractGlosses(dec)
			}
		}
	}
}

func extractGlosses(dec *json.Decoder) string {
	var meaning strings.Builder
	for {
		w, _ := dec.Token()
		if fmt.Sprint(w) == "]" {
			break
		}
		meaning.WriteString(fmt.Sprint(w))
	}

	return meaning.String()
}
