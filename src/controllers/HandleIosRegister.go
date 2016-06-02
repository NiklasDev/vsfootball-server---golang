package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleIosRegister(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]
	deviceToken := request.FormValue("devicetoken")

	success, message, status := models.RegisterIosDevice(guid, deviceToken)

	output := jsonOutputs.SuccessMessage{Success: success, Message: message, Statuscode: status}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))

}
