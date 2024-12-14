package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Translation struct {
	Senses string `json:"sense"`
	Word   string `json:"word"`
}

type Senses struct {
	Glosses      []string      `json:"glosses"`
	Transalation []Translation `json:"translations"`
}

type Word struct {
	Word   string   `json:"word"`
	Senses []Senses `json:"senses"`
}

type dbHandler interface {
	Insert(key, val string) error
	Get(key string) (string, error)
	Close() error
}

type Parser struct {
	fileloc   string
	dbhandler dbHandler
}

type Settings struct {
	FileLoc   string
	DBhandler dbHandler
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
	defer fp.Close()

	scn := bufio.NewScanner(fp)
	word := &Word{}

	for scn.Scan() {
		if len(scn.Bytes()) == 0 {
			return nil
		}
		err = json.Unmarshal(scn.Bytes(), word)
		if err != nil {
			return fmt.Errorf("unmarshalling json: %w", err)
		}
		err = insertWordContents(p.dbhandler, word)
		if err != nil {
			return fmt.Errorf("inserting word contents: %w", err)
		}
	}
	return nil
}

func insertWordContents(p dbHandler, word *Word) error {
	for _, sense := range word.Senses {
		if len(sense.Glosses) != 0 {
			var meaning strings.Builder
			for _, gloss := range sense.Glosses {
				_, err := meaning.WriteString(gloss + ",")
				if err != nil {
					return fmt.Errorf("concatinating glosses: %w", err)
				}
			}

			err := p.Insert(word.Word, meaning.String())
			if err != nil {
				return fmt.Errorf("inserting meaing values to db: %w", err)
			}
		}

		err := insertTranslations(p, sense.Transalation)
		if err != nil {
			return fmt.Errorf("inserting transations: %w", err)
		}
	}

	return nil
}

func insertTranslations(p dbHandler, translations []Translation) error {
	for _, lang := range translations {
		if len(lang.Senses) != 0 {
			err := p.Insert(lang.Word, lang.Senses)
			if err != nil {
				return fmt.Errorf("inserting meaing values to db: %w", err)
			}
		}
	}
	return nil
}
