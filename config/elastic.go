package config

import (
	"github.com/elastic/go-elasticsearch/v8"
)

func NewAppElastic() *elasticsearch.Client {
	// Connect elastic
	cfg := elasticsearch.Config{
		Addresses: []string{
			Env.ElasticHost,
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	res, err := es.Ping()
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic(res)
	}

	return es
}
