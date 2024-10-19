package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sanjayJ369/LangApp/database"
	"github.com/sanjayJ369/LangApp/parser"
)

func main() {
	dictLoc := flag.String("dict", "", "dict is the location of dictionary which contains the word and its meaning")
	dbLoc := flag.String("db", "", "db is the location of the sqlite3 database where the data should be loaded into, new db is created is it is not present")
	tcount := flag.Int("t", 1, "t is the number of parallel threads to be run")
	parallel := flag.Bool("p", false, "p flag inicates if the parser should be run parallely ")

	flag.Parse()
	if *dictLoc == "" {
		log.Fatalln("dict flag is empty")
	}
	if *dbLoc == "" {
		log.Fatalln("db flag is empty")
	}

	handler, err := database.NewSqlite(*dbLoc)
	if err != nil {
		log.Fatalf("creating database: %s", err)
	}

	setting := parser.Settings{
		FileLoc:   *dictLoc,
		DBhandler: handler,
	}

	p := parser.New(setting)

	if *parallel {
		fmt.Println(*tcount)
		err = p.ParallelParse(*tcount)
	} else {
		err = p.Parse()
	}
	if err != nil {
		log.Fatalf("parsing file: %s", err)
	}

	log.Printf("Sucessfully Parsed File:%s, into Database:%s", *dictLoc, *dbLoc)
}
