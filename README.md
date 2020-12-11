# go-backend

simple golang backend

install golang

install mysql and create database scrap with scrap.sql

# Run Scrapper

install mysql driver for go \n
run "go get -u github.com/go-sql-driver/mysql"

install gocolly, a web scrapper library for go \n
run "go get -u github.com/gocolly/colly/..."

run "go run scrapper/scrapper.go" to get data from the web

# Run Server

install gorilla/mux a request router for go \n
run "go get -u github.com/gorilla/mux"

run "go run main.go" to run server \n

server run at http://localhost:8080/api/v1

get all categories : http://localhost:8080/api/v1/books/categories

get books by category : http://localhost:8080/api/v1/categories/{categoryKode}

get books by kode : http://localhost:8080/api/v1/book/detail/{bookKode}

get books by filter : http://localhost:8080/api/v1/books/filter

