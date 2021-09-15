package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"services/internal/models"
	"services/utils"

	uuid "github.com/satori/go.uuid"
)

func (h *Handler) CheckWarrantyStatusInWarranty(itemUID uuid.UUID) (code int, warranty models.Warranty, err error) {
	var (
		resp *http.Response
		body []byte
	)

	urlString := h.conf.ServiceWarranty.URL + "/api/v1/warranty/" + itemUID.String()
	fmt.Println(urlString)
	resp, err = utils.Get(urlString, []byte{})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	code = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)

	if code != http.StatusOK {
		return
	}
	err = json.Unmarshal(body, &warranty)
	return
}
