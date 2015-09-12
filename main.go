package main

import (
	"flag"
	"./checker"
	"./notifier"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var checkClient checker.Checker

func main() {
	c := make(chan int)

	//Intialize Flags
	var dbFile = flag.String("db", "./redfactor.db", "SQLite Database File")
	var hostname = flag.String("host", "127.0.0.1", "InfluxDB Hostname")
	var port = flag.Int("port", 8086, "InfluxDB Port")
	var database = flag.String("database", "redfactor", "InfluxDB Database Name")
	flag.Parse()

	db, err := sql.Open("sqlite3", *dbFile)
	if (err != nil) {
		log.Fatal(err)
	}

	checkClient.NewChecker(db, *hostname, *port, *database)

	go checkClient.StartChecker(c)
	go notifier.StartNotifier(c)
	select {}
}
