package services

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/projectdiscovery/gologger"
	"github.com/wahyuhadi/ESgo/es"
	"github.com/wahyuhadi/postman-fuzz/models"
)

// Model data

func Elastic(opts *models.Opts, data *models.Elastic) {

	cfg := elasticsearch.Config{
		Addresses: []string{opts.ElasticURI},
		Username:  opts.ElasticUser, // if ES need this
		Password:  opts.ElasticPass,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
				// ...
			},
		},
	}
	c, _ := elasticsearch.NewClient(cfg)
	// PushexamplePushData(c)
	pushdata(opts, data, c)

}

func pushdata(opts *models.Opts, datas *models.Elastic, c *elasticsearch.Client) {

	// parsing with esutil from elastic
	fmt.Println(datas)
	data := esutil.NewJSONReader(&datas)
	// Push data to elastic
	response, err := es.PushData(c, opts.ElasticIndex, data)
	if err != nil {
		gologger.Info().Str("Error", fmt.Sprintf("%v", err.Error())).Msg("Error push data")

	}
	gologger.Info().Str("Is Error ", fmt.Sprintf("%v", response.IsError())).Msg("Success Push data to elastic")

}
