package config

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client

func InitSearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{Config.SearchUrl},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	log.Println(elasticsearch.Version)
	log.Println(es.Info)
	ES = es
}
