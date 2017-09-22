package bling

import (
	"net/http"
	"io"
)

// HTTPClient is an interface for testing a request object.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type CheckRetry func(resp *http.Response, err error) (bool, error)

type Bling struct {
	HTTPClient HTTPClient
	//Logger     *log.Logger

	//CheckRetry CheckRetry
}

//all the backoff ans try managers will go here
type Request struct {
	client      HTTPClient
	verb        string
	rawUrl      string
	headers     http.Header
	queryParams []interface{}
	body        io.Reader
	err         error
}

type Result struct {
	body        []byte
	contentType string
	err         error
	statusCode  int
}
