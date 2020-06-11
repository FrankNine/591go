package main

import (
	"fmt"
	"log"
	"strings"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	rent "github.com/neighborhood999/fiveN1-rent-scraper"
)

func main() {
	options := rent.NewOptions()
	// https://github.com/neighborhood999/fiveN1-rent-scraper/blob/master/list/url-jump-ip.md
	options.Region = 1
	// https://github.com/neighborhood999/fiveN1-rent-scraper/blob/master/list/section.md
	options.Section = "4,7"
	options.Kind = 0
	options.RentPrice = "10000, 40000"

	url, err := rent.GenerateURL(options)
	if err != nil {
		log.Fatalf("\x1b[91;1m%s\x1b[0m", err)
	}

	f := rent.NewFiveN1(url)
	if err := f.Scrape(f.GetTotalPage()); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)

	db.Exec("TRUNCATE rentinfo")

	stmt, err := db.Prepare("INSERT INTO rentinfo(title, url, address, floor, max_floor, is_new, ping, price, rent_type, option_type) values(?,?,?,?,?,?,?,?,?,?)")
	checkErr(err)

	for p, page := range f.RentList {
		fmt.Printf("Page: %d", p)
		for _, r := range page {
			fmt.Println(r.Title)

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
			_, err := stmt.Exec(r.Title, r.URL, r.Address, floor, maxFloor, r.IsNew, r.Ping, r.Price, r.RentType, r.OptionType)
			checkErr(err)
		}
	}

	db.Close()
}

func checkErr(err error) {
	if err != nil && !strings.Contains(err.Error(), "UNIQUE") {
		panic(err)
	}
}
