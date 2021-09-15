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

func (h *Handler) GetItemInWarehouse(itemUID uuid.UUID) (code int, item models.Item, err error) {
	var (
		resp *http.Response
		body []byte
	)

	urlString := h.conf.ServiceWarehouse.URL + "/api/v1/warehouse/" + itemUID.String()
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
	err = json.Unmarshal(body, &item)

	return
}
