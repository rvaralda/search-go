package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/elastic/go-elasticsearch/v7"
)

type HttpResponse struct {
	url      string
	response *http.Response
	content  string
	id       int
	err      error
}

type People struct {
	Name  string        `json:"name"`
	Films []interface{} `json:"films"`
}

type Response struct {
	Query  string
	People []People `json:"people"`
}

func es_conn() *elasticsearch.Client {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	return es
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

var validPath = regexp.MustCompile("^/(index|search)/*([a-zA-Z0-9]*)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func main() {
	log.Println("Running")
	http.HandleFunc("/search/", makeHandler(searchHandler))
	http.HandleFunc("/index", makeHandler(indexHandler))

	log.Fatal(http.ListenAndServe("localhost:9000", nil))
}
