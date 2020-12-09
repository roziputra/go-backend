package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

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
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("ul#zg_browseRoot ul ul", func(e *colly.HTMLElement) {
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

			//visit each category link
			c2 := colly.NewCollector()

			c2.OnHTML("ol#zg-ordered-list", func(e *colly.HTMLElement) {
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
				})

			})

			c2.OnRequest(func(r *colly.Request) {
				fmt.Println("Visit Category", r.URL)
			})

			c2.Visit(categorylink)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.amazon.com/best-sellers-books-Amazon/zgbs/books/ref=zg_bs_nav_0")
}
