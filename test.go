package main

import (
	"log"
	"strings"

	"database/sql"

	rent "github.com/FrankNine/fiveN1-rent-scraper"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	options := rent.NewOptions()
	url, err := rent.GenerateURL(options)
	if err != nil {
		log.Fatalf("\x1b[91;1m%s\x1b[0m", err)
	}

	f := rent.NewFiveN1(url)
	if err := f.Scrape(1); err != nil {
		log.Fatal(err)
	}

	json := rent.ConvertToJSON(f.RentList)
	log.Println(string(json))

	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)

	stmt, err := db.Prepare("INSERT INTO rentinfo(title, url, address, floor, max_floor) values(?,?,?,?,?)")
	checkErr(err)

	for _, page := range f.RentList {
		for _, r := range page {
			floorWithoutPrefix := strings.TrimPrefix(r.Floor, "樓層：")
			hasMaxFloor := strings.Contains(floorWithoutPrefix, "/")
			var floor string
			var maxFloor string
			if hasMaxFloor {
				floors := strings.Split(floorWithoutPrefix, "/")
				floor = floors[0]
				maxFloor = floors[1]
			} else {
				floor = floorWithoutPrefix
				maxFloor = ""
			}

			_, err := stmt.Exec(r.Title, r.URL, r.Address, floor, maxFloor)
			checkErr(err)
		}
	}

	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
