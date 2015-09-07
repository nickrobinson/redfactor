package checker

import (
	"database/sql"
	"fmt"
	"github.com/influxdb/influxdb/client"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/url"
	"os"
	"time"
	"encoding/json"
	"github.com/nickrobinson/funcmon"
)

const (
	MyHost        = "192.168.1.206"
	MyPort        = 8086
	MyDB          = "redfactor"
	MyMeasurement = "cpu"
)

const (
	ABOVE_THRESHOLD = 0.0
	BELOW_THRESHOLD = 1.0
)

var fmClient *funcmon.Client

func RunChecker(c chan int) {
	db, err := sql.Open("sqlite3", "./redfactor.db")
	checkErr(err)

	u, err := url.Parse(fmt.Sprintf("http://%s:%d", MyHost, MyPort))
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL:      *u,
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PWD"),
	}

	con, err := client.NewClient(conf)
	checkErr(err)

	dur, ver, err := con.Ping()
	if err != nil {
		log.Fatal(err)

	}
	log.Printf("Happy as a Hippo! %v, %s", dur, ver)

	config := funcmon.Config{
		Host: "127.0.0.1",
		Port: 8086,
		DB: "funcmon",
	}
	fmClient,err = funcmon.NewClient(config)

	for {
		// query
		rows, err := db.Query("SELECT * FROM queries")
		checkErr(err)

		for rows.Next() {
			var rowId int
			var query string
			err = rows.Scan(&rowId, &query)
			checkErr(err)
			fmt.Println(rowId)
			fmt.Println(query)
			res, err := queryDB(con, query)
			checkErr(err)

			vals := res[0].Series[0].Values
			fmt.Printf("%v\n", vals[len(vals)-1])

			val, err := vals[len(vals)-1][1].(json.Number).Float64()
			checkErr(err)

			go evaluate(db, rowId, float64(val), c)
		}
		time.Sleep(5 * time.Minute)
	}
}

func evaluate(db *sql.DB, id int, value float64, c chan int) {
	fmClient.StartMonitoring("evaluate")
	query := fmt.Sprintf("SELECT * FROM thresholds WHERE id = %d", id)
	rows, err := db.Query(query)

	checkErr(err)

	for rows.Next() {
		var id int
		var threshold_type int
		var threshold float64
		var description string
		err = rows.Scan(&id, &threshold_type, &threshold, &description)
		checkErr(err)

		switch {
		case threshold_type == BELOW_THRESHOLD:
			if value < threshold {
				log.Printf("Detected Threshold Breach: %s", description)
				c <- id
			}
		case threshold_type == ABOVE_THRESHOLD:
			if value > threshold {
				log.Printf("Detected Threshold Breach: %s", description)
				c <- id
			}
		}
	}
	fmClient.EndMonitoring("evaluate")
}

// queryDB convenience function to query the database
func queryDB(con *client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: MyDB,
	}
	if response, err := con.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return res, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
