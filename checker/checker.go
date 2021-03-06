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
)

const (
	MyHost        = "192.168.1.206"
	MyPort        = 8086
	MyDB          = "redfactor"
)

const (
	ABOVE_THRESHOLD = 0.0
	BELOW_THRESHOLD = 1.0
)

type Checker struct {
	db *sql.DB
	host string
	port int
	influxDb string
}

func (checkerPtr *Checker) NewChecker(filename_ *sql.DB, host_ string, port_ int, influxDb_ string) {
	checkerPtr.db = filename_
	checkerPtr.host = host_
	checkerPtr.port = port_
	checkerPtr.influxDb = influxDb_
	return
}

func (checkerPtr *Checker) StartChecker(c chan int) {
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", checkerPtr.host, checkerPtr.port))
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

	for {
		// query
		rows, err := checkerPtr.db.Query("SELECT * FROM queries")
		if (err != nil) {
			log.Fatal(err)
		} else {
			for rows.Next() {
				var rowId int
				var query string
				var thresholdId int
				err = rows.Scan(&rowId, &query, &thresholdId)
				checkErr(err)
				fmt.Println(rowId)
				fmt.Println(query)
				res, err := queryDB(con, query)
				checkErr(err)

				vals := res[0].Series[0].Values

				//Continue if value is nil
				if (vals[len(vals) - 1][1] == nil) {
					log.Printf("%s: Got nil value from db query", "StartChecker")
					continue
				}

				val, err := vals[len(vals) - 1][1].(json.Number).Float64()
				checkErr(err)

				go evaluate(checkerPtr.db, thresholdId, float64(val), c)
			}
		}

		time.Sleep(5 * time.Minute)
	}
}

func evaluate(db *sql.DB, id int, value float64, c chan int) {
	query := fmt.Sprintf("SELECT * FROM thresholds WHERE id = %d", id)
	rows, err := db.Query(query)

	if (err != nil) {
		log.Fatal(err)
		return
	}

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
