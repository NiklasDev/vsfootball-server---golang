package controllers

import (
	"encoding/json"
	"fmt"
	"models"
	"net/http"
)

func HandleGcmTest(res http.ResponseWriter, request *http.Request) {
	deviceId := request.FormValue("deviceid")
	message := request.FormValue("message")
	data := map[string]string{"gameid": "123"}
	payload, _ := json.Marshal(data)
	models.SendGCMMessage(deviceId, message, string(payload))
	fmt.Fprintf(res, "Test if you have push.")
}
