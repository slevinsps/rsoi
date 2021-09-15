package api

import (
	"net/http"
	"restService/internal/utils"
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
