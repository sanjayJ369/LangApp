package main

import (
	"flag"
	"log"
	"os"

	"github.com/sanjayJ369/LangApp/database"
	"github.com/sanjayJ369/LangApp/parser"
)

func main() {
	dictLoc := flag.String("dict", "", "dict is the location of dictionary which contains the word and its meaning")
	dbLoc := flag.String("db", "", "db is the location of the sqlite3 database where the data should be loaded into, new db is created is it is not present")

	flag.Parse()
	if *dictLoc == "" {
		log.Fatalln("dict flag is empty")
	}
	if *dbLoc == "" {
		log.Fatalln("db flag is empty")
	}

	handler, err := database.NewBadger(*dbLoc)
	defer func() {
		err := handler.Close()
		if err != nil {
			log.Fatalf("closing database: %s", err.Error())
		}
	}()

	if err != nil {
		log.Fatalf("creating database: %s", err)
	}

	dict, err := os.Open(*dictLoc)
	if err != nil {
		log.Fatalf("opening dictionary: %s", err)
	}
	defer dict.Close()

	setting := parser.Settings{
		Content:   dict,
		DBhandler: handler,
	}

	p, err := parser.New(setting)
	if err != nil {
		log.Fatalf("creating parser: %s", err)
	}

	err = p.Parse()
	if err != nil {
		log.Fatalf("parsing file: %s", err)
	}

	log.Printf("Sucessfully Parsed File:%s, into Database:%s", *dictLoc, *dbLoc)
}
