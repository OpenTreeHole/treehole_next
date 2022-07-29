package models

import (
	"log"
	"treehole_next/config"

	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client

func InitSearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{config.Config.SearchUrl},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	log.Println(elasticsearch.Version)
	log.Println(es.Info)
	ES = es
}
