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

func (h *Handler) WarrantySendRequest(warrantyParams models.WarrantyParams, orderItemUID uuid.UUID) (code int, body []byte, err error) {
	var (
		resp *http.Response
	)
	sendBody, _ := json.Marshal(warrantyParams)
	fmt.Println("warrantyParams = ", warrantyParams)
	urlString := h.conf.ServiceWarranty.URL + "/api/v1/warranty/" + orderItemUID.String() + "/warranty"
	fmt.Println(urlString)
	resp, err = utils.Post(urlString, sendBody)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	code = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)

	return
}
