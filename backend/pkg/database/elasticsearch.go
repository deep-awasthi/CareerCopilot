package database

import (
	"fmt"

	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/logger"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

var ESClient *elasticsearch.Client

func InitElasticsearch(cfg *config.ElasticsearchConfig) (*elasticsearch.Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
	}

	if cfg.Username != "" {
		esCfg.Username = cfg.Username
		esCfg.Password = cfg.Password
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch returned error: %s", res.Status())
	}

	ESClient = client
	logger.Info("Elasticsearch connected", zap.Strings("addresses", cfg.Addresses))
	return client, nil
}
