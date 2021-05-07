# Search-go
---

Sample search application written in go using Elasticsearch.

It consumes https://swapi.dev to index, and return star wars characters with movies.

### To run

```bash
$ export ELASTICSEARCH_URL=http://elasticsearch:9200
 
$ go build
$ ./search-go

```

### To run with docker, use

```bash
$ docker network create search-network
$ docker run -d --net search-network -e "discovery.type=single-node" -p 9200:9200 -p 9300:9300 elasticsearch:7.12.0
$ docker run --name search-go --net search-network -d -p 9000:9000 search-go
```

#### Checking logs
```bash
$ docker logs -f search-go
```

### Usage

#### Indexing
```bash
$ curl localhost:9000/index
```

#### Searching
```bash
$ curl localhost:9000/search/luke
```
