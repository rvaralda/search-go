#Search-go

Search example using Elasticsearch written in GO.

It consumes https://swapi.co/ to index, and return sw caracters with movies.

To run:

```
$ export ES_HOST=<elasticsearch_host>
$ export ES_PORT=<elasticsearch_port>
 
$ go get -d
$ go build -o es es.go
$ ./es

```

To run with docker, use:

```
$ docker run -e "ES_HOST=<elasticsearch_host> -e "ES_PORT=9200" -d -p 9000:9000 \
--name search-go search-go
```