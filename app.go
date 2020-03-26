package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
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

	payload := getPayload(r)

	url := &ShortUrl{
		OriginalUrl: payload["url"],
		ShortUrl:    1,
	}

	j, err := json.Marshal(url)

	if err != nil {
		http.Error(w, `{"message": "Something went wrong"}`, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(j))
}

func getPayload(r *http.Request) map[string]string {
	body := make(map[string]string)
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &body)

	return body
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Allow-Access-Control-Origin", "*")

	id := r.URL.Path[len("/api/shorturl/"):]

	fmt.Fprintf(w, `{"message": "hello, %s"}`, id)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT is not set.")
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/shorturl/new", shortenerHandler)
	http.HandleFunc("/api/shorturl/", redirectHandler)

	log.Print("Serving assets on /static/")
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Print("Starting server on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
