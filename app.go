package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type ShortUrl struct {
	OriginalUrl string `json:"original_url"`
	ShortUrl    int    `json:"short_url"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("views/index.html"))
	t.Execute(w, nil)
}

func shortenerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Allow-Access-Control-Origin", "*")

	url := &ShortUrl{
		OriginalUrl: "https://example.com/",
		ShortUrl:    1,
	}

	j, err := json.Marshal(url)

	if err != nil {
		http.Error(w, `{"message": "Something went wrong"}`, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(j))
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT is not set.")
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/shorturl/new", shortenerHandler)

	log.Print("Serving assets on /static/")
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Print("Starting server on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
