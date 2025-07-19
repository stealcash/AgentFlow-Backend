package db

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stealcash/AgentFlow/app/globals"
	"log"
)

var EsDB *elasticsearch.Client

func ElasticConnection() error {
	if !globals.Config.ElasticDatabase.RequiredElasticConnection {
		log.Println(" Elastic connection is disabled by config")
		return nil
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("%s:%s",
				globals.Config.ElasticDatabase.Host,
				globals.Config.ElasticDatabase.Port),
		},
		Username: globals.Config.ElasticDatabase.User,
		Password: globals.Config.ElasticDatabase.Password,
	}

	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	EsDB = esClient

	res, err := EsDB.Info()
	if err != nil {
		return fmt.Errorf("failed to get Elasticsearch info: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error connecting to Elasticsearch: %s", res.String())
	}

	log.Println("Connected to Elasticsearch")
	return nil
}
