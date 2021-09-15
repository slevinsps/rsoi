package utils

import (
	"bytes"
	"net/http"
)

// Post
func Post(urlString string, body []byte) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", urlString, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	return
}

// Delete
func Delete(urlString string, body []byte) (resp *http.Response, err error) {
	req, err := http.NewRequest("DELETE", urlString, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	return
}

// Get
func Get(urlString string, body []byte) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", urlString, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	return
}
