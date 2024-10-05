package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sanjayJ369/LangApp/database"
)

type Parser struct {
	fileloc   string
	dbhandler database.Handler
}

type Settings struct {
	FileLoc   string
	DBhandler database.Handler
}

func New(s Settings) *Parser {
	return &Parser{
		fileloc:   s.FileLoc,
		dbhandler: s.DBhandler,
	}
}

func (p *Parser) Parse() error {

	fp, err := os.OpenFile(p.fileloc, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}

	dec := json.NewDecoder(fp)
	var word string
	var meaning string
	for {
		t, err := dec.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("getting token: %w", err)
		}
		if t == "word" {
			wordtkn, err := dec.Token()
			if err != nil {
				return fmt.Errorf("getting word token: %w", err)
			}
			word = fmt.Sprintf("%s", wordtkn)
			meaning, err = extractparser(dec)

			err = p.dbhandler.Insert(word, meaning)
			if err != nil {
				return fmt.Errorf("inserting meaing values to db: %w", err)
			}
		}
	}
}

func extractparser(dec *json.Decoder) (string, error) {

	for {
		w, err := dec.Token()
		if err != nil {
			return "", fmt.Errorf("getting meaning token: %w", err)
		}
		if w == "glosses" {
			w, err = dec.Token()
			if err != nil {
				return "", fmt.Errorf("getting glosses token: %w", err)
			}
			if fmt.Sprint(w) == "[" {
				return extractGlosses(dec)
			}
		}
	}
}

func extractGlosses(dec *json.Decoder) (string, error) {
	var parser strings.Builder
	for {
		w, err := dec.Token()
		if err != nil {
			return "", fmt.Errorf("getting glosses value token: %w", err)
		}
		if fmt.Sprint(w) == "]" {
			break
		}
		parser.WriteString(fmt.Sprint(w))
	}

	return parser.String(), nil
}
