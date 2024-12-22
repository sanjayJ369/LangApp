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
	DBhandler DBHandler
}

func check(s Settings) error {
	var aErr error

	if s.DBhandler == nil {
		aErr = errors.Join(aErr, errors.New("no db handler"))
	}

	return aErr
}

type Parser struct {
	dbhandler DBHandler
}

func New(settings Settings) (*Parser, error) {
	err := check(settings)
	if err != nil {
		return nil, fmt.Errorf("checking settings: %w", err)
	}

	return &Parser{
		dbhandler: settings.DBhandler,
	}, nil
}

type token struct {
	Word   string  `json:"word"`
	Senses []sense `json:"senses"`
}

type sense struct {
	Glosses []string `json:"glosses"`
}

func (p *Parser) Parse(content io.Reader) error {
	scn := bufio.NewScanner(content)
	maxTokenSize := bufio.MaxScanTokenSize
	buf := make([]byte, bufio.MaxScanTokenSize)
	scn.Buffer(buf, maxTokenSize*64)

	var (
		tok     token
		meaning strings.Builder
		err     error
	)

	for scn.Scan() {
		if len(scn.Bytes()) == 0 {
			return nil
		}

		err = json.Unmarshal(scn.Bytes(), &tok)
		if err != nil {
			return fmt.Errorf("unmarshalling token: %w", err)
		}

		err = insertToken(p.dbhandler, &tok, meaning)
		if err != nil {
			return fmt.Errorf("inserting token: %w", err)
		}
	}

	if err := scn.Err(); err != nil {
		return fmt.Errorf("scanning token: %w", err)
	}

	return nil
}

func insertToken(p DBHandler, tok *token, meaning strings.Builder) error {
	meaning.Reset()

	for i, sense := range tok.Senses {
		for j, gloss := range sense.Glosses {
			_, err := meaning.WriteString(gloss)
			if err != nil {
				return fmt.Errorf("writing gloss: %w", err)
			}

			if j < len(sense.Glosses)-1 {
				_, err = meaning.WriteString("\n")
				if err != nil {
					return fmt.Errorf("writing new line")
				}
			}
		}

		if i < len(tok.Senses)-1 {
			_, err := meaning.WriteString("\n")
			if err != nil {
				return fmt.Errorf("writing new line")
			}
		}
	}

	if meaning.Len() > 0 {
		err := p.Insert(tok.Word, meaning.String())
		if err != nil {
			return fmt.Errorf("inserting meaing: %w", err)
		}
	}

	return nil
}
