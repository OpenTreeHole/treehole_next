package elastic

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/goccy/go-json"
	"io"
	"log"
	"strings"

	. "treehole_next/config"
	. "treehole_next/models"
)

var ES *elasticsearch.Client

const IndexName = "floors"

func Init() {
	if Config.Mode == "test" || Config.Mode == "bench" || Config.ElasticsearchUrl == "" {
		return
	}

	// export ELASTICSEARCH_URL environment variable to set the ElasticSearch URL
	// example: http://user:pass@127.0.0.1:9200
	var err error
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{Config.ElasticsearchUrl},
	})
	if err != nil {
		log.Fatalf("Error creating elasticsearch client: %s", err)
	}

	res, err := ES.Info()
	if err != nil {
		log.Fatalf("Error getting elasticsearch response: %s", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	var r Map
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the elasticsearch response body: %s", err.Error())
	}

	// print Client and Server Info
	log.Printf("elasticsearch Client: %s\n", elasticsearch.Version)
	log.Printf("elasticsearch Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))
	ES = es
}
