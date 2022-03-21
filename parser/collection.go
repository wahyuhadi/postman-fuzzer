package parser

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/rbretecher/go-postman-collection"
	"github.com/wahyuhadi/postman-fuzz/models"
)

type Header struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Disabled    bool   `json:"disabled,omitempty"`
	Description string `json:"description,omitempty"`
}

func ParseCollection(options *models.Opts) (reqmodels *[]models.Request, err error) {
	file, err := os.Open(options.Location)
	if err != nil {
		panic(err)
	}

	c, err := postman.ParseCollection(file)
	if err != nil {
		return nil, errors.New("Error when parse collection")
	}

	var arr []models.Request
	for _, col := range c.Items {
		// Generate new http request
		req, err := http.NewRequest(string(col.Request.Method), col.Request.URL.String(), nil)
		if err != nil {
			log.Printf("Error sending request to API endpoint. %+v", err)
			return nil, errors.New("Error")
		}
		u, _ := url.Parse(col.Request.URL.String())
		req.Proto = col.Request.URL.Protocol
		req.URL = u
		req.Host = u.Host
		for _, header := range col.Request.Header {
			req.Header.Set(header.Key, header.Value)
		}
		value := models.Request{Req: req, URI: &col.Request.URL.Raw}
		if col.Request.Body.Raw != "" {
			value.Body = &col.Request.Body.Raw
			value.IsHaveBody = true                         // if this have a body
			value.IsBodyJson = isJSON(col.Request.Body.Raw) // check if json
		}
		arr = append(arr, value)
	}

	defer file.Close()

	return &arr, nil
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
