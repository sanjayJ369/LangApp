package parser

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	content   io.Reader
	dbhandler dbHandler
}

type Settings struct {
	Content   io.Reader
	DBhandler dbHandler
}

func check(s Settings) error {
	var aErr error

	if s.Content == nil {
		aErr = errors.Join(aErr, errors.New("no content"))
	}

	if s.DBhandler == nil {
		aErr = errors.Join(aErr, errors.New("no db handler"))
	}

	return aErr
}

func New(settings Settings) (*Parser, error) {
	err := check(settings)
	if err != nil {
		return nil, fmt.Errorf("checking settings: %w", err)
	}

	return &Parser{
		content:   settings.Content,
		dbhandler: settings.DBhandler,
	}, nil
}

func (p *Parser) Parse() error {
	fp := p.content

	scn := bufio.NewScanner(fp)
	word := &Word{}
	var err error

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
