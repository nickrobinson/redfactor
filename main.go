package main

import (
	"fmt"
	"flag"
	"./checker"
	"./notifier"
)

var checkClient checker.Checker

func main() {
	c := make(chan int)
	fmt.Println("Hello World")

	var dbFile = flag.String("db", "./redfactor.db", "SQLite Database File")
	var hostname = flag.String("host", "127.0.0.1", "InfluxDB Hostname")
	var port = flag.Int("port", 8086, "InfluxDB Port")
	var database = flag.String("database", "redfactor", "InfluxDB Database Name")
	flag.Parse()

	checkClient.NewChecker(*dbFile, *hostname, *port, *database, "cpu")

	go checkClient.StartChecker(c)
	go notifier.StartNotifier(c)
	select {}
}
