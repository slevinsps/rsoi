package api

import (
	"encoding/json"
	"net/http"
	"services/internal/models"
	"services/utils"
)

func printResult(catched error, number int, place string) {
	if catched != nil {
		utils.PrintDebug("api/"+place+" failed(code:", number, "). Error message:"+catched.Error())
	} else {
		utils.PrintDebug("api/"+place+" success(code:", number, ")")
	}
}

func sendJSON(rw http.ResponseWriter, result []byte) {
	// bytes,_ := 	result.MarshalJSON()
	rw.Write(result)
	//json.NewEncoder(rw).Encode(result)
}

func sendMessage(rw http.ResponseWriter, messageText string, status int, err error) {
	rw.WriteHeader(status)
	message := models.Message{Message: messageText}
	if err != nil {
		additionalProps := &models.AdditionalProp{AdditionalProp1: err.Error()}
		message.Errors = additionalProps
	}
	resBytes, _ := json.Marshal(message)
	sendJSON(rw, resBytes)
}
