package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func search(querystr string) Response {
	es := es_conn()

	resp, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("sw"),
		es.Search.WithQuery(querystr),
		es.Search.WithPretty(),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Println("error during search:" + err.Error())
		log.Fatal(err)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	res := Response{}

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		p := People{
			source.(map[string]interface{})["name"].(string),
			source.(map[string]interface{})["films"].([]interface{}),
		}

		res.People = append(res.People, p)
	}

	res.Query = querystr
	return res

}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Path[len("/search/"):]
	res := search(query)
	err := templates.ExecuteTemplate(w, "search.html", &res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
