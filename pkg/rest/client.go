package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jfbramlett/go-aop/pkg/tracing"
	"github.com/jfbramlett/go-aop/pkg/web"
)

func NewTransport() *AopTransport {
	return &AopTransport{underlying: http.DefaultTransport}
}

func NewTransportFor(roundTripper http.RoundTripper) *AopTransport {
	return &AopTransport{underlying: roundTripper}
}

type AopTransport struct {
	underlying http.RoundTripper
}

func (t *AopTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	requestId := tracing.GetTraceFromContext(r.Context())
	r.Header.Add(web.HeaderRequestId, requestId)
	return t.underlying.RoundTrip(r)
}

func NewHttpClient() *http.Client {
	return &http.Client{Timeout: time.Second * 10, Transport: NewTransport()}
}

func HttpGetRequest(ctx context.Context, url string, response interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)

	client := NewHttpClient()

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
