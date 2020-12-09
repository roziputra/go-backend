package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "rozi"
	dbName := "scrap"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {

	db := dbConn()

	categoriesCollector := colly.NewCollector()
	booksCollector := colly.NewCollector()

	curtime := time.Now()
	datetime := curtime.Format("2006-01-02 15:04:05")

	var lastID int64

	// Find and collect Categories
	categoriesCollector.OnHTML("ul#zg_browseRoot ul ul", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, e *colly.HTMLElement) {

			var name, kode, categorylink string

			name = e.ChildText("a")
			categorylink = e.ChildAttr("a", "href")
			re := regexp.MustCompile("/zgbs/books/(.*)/ref=")
			match := re.FindStringSubmatch(categorylink)
			kode = match[1]

			if categorylink == "" || name == "" {
				// If we can't name or link , we return and go directly to the next element
				return
			}
			fmt.Printf("category -> (%s) %s \n", kode, name)

			// insert categories
			sql := "INSERT INTO categories(kode, name, date) VALUES (?, ?, ?)"
			insertCategories, err := db.Prepare(sql)

			if err != nil {
				panic(err.Error())
			}

			res, err := insertCategories.Exec(kode, name, datetime)
			if err != nil {
				panic(err.Error())
			}

			lastinsert, err := res.LastInsertId()

			if err != nil {
				log.Fatal(err)
			}

			lastID = lastinsert

			//visit category link
			booksCollector.Visit(categorylink)
		})
	})

	// Find and collect Books
	booksCollector.OnHTML("ol#zg-ordered-list", func(e *colly.HTMLElement) {
		e.ForEach("li.zg-item-immersion", func(_ int, e *colly.HTMLElement) {
			var booktitle, bookkode, bookformat, bookprice, booklink string

			booklink = e.ChildAttr("span.zg-item a", "href")
			re := regexp.MustCompile("/dp/(.*)/ref=")
			match := re.FindStringSubmatch(booklink)
			bookkode = match[1]
			booktitle = e.ChildAttr(".zg-item a div.p13n-sc-truncate-desktop-type2", "title")
			if booktitle == "" {
				booktitle = e.ChildText(".zg-item a div.p13n-sc-truncate-desktop-type2")
			}
			bookformat = e.ChildText(".zg-item .a-row .a-size-small.a-color-secondary")
			bookprice = strings.TrimLeft(e.ChildText(".zg-item a .p13n-sc-price"), "$")

			fmt.Printf("book -> (%s) %s \n %s price: %s \n", bookkode, booktitle, bookformat, bookprice)

			// insert categories
			sqlbooks := "INSERT INTO books(kode, title, category, format, price, date) VALUES (?, ?, ?, ?, ?, ?)"
			insertBooks, err := db.Prepare(sqlbooks)

			if err != nil {
				panic(err.Error())
			}
			insertBooks.Exec(bookkode, booktitle, lastID, bookformat, bookprice, datetime)

		})
	})

	//visit next page
	booksCollector.OnHTML("ul.a-pagination li.a-last", func(e *colly.HTMLElement) {
		e.Request.Visit(e.ChildAttr("a", "href"))
	})
	categoriesCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	booksCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visit Category", r.URL)
	})

	categoriesCollector.Visit("https://www.amazon.com/best-sellers-books-Amazon/zgbs/books/ref=zg_bs_nav_0")
}
