package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Url     string
	Params  map[string]string
	Headers map[string]string
}

func New(url string) *Request {
	return &Request{
		Url:     url,
		Params:  make(map[string]string),
		Headers: make(map[string]string),
	}
}

func (req *Request) Get() ([]byte, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	url := req.Url
	if len(req.Params) > 0 {
		url = fmt.Sprint(url, "?")
		for key, value := range req.Params {
			url = fmt.Sprint(url, key, "=", value, "&")
		}
		url = strings.TrimSuffix(url, "&")
	}

	fmt.Printf("Total url: %v\n", url)

	r, err := client.Get(url)
	if err != nil {
		err = fmt.Errorf("Request to url (%s) is failed. Error: %w\n", url, err)
		return nil, err
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("Reading of response is failed. Error: %w\n", err)
		return nil, err
	}

	return body, nil
}

func (req *Request) Post(body []byte) ([]byte, error) {
	r, err := http.NewRequest("POST", req.Url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("Can't create post request to %s. Error: %v\n", req.Url, err)
	}
	if len(req.Headers) > 0 {
		for key, value := range req.Headers {
			r.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("Post request is failed. Url: %s\n, Body: %v\n, Headers: %v\n, Error: %v\n", req.Url, body, req.Headers, err)
	}

	defer res.Body.Close()

	resbody, err := io.ReadAll(res.Body)
	//StatusCreated = 201
	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Post request has unsuccessfull status code. Url: %s\n, Body: %v\n, Headers: %v\n, Status: %s, Response: %x\n", req.Url, body, req.Headers, res.Status, resbody)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed reading post response. Url: %s\n, Body: %v\n, Headers: %v\n, Status: %s, Response: %x, Error: %v\n", req.Url, body, req.Headers, res.Status, resbody, err)
	}

	return resbody, nil
}
