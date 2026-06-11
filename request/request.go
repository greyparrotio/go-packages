package request

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

/*
Request sends a http request to the specified url with the accompanying data as the payload and the specified method as
the http verb of the http request.

*/
func Request(ctx context.Context, url string, data io.Reader, method string, header, query map[string]string, disableSSL bool) (*http.Response, error) {
	client := &http.Client{
		Transport: otelhttp.NewTransport(
			&http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: disableSSL,
				},
			}),
	}
	request, err := http.NewRequestWithContext(ctx, method, url, data)
	if err != nil {
		return nil, err
	}

	if len(header) >= 1 {
		for key, value := range header {
			request.Header.Set(key, value)
		}
	}
	if len(query) >= 1 {
		requestQuery := request.URL.Query()
		for key, value := range query {
			requestQuery.Add(key, value)
		}
		request.URL.RawQuery = requestQuery.Encode()
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
