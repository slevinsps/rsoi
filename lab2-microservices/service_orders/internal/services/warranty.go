package api

import (
	"fmt"
	"net/http"
	"services/utils"

	uuid "github.com/satori/go.uuid"
)

func (h *Handler) startWarranty(itemUID uuid.UUID) (err error) {
	var (
		resp *http.Response
	)

	urlString := h.conf.ServiceWarranty.URL + "/api/v1/warranty/" + itemUID.String()
	fmt.Println(urlString)
	resp, err = utils.Post(urlString, []byte{})
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return
}

func (h *Handler) stopWarranty(itemUID uuid.UUID) (code int, err error) {
	var (
		resp *http.Response
	)

	urlString := h.conf.ServiceWarranty.URL + "/api/v1/warranty/" + itemUID.String()
	fmt.Println(urlString)

	resp, err = utils.Delete(urlString, []byte{})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	code = resp.StatusCode

	return
}
