package parser

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type DBHandler interface {
	Insert(key, val string) error
}

type Settings struct {
	Content   io.Reader
	DBhandler DBHandler
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

type Parser struct {
	content   io.Reader
	dbhandler DBHandler
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

type token struct {
	Word   string  `json:"word"`
	Senses []sense `json:"senses"`
}

type sense struct {
	Glosses       []string      `json:"glosses"`
	Transalations []translation `json:"translations"`
}

type translation struct {
	Word  string `json:"word"`
	Sense string `json:"sense"`
}

func (p *Parser) Parse() error {
	scn := bufio.NewScanner(p.content)

	var (
		tok token
		err error
	)

	for scn.Scan() {
		if len(scn.Bytes()) == 0 {
			return nil
		}

		err = json.Unmarshal(scn.Bytes(), &tok)
		if err != nil {
			return fmt.Errorf("unmarshalling token: %w", err)
		}

		err = insertToken(p.dbhandler, &tok)
		if err != nil {
			return fmt.Errorf("inserting token: %w", err)
		}
	}

	return nil
}

func insertToken(p DBHandler, tok *token) error {
	for _, sense := range tok.Senses {
		if len(sense.Glosses) != 0 {
			var meaning strings.Builder

			for _, gloss := range sense.Glosses {
				_, err := meaning.WriteString(gloss + ",")
				if err != nil {
					return fmt.Errorf("concatinating glosses: %w", err)
				}
			}

			err := p.Insert(tok.Word, meaning.String())
			if err != nil {
				return fmt.Errorf("inserting meaing values to db: %w", err)
			}
		}

		err := insertTranslations(p, sense.Transalations)
		if err != nil {
			return fmt.Errorf("inserting transations: %w", err)
		}
	}

	return nil
}

func insertTranslations(p DBHandler, translations []translation) error {
	for _, lang := range translations {
		if len(lang.Sense) != 0 {
			err := p.Insert(lang.Word, lang.Sense)
			if err != nil {
				return fmt.Errorf("inserting meaing values to db: %w", err)
			}
		}
	}

	return nil
}
