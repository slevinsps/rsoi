package api

import (
	"io/ioutil"
	"net/http"
	"services/utils"
)

func (h *Handler) SendRequest(r *http.Request, serviceBase string, urlString string, method string) (code int, body []byte, err error) {
	var (
		resp  *http.Response
		token string
	)

	if token, err = h.redisClient.Get(serviceBase).Result(); err != nil {
		if token, err = h.RegisterService(serviceBase); err != nil {
			return
		}
	}
	r.Header.Set("Authorization", "B "+token)

	if resp, err = utils.SendRequest(serviceBase+urlString, r, method); err != nil {
		return
	}
	defer resp.Body.Close()
	code = resp.StatusCode
	if code == http.StatusUnauthorized {
		if token, err = h.RegisterService(serviceBase); err != nil {
			return
		}
		r.Header.Set("Authorization", "B "+token)
		if resp, err = utils.SendRequest(serviceBase+urlString, r, method); err != nil {
			return
		}
	}
	body, err = ioutil.ReadAll(resp.Body)

	return
}
