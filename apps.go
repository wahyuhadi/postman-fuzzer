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
)

func initOps() *models.Opts {
	flag.Parse()
	return &models.Opts{
		Location:    *location,
		Proxy:       *proxy,
		KeyHeader:   *keyheader,
		ValueHeader: *valueheader,
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
