package bling

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/json"
)

func NewRequest(client HTTPClient, verb string) *Request {
	r := &Request{
		client: client,
		verb:   verb,
	}
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Set(key, value)
	return r
}

func (r *Request) Path(path string) *Request {
	baseUrl, baseErr := url.Parse(r.rawUrl)
	pathUrl, pathErr := url.Parse(path)
	if baseErr == nil && pathErr == nil {
		r.rawUrl = baseUrl.ResolveReference(pathUrl).String()
		return r
	}
	return r
}

func (r *Request) URL() string {
	return ""
}

func (r *Request) Body(obj interface{}) *Request {
	if r.err != nil {
		return r
	}
	switch t := obj.(type) {
	case string:
		//try to read from a file
		data, err := ioutil.ReadFile(t)
		if err != nil {
			r.err = err
			return r
		}
		r.body = bytes.NewReader(data)
	case []byte:
		r.body = bytes.NewReader(t)
	case io.Reader:
		r.body = t
	default:
		r.err = fmt.Errorf("unknown type used for body: %+v", obj)
	}
	return r
}

func (r *Request) Do() Result {
	var result Result
	err := r.performRequest(func(req *http.Request, resp *http.Response) {
		result = r.transformResponse(resp, req)
	})
	if err != nil {
		return Result{err: err}
	}
	return result
}

func (r *Request) transformResponse(resp *http.Response, req *http.Request) Result {
	var body []byte
	if resp.Body != nil {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			unexpectedErr := fmt.Errorf("Unexpected error %#v when reading response body. Please retry.", err)
			return Result{
				err: unexpectedErr,
			}
		}
		body = data
	}
	return Result{
		body:        body,
		contentType: resp.Header.Get("Content-Type"),
		statusCode:  resp.StatusCode,
	}
}

func (r *Request) DoRaw() (*http.Response, error) {
	var _resp *http.Response
	err := r.performRequest(func(request *http.Request, response *http.Response) {
		_resp = response
	})
	if err != nil {
		return nil, err
	}
	return _resp, err
}

func (r *Request) GetHttpRequest() (*http.Request, error) {
	reqURL, err := url.Parse(r.rawUrl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(r.verb, reqURL.String(), r.body)
	if err != nil {
		return nil, err
	}
	req.Header = r.headers
	return req, err
}

func (r *Request) performRequest(fn func(*http.Request, *http.Response)) error {
	//Retry logic should go here
	client := r.client
	if client == nil {
		client = http.DefaultClient
	}
	req, err := r.GetHttpRequest()
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	fn(req, resp)
	return nil
}

func (r Result) StatusCode(statusCode *int) Result {
	*statusCode = r.statusCode
	return r
}

func (r Result) Into(successObj, errorObj interface{}) error {
	if r.Has2xxStatus() {
		if successObj != nil {
			return json.Unmarshal(r.body, successObj)
		}
	} else {
		if errorObj != nil {
			return json.Unmarshal(r.body, errorObj)
		}
	}
	return nil
}

func (r Result) Has2xxStatus() bool {
	return 200 <= r.statusCode && r.statusCode <= 299
}
