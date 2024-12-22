package main

import (
	"flag"
	"log"
	"os"

	"github.com/sanjayJ369/LangApp/database"
	"github.com/sanjayJ369/LangApp/exporter"
	"github.com/sanjayJ369/LangApp/flashcard"
	"github.com/sanjayJ369/LangApp/learner"
	"github.com/sanjayJ369/LangApp/lemmatizer"
	"github.com/sanjayJ369/LangApp/meaning"
)

func main() {

	text := flag.String("t", "", "text from which we want to extract the meaning from")
	fileLoc := flag.String("l", "./cards.txt", "location to store the exported cards")
	flag.Parse()
	if *text == "" {
		log.Fatalln("input text is empty")
	}

	learnerId := "test"

	dbhandler, err := database.NewSqlite("./assets/meaning.db")
	if err != nil {
		log.Fatalf("opening db: %s", err)
	}
	defer dbhandler.Close()
	
	meaningH, err := meaning.New(meaning.Settings{
		GetMeaning: dbhandler,
	})
	if err != nil {
		log.Fatalf("creating mining handler: %s", err)
	}

	settings := flashcard.Settings{
		Learner:    learner.New("../../assets/learnerData"),
		Meaning:    meaningH,
		Exporter:   exporter.New(),
		Lemmatizer: lemmatizer.New(),
	}
	fc, err := flashcard.New(settings)
	if err != nil {
		log.Fatalf("creating flashcards: %s", err)
	}

	fc.CreateFlashCards(learnerId, *text)
	data := fc.Export(learnerId)
	fp, err := os.Create(*fileLoc)
	if err != nil {
		log.Fatalf("creating file: %s", err)
	}

	_, err = fp.Write(data)
	if err != nil {
		log.Fatalf("writing data to file: %s", err)
	}

	err = fp.Close()
	if err != nil {
		log.Fatalf("closing file: %s", err)
	}
}
