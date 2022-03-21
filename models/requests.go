package models

import "net/http"

type Request struct {
	Req        *http.Request
	URI        *string
	IsHaveBody bool
	Body       *string
	IsBodyJson bool
}
