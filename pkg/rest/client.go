package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

var transport *AopTransport
var client *http.Client

func InitRestClient() {
	transport = newTransport()
	client = newHttpClient(transport)
}

func AddRequestProxy(p RequestProxy) {
	transport.proxies = append(transport.proxies, p)
}

func newTransport() *AopTransport {
	return &AopTransport{underlying: http.DefaultTransport, proxies: make([]RequestProxy, 0)}
}

type AopTransport struct {
	underlying http.RoundTripper
	proxies    []RequestProxy
}

func (t *AopTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	var err error
	for _, p := range t.proxies {
		r, err = p.Before(r.Context(), r)
		if err != nil {
			return nil, err
		}
	}

	var resp *http.Response
	resp, err = t.underlying.RoundTrip(r)

	for _, p := range t.proxies {
		resp, err = p.After(r.Context(), resp, err)
	}

	return resp, err
}

func newHttpClient(transport http.RoundTripper) *http.Client {
	return &http.Client{Timeout: time.Second * 10, Transport: transport}
}

func GetRequest(ctx context.Context, url string, response interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)

	htmlResp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer htmlResp.Body.Close()

	decoder := json.NewDecoder(htmlResp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return err
	}
	return nil
}
