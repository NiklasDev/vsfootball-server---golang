package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleAndroidRegister(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]
	registrationid := request.FormValue("registrationid")

	success, message, status := models.RegisterAndroidDevice(guid, registrationid)

	output := jsonOutputs.SuccessMessage{Success: success, Message: message, Statuscode: status}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))

}
