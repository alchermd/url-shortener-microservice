package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dns := os.Getenv("DATABASE_URL")

	if dns == "" {
		log.Fatal("$DATABASE_URL is not set.")
	}

	db, err := sql.Open("mysql", dns)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Connected to the database.")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS urls(id INT AUTO_INCREMENT PRIMARY KEY, original_url TEXT NOT NULL);")
	
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("Table created.")
	}
}
