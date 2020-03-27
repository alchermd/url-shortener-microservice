package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type ShortUrl struct {
	OriginalUrl string `json:"original_url"`
	ShortUrl    int64  `json:"short_url"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("views/index.html"))
	t.Execute(w, nil)
}

func shortenerHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Allow-Access-Control-Origin", "*")

	payload := getPayload(r)

	result, err := db.Exec("INSERT INTO urls(original_url) VALUES(?)", payload["url"])

	if err != nil {
		http.Error(w, `{"message": "Something went wrong"}`, http.StatusInternalServerError)
		return
	}

	shortUrl, err := result.LastInsertId()

	if err != nil {
		http.Error(w, `{"message": "Something went wrong"}`, http.StatusInternalServerError)
		return
	}

	url := &ShortUrl{
		OriginalUrl: payload["url"],
		ShortUrl:    shortUrl,
	}

	j, err := json.Marshal(url)

	if err != nil {
		http.Error(w, `{"message": "Something went wrong"}`, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(j))
}

func redirectHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Allow-Access-Control-Origin", "*")

	id := r.URL.Path[len("/api/shorturl/"):]

	fmt.Fprintf(w, `{"message": "hello, %s"}`, id)
}

func getPayload(r *http.Request) map[string]string {
	body := make(map[string]string)
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &body)

	return body
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT is not set.")
	}

	dns := os.Getenv("DATABASE_URL")

	if dns == "" {
		log.Fatal("$DATABASE_URL is not set.")
	}

	db, err := sql.Open("mysql", dns)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Connected to the database.")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/shorturl/new", func(w http.ResponseWriter, r *http.Request) {
		shortenerHandler(w, r, db)
	})
	http.HandleFunc("/api/shorturl/", func(w http.ResponseWriter, r *http.Request) {
		redirectHandler(w, r, db)
	})

	log.Print("Serving assets on /static/")
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Print("Starting server on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
