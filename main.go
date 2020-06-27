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
	// https://github.com/neighborhood999/fiveN1-rent-scraper/blob/master/list/url-jump-ip.md
	// https://github.com/neighborhood999/fiveN1-rent-scraper/blob/master/list/section.md
	songXinInfo := dumpRegion(1, "4,7", "10000, 40000")
	writeDatabase(songXinInfo, "松山信義.db")

	daInfo := dumpRegion(1, "5", "10000, 40000")
	writeDatabase(daInfo, "大安.db")

	zongInfo := dumpRegion(1, "1,3", "10000, 40000")
	writeDatabase(zongInfo, "中山中正.db")
}

func dumpRegion(region int, section string, rentPrice string) rent.HouseInfoCollection {
	options := rent.NewOptions()
	options.Region = region
	options.Section = section
	options.Kind = 0
	options.RentPrice = rentPrice

	url, err := rent.GenerateURL(options)
	if err != nil {
		log.Fatalf("\x1b[91;1m%s\x1b[0m", err)
	}

	f := rent.NewFiveN1(url)
	if err := f.Scrape(f.GetTotalPage()); err != nil {
		log.Fatal(err)
	}

	return f.RentList
}

func writeDatabase(infoCollection rent.HouseInfoCollection, databaseName string) {
	db, err := sql.Open("sqlite3", databaseName)
	checkErr(err)

	db.Exec(`CREATE TABLE IF NOT EXISTS"rentInfo" (
		"title"	TEXT,
		"url"	TEXT UNIQUE,
		"address"	TEXT,
		"floor"	TEXT,
		"max_floor"	TEXT,
		"is_new"	TEXT,
		"ping"	TEXT,
		"price"	TEXT,
		"rent_type"	TEXT,
		"option_type"	TEXT
	)`)

	db.Exec("TRUNCATE rentinfo")

	stmt, err := db.Prepare("INSERT INTO rentinfo(title, url, address, floor, max_floor, is_new, ping, price, rent_type, option_type) values(?,?,?,?,?,?,?,?,?,?)")
	checkErr(err)

	for p, page := range infoCollection {
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
