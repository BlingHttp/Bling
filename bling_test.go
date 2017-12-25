package bling

import (
	"testing"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"net/http"
)

func assertMethod(t *testing.T, expectedMethod string, req *http.Request) {
	if actualMethod := req.Method; actualMethod != expectedMethod {
		t.Errorf("expected method %s, got %s", expectedMethod, actualMethod)
	}
}

func assertBody(t *testing.T, expectedBody string, req *http.Request) {
	actualBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Errorf("expected nil error , got: %v", err)
	}
	if string(actualBody) != expectedBody {
		t.Errorf("expected body %s, got %s", actualBody, expectedBody)
	}
}

func testServer() (*http.Client, *http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	client := &http.Client{Transport: transport}
	return client, mux, server
}

func TestGetBasicUsage(t *testing.T) {
	client, mux, server := testServer()
	defer server.Close()
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "okay", "data": "bling"}`)
	})

	blingClient := New().Client(client)
	resp, err := blingClient.Get("http://google.com/get").DoRaw()
	if err != nil {
		t.Errorf("expected nil error , got: %v", err)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected nil error , got: %v", err)
	}
	fmt.Println(string(responseBody))
	return
}

func TestPostBasicUsage(t *testing.T) {
	client, mux, server := testServer()
	bodyString := `{"status": "okay", "data": "bling"}`
	defer server.Close()
	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		assertBody(t, bodyString, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, bodyString)
	})

	blingClient := New().Client(client)
	resp, err := blingClient.Post("http://example.com").Path("post").Body([]byte(bodyString)).DoRaw()
	if err != nil {
		t.Errorf("expected nil error , got: %v", err)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected nil error , got: %v", err)
	}
	fmt.Println(string(responseBody))
	return
}

func TestPathSetter(t *testing.T) {
	cases := []struct {
		base           string
		path           string
		expectedRawUrl string
	}{
		{"http://google.com/", "foo", "http://google.com/foo"},
		{"http://google.com/", "/foo", "http://google.com/foo"},
		{"http://google.com", "foo", "http://google.com/foo"},
		{"http://google.com", "/foo", "http://google.com/foo"},
		{"http://google.com/foo/", "bar", "http://google.com/foo/bar"},
		// base should end in trailing slash if it is to be Path extended
		{"http://google.com/foo", "bar", "http://google.com/bar"},
		{"http://google.com/foo", "/bar", "http://google.com/bar"},
		// path extension is absolute
		{"http://google.com", "http://bling.com/", "http://bling.com/"},
		{"http://google.com/", "http://bling.com/", "http://bling.com/"},
		{"http://google.com", "http://bling.com", "http://bling.com"},
		{"http://google.com/", "http://bling.com", "http://bling.com"},
		// empty base, empty path
		{"", "http://bling.com", "http://bling.com"},
		{"http://google.com", "", "http://google.com"},
		{"", "", ""},
	}
	for _, c := range cases {
		client := New().Get(c.base).Path(c.path)
		if client.rawUrl != c.expectedRawUrl {
			t.Errorf("expected %s, got %s", c.expectedRawUrl, client.rawUrl)
		}
	}
}

func TestMethodSetters(t *testing.T) {
	cases := []struct {
		blingRequest   *Request
		expectedMethod string
	}{
		{New().Get("http://google.com"), "GET"},
		{New().Post("http://google.com"), "POST"},
		{New().Put("http://google.com"), "PUT"},
		{New().Delete("http://google.com"), "DELETE"},
	}
	for _, c := range cases {
		if c.blingRequest.verb != c.expectedMethod {
			t.Errorf("expected method %s, got %s", c.expectedMethod, c.blingRequest.verb)
		}
	}
}

func TestRequest_Do(t *testing.T) {
	client, mux, server := testServer()
	defer server.Close()
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "okay", "data": "bling"}`)
	})

	blingClient := New().Client(client)
	result := blingClient.Get("http://google.com/get").Do()
	if result.err != nil {
		t.Errorf("result should not return err")
	}
	if string(result.body) != `{"status": "okay", "data": "bling"}` {
		t.Errorf("result body should match expected body")
	}
	if result.statusCode != http.StatusOK {
		t.Errorf("result status code should be 200")
	}
}
