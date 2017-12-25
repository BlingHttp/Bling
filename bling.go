package bling

import "net/http"

func New() *Bling {
	return &Bling{
		HTTPClient: http.DefaultClient,
	}
}

func (b *Bling) Client(client HTTPClient) *Bling {
	b.HTTPClient = client
	return b
}

func (b *Bling) Verb(verb string) *Request {
	return NewRequest(b.HTTPClient, verb)
}

func (b *Bling) Get(pathUrl string) *Request {
	return b.Verb("GET").Path(pathUrl)
}

func (b *Bling) Post(pathUrl string) *Request {
	return b.Verb("POST").Path(pathUrl)
}

func (b *Bling) Put(pathUrl string) *Request {
	return b.Verb("PUT").Path(pathUrl)
}

func (b *Bling) Patch(pathUrl string) *Request {
	return b.Verb("PATCH").Path(pathUrl)
}

func (b *Bling) Delete(pathUrl string) *Request {
	return b.Verb("DELETE").Path(pathUrl)
}

func (b *Bling) Head(pathUrl string) *Request {
	return b.Verb("HEAD").Path(pathUrl)
}
