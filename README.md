# go-backend

simple golang backend

install golang

install mysql and create database scrap with scrap.sql

# Run Scrapper

run "go run scrapper/scrapper.go" to get data from the web

# Run Server

run "go run main.go" to run server

server run at http://localhost:8080/api/v1

get all categories : http://localhost:8080/api/v1/books/categories

get books by category : http://localhost:8080/api/v1/categories/{categoryKode}

get books by kode : http://localhost:8080/api/v1/book/detail/{bookKode}

get books by filter : http://localhost:8080/api/v1/books/filter

