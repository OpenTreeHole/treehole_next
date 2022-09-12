package config

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client

func InitSearch() {
	if Config.Mode == "test" || Config.Mode == "bench" {
		return
	}

	// export ELASTICSEARCH_URL environment variable to set the ElasticSearch URL
	// example: http://user:pass@127.0.0.1:9200
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
	ES = es
}
