package web

import (
	"io"
	"net/http"
)

type HttpBody struct {
	Key   string
	Value string
}

var client = &http.Client{}

type HttpHeaders struct {
	Cookie      *http.Cookie
	ContentType string
}

func SendRequest(method string, url string, headers *HttpHeaders, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		if headers.Cookie != nil {
			req.AddCookie(headers.Cookie)
		}

		req.Header.Add("Content-Type", headers.ContentType)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
