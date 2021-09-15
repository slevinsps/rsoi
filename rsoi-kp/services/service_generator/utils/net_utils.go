package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// Post
func SendRequest(urlString string, r *http.Request, method string) (resp *http.Response, err error) {
	var (
		body []byte
	)
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	req, err := http.NewRequest(method, urlString, bytes.NewBuffer(body))
	req.Header = r.Header
	// req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	return
}
