package parser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

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

type wordMeaning struct {
	word    string
	meaning string
}

type noMeaning struct {
	error
}

func (n *noMeaning) Error() string {
	return "word has no meaning"
}

var noMeaningError *noMeaning = &noMeaning{}

type ParallelParser struct {
	sync.Mutex
	insert    chan wordMeaning
	reader    *bufio.Reader
	BytesRead int64
	limit     int64
}

func (p *ParallelParser) Read() error {

	for {
		line, err := p.reader.ReadBytes('\n')
		p.BytesRead += int64(len(line))
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("reading bytes: %w", err)
		}

		// check the first character
		// if it is not '{' skip
		// as the thread might have started in middle of a json entry
		if len(line) < 7 || string(line[:7]) != `{"pos":` {
			continue
		}

		wm, err := extractWordMeaning(line)
		if errors.As(err, &noMeaningError) {
			continue
		}

		if err != nil {
			return fmt.Errorf("extracting word and meaning: %w", err)
		}

		p.insert <- *wm
		if p.BytesRead >= p.limit {
			return nil
		}
	}
}

func extractWordMeaning(line []byte) (*wordMeaning, error) {
	dec := json.NewDecoder(bytes.NewReader(line))
	wm := &wordMeaning{}
	for {
		t, err := dec.Token()
		if err == io.EOF {
			return wm, nil
		}
		if err != nil {
			return wm, fmt.Errorf("getting token: %w", err)
		}
		if t == "word" {
			wordtkn, err := dec.Token()
			if err != nil {
				return wm, fmt.Errorf("getting word token: %w", err)
			}
			wm.word = fmt.Sprintf("%s", wordtkn)
			wm.meaning, err = extractparser(dec)

			if err != nil {
				if errors.As(err, &noMeaningError) {
					return wm, noMeaningError
				}
				return wm, fmt.Errorf("extracting meaning: %w", err)
			}
			return wm, nil
		}
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
			if err != nil {
				return fmt.Errorf("extracting phrase: %w", err)
			}

			err = p.dbhandler.Insert(word, meaning)
			if err != nil {
				return fmt.Errorf("inserting meaing values to db: %w", err)
			}
		}
	}
}

func (p *Parser) ParallelParse(threadCount int) error {

	insertChan := make(chan wordMeaning)
	errChan := make(chan error)

	var insertedVals atomic.Int64
	insertedVals.Store(0)

	activeThreads := threadCount

	var parserwg sync.WaitGroup
	var insertwg sync.WaitGroup

	readers := []*ParallelParser{}
	fi, err := os.Stat(p.fileloc)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	fileSize := fi.Size()
	acc := int64(0)
	limit := fileSize / int64(threadCount)

	for i := 0; i < int(threadCount); i++ {

		fptr, err := os.Open(p.fileloc)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}

		fptr.Seek(acc, io.SeekStart)
		r := &ParallelParser{
			insert:    insertChan,
			reader:    bufio.NewReader(fptr),
			BytesRead: 0,
			limit:     int64(float64(limit) * 1.4),
		}

		acc += limit
		readers = append(readers, r)
		parserwg.Add(1)
	}

	for _, r := range readers {
		go func(r *ParallelParser) {
			err := r.Read()
			if err != nil {
				errChan <- err
			}
			parserwg.Done()
			activeThreads -= 1
		}(r)
	}

	go func() {
		for {
			time.Sleep(5 * time.Second)
			var total int64
			for _, r := range readers {
				r.Lock()
				total += r.BytesRead
				r.Unlock()
			}
			fmt.Println("\nread: ", total, "bytes")
			fmt.Println("active threads: ", activeThreads)
			fmt.Printf("completed: %f%%\n", float64(total)/float64(fileSize)*100)
			fmt.Println("insert buffer:", len(insertChan))
			fmt.Println("\ninserted vals:", insertedVals.Load())
		}
	}()

	go func() {
		for i := 0; i < threadCount*2; i++ {
			insertwg.Add(1)
			go func() {
				for {
					wm, ok := <-insertChan
					if !ok {
						insertwg.Done()
						return
					}
					insertedVals.Store(insertedVals.Load() + 1)
					err := p.dbhandler.Insert(wm.word, wm.meaning)
					if err != nil {
						errChan <- err
					}
				}
			}()
		}

	}()

	go func() {
		parserwg.Wait()
		close(insertChan)
		insertwg.Wait()

		// complete remaining insertions in the buffer
		for wm := range insertChan {
			err = p.dbhandler.Insert(wm.word, wm.meaning)
			if err != nil {
				errChan <- err
			}
		}
		errChan <- nil
	}()

	return <-errChan
}

func extractparser(dec *json.Decoder) (string, error) {

	for {
		w, err := dec.Token()
		if err != nil {
			return "", fmt.Errorf("getting meaning token: %w", err)
		}
		if w == "no-gloss" {
			return "", noMeaningError
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
