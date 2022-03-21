package fuzz

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.com/michenriksen/jdam/pkg/jdam"
	"gitlab.com/michenriksen/jdam/pkg/jdam/mutation"

	"github.com/wahyuhadi/postman-fuzz/models"
)

var localCertFile = "/home/samalas/.key.pem"

func DoFuzz(opt *models.Opts, reqs *[]models.Request) {
	for _, req := range *reqs {

		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}

		// Read in the cert file
		certs, err := ioutil.ReadFile(localCertFile)
		if err != nil {
			log.Fatalf("Failed to append %q to RootCAs: %v", localCertFile, err)
		}

		// Append our cert to the system pool
		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			log.Println("No certs appended, using system certs only")
		}

		// Trust the augmented cert pool in our client
		config := &tls.Config{
			InsecureSkipVerify: false,
			RootCAs:            rootCAs,
		}
		tr := &http.Transport{TLSClientConfig: config}
		proxyStr := opt.Proxy
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			log.Println(err)
		}
		log.Println("Do Fuzz in : ", req.Req.Method, *req.URI)
		if opt.Proxy != "" {
			tr.Proxy = http.ProxyURL(proxyURL)
		}
		//adding the Transport object to the http Client
		client := &http.Client{
			Transport: tr,
		}

		req.Req.Header.Set("User-Agent", "Please Plant More Trees")
		if opt.KeyHeader != "" {
			req.Req.Header.Set(opt.KeyHeader, opt.ValueHeader)
		}

		// Do fuzz body json
		if req.Body != nil && req.IsBodyJson {
			data := (*req.Body)
			subject := map[string]interface{}{}
			err := json.Unmarshal([]byte(data), &subject)
			if err != nil {
				panic(err)
			}
			req.Req.Header.Set("Content-Type", "application/json")
			// Create a new fuzzer with all available mutators.
			// We instruct the fuzzer to ignore the id field so that its
			// value is never changed. This is useful in cases where you know
			// that altering this value would result in uninteresting errors.
			fuzzer := jdam.New(mutation.Mutators).MaxDepth(1000).NilChance(1)
			for i := 0; i < 1500; i++ {
				// Fuzz a random field with a random mutator.
				fuzzed := fuzzer.Fuzz(subject)

				// Encode the fuzzed object into JSON.
				fuzzedJSON, err := json.Marshal(fuzzed)
				if err != nil {
					panic(err)
				}

				req.Req.Body = io.NopCloser(strings.NewReader(string(fuzzedJSON)))
				// Send request to the server.
				resp, err := client.Do(req.Req)
				if err != nil {
					log.Fatal(err)
					return
				}
				if resp.StatusCode != 200 {
					defer resp.Body.Close()
					// Our payload has caused some sort of internal server error!
					// Write the payload to a file for further research.
					b, err := io.ReadAll(resp.Body)
					// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
					if err != nil {
						log.Fatalln(err)
					}

					save := fmt.Sprintf("\nRequest:\n%v\n\nResponses:\n%v", string(fuzzedJSON), string(b))
					ioutil.WriteFile(fmt.Sprintf("crash/error-%v.json", i), []byte(save), 0644)
				}

				// Sleep for a bit to be nice to the server (lol).
				time.Sleep(200 * time.Millisecond)

			}

		}

		client.Do(req.Req)
	}
}
