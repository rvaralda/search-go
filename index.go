package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
)

func index(es *elasticsearch.Client, id int, data string) {
	resp, err := es.Index(
		"sw",
		strings.NewReader(data),
		es.Index.WithDocumentID(strconv.Itoa(id)),
		es.Index.WithRefresh("true"),
		es.Index.WithPretty(),
		es.Index.WithFilterPath("result", "_id"),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
	} else {
		log.Printf("[%s] %s", resp.Status(), r["result"])
	}

}

func sw_get() int {
	es := es_conn()
	peopleNumber := 83
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}

	for i := 1; i <= peopleNumber; i++ {
		time.Sleep(1)
		url := "https://swapi.dev/api/people/" + strconv.Itoa(i)
		go func(url string, i int) {
			log.Printf("Fetching %s \n", url)
			resp, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}

			b, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			ch <- &HttpResponse{url, resp, string(b[:]), i, err}
		}(url, i)
	}

	for {
		select {
		case r := <-ch:
			log.Println("Done case: " + strconv.Itoa(r.id))
			log.Printf("Indexing %s, %d\n", r.url, r.id)
			index(es, r.id, r.content)

			responses = append(responses, r)
			if len(responses) == peopleNumber {
				return peopleNumber
			}
		}
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var indexes int
	indexes = sw_get()

	result := fmt.Sprintf("%d Indexes Ok\n", indexes)
	err := templates.ExecuteTemplate(w, "index.html", result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
