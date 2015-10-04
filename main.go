package main

import (
	"flag"
	"./checker"
	"./notifier"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
)

var checkClient checker.Checker

func main() {
	//Start Martini
	m := martini.Classic()
	// render html templates from templates directory
	m.Use(render.Renderer(render.Options{
		Directory: "templates",
		Layout: "layout",
	}))

	m.Get("/", func(r render.Render) {
		r.HTML(http.StatusOK, "index", "info")
	})

	m.Get("/info", func(r render.Render) {
		r.HTML(http.StatusOK, "info", "info")
	})

	m.Get("/alarms", func(r render.Render) {
		var retData struct {
			ID string
		}
		retData.ID = "nick_test"

		r.HTML(http.StatusOK, "alarms", retData)
	})
	m.Run()

	c := make(chan int)

	//Intialize Flags
	var dbFile = flag.String("db", "./redfactor.db", "SQLite Database File")
	var hostname = flag.String("host", "192.168.1.206", "InfluxDB Hostname")
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
