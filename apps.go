package main

import (
	"flag"
	"log"

	"github.com/wahyuhadi/postman-fuzz/fuzz"
	"github.com/wahyuhadi/postman-fuzz/models"
	"github.com/wahyuhadi/postman-fuzz/parser"
)

var (
	location    = flag.String("loc", "", "Location postman collection")
	proxy       = flag.String("p", "", "Proxy server")
	keyheader   = flag.String("key", "", "key headers")
	valueheader = flag.String("value", "", "value headers")
	elastic     = flag.Bool("elastic", false, "push to elastic")
	elastricURI = flag.String("elasticurl", "http://127.0.0.1:9200", "elastic url")
)

func initOps() *models.Opts {
	flag.Parse()
	return &models.Opts{
		Location:     *location,
		Proxy:        *proxy,
		KeyHeader:    *keyheader,
		ValueHeader:  *valueheader,
		Elastic:      *elastic,
		ElasticURI:   *elastricURI,
		ElasticUser:  "tes",
		ElasticPass:  "tes",
		ElasticIndex: "ngetest",
	}
}

func main() {
	options := initOps()
	reqs, err := parser.ParseCollection(options)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	fuzz.DoFuzz(options, reqs)

}
