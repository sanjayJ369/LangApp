package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sanjayJ369/LangApp/database"
)

type Transation struct {
	Senses string `json:"sense"`
	Word   string `json:"word"`
}

type Senses struct {
	Glosses      []string     `json:"glosses"`
	Transalation []Transation `json:"translations"`
}

type Word struct {
	Word   string   `json:"word"`
	Senses []Senses `json:"senses"`
}

type Insertable interface {
	Insert(string, string) error
}

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

func (p *ParallelParser) Insert(word, meaning string) error {
	p.insert <- wordMeaning{
		word:    word,
		meaning: meaning,
	}
	return nil
}

func (p *ParallelParser) Read() error {

	maxTokenSize := bufio.MaxScanTokenSize
	scn := bufio.NewScanner(p.reader)
	buf := make([]byte, bufio.MaxScanTokenSize)
	scn.Buffer(buf, maxTokenSize*64)
	word := &Word{}

	for scn.Scan() {
		line := scn.Text()
		p.BytesRead += int64(len(line))

		if len(line) < 7 || string(line[:7]) != `{"pos":` {
			continue
		}

		err := json.Unmarshal([]byte(line), word)
		if err != nil {
			log.Fatalln("unmarshing json: %w", err)
		}

		err = InsertWordContents(p, word)
		if err != nil {
			return fmt.Errorf("inserting word contents: %w", err)
		}

		if p.BytesRead >= p.limit {
			return nil
		}
	}

	if scn.Err() != nil {
		return fmt.Errorf("scanning file: %w", scn.Err())
	}

	return nil
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
		err = InsertWordContents(p.dbhandler, word)
		if err != nil {
			return fmt.Errorf("inserting word contents: %w", err)
		}
	}
	return nil
}

func InsertWordContents(p Insertable, word *Word) error {
	// insert word and glosses present in the senses
	for _, sense := range word.Senses {
		// insert sense gloss
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

		// insert sense translations
		err := InsertTranslations(p, sense.Transalation)
		if err != nil {
			return fmt.Errorf("inserting transations: %w", err)
		}
	}

	return nil
}

func InsertTranslations(p Insertable, translations []Transation) error {
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

func (p *Parser) Close() error {
	return p.dbhandler.Close()
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
			limit:     int64(float64(limit) * 1.2),
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
			preInsertedVals := insertedVals.Load()
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
			fmt.Println("inserted vals:", insertedVals.Load())
			fmt.Printf("rate of insertion: %d vals/sec\n", (insertedVals.Load()-preInsertedVals)/5)
		}
	}()

	go func() {
		for i := 0; i < threadCount; i++ {
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

func GetLineCount(path string) int {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("getting line count: %s", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	lineCount := 0
	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		lineCount++
	}

	return lineCount
}
