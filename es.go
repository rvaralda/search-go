package main

import (
    "net/http"
    "log"
    "fmt"
    "io/ioutil"
    "strconv"
    "encoding/json"
    "os"
    "html/template"
    "regexp"

    elastigo "github.com/mattbaird/elastigo/lib"
)

type HttpResponse struct {
  url      string
  response *http.Response
  content  string
  id       int
  err      error
}

type People struct {
    Name      string  `json:"name"`
    Films   []string  `json:"films"`
}

type Response struct {
    Query   string
    People  []People `json:"people"`
}

//Elastic Search connectior
func es_conn() ( *elastigo.Conn ) {
    c := elastigo.NewConn()
    c.Domain = os.Getenv("ES_HOST")
    c.Port = os.Getenv("ES_PORT")
    return c
}

//index
func index(id int, data string) {
    c := es_conn()

    response, err := c.Index("sw", "people", strconv.Itoa(id), nil, data)
    if err != nil {
        log.Fatal(err)
    }
    c.Flush()
    log.Printf("Index %d %v", id, response)

}

func sw_get() []*HttpResponse {
    urls := 87
    ch := make(chan *HttpResponse)
    responses := []*HttpResponse{}
    for i := 1; i <= urls; i++ {
        url := "http://swapi.co/api/people/" + strconv.Itoa(i)
        go func(url string, i int) {
            fmt.Printf("Fetching %s \n", url)
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
            
            fmt.Printf("Index %s, %d\n", r.url, r.id)
            index(r.id, r.content)

            responses = append(responses, r)
            if len(responses) == urls {
                return responses
            }
        }
    }
    return responses

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    sw_get()
    err := templates.ExecuteTemplate(w, "index.html", "Index Ok")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

//search
func search(query string) Response {
    c := es_conn()

    searchJson := `{
        "query": {
            "fuzzy_like_this_field" : {
             "name" : {
                "like_text": "` + query +  `", "max_query_terms": 5
             }
            }
        }
    }`
    
    searchresponse, err := c.Search("sw", "people", nil, searchJson)
    if err != nil {
        log.Println("error during search:" + err.Error())
        log.Fatal(err)
    }

    res := Response{}
    p   := People{}
    res.Query = query

    for _, response := range searchresponse.Hits.Hits {
        bytes, err := response.Source.MarshalJSON()
        if err != nil {
            log.Fatal(err)
        }
        json.Unmarshal(bytes, &p)
        res.People = append(res.People, p)
        // fmt.Println(res)
    }

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

var templates = template.Must(template.ParseGlob("templates/*.html"))


var validPath = regexp.MustCompile("^/(index|search)/*([a-zA-Z0-9]*)$")

func makeHandler(fn func (http.ResponseWriter, *http.Request)) http.HandlerFunc {
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
    fmt.Println("Running")
    http.HandleFunc("/search/", makeHandler(searchHandler))
    http.HandleFunc("/index", makeHandler(indexHandler))
    err := http.ListenAndServe(":9000", nil)
    if err != nil {
        log.Fatal(err)
    }
}