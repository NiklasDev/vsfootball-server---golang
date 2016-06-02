package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleChangePassword(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]
	oldPassword := request.FormValue("oldpassword")
	newPassword := request.FormValue("newpassword")
	success, message, status := models.ChangePassword(guid, oldPassword, newPassword)
	output := jsonOutputs.SuccessMessage{Success: success, Message: message, Statuscode: status}

	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
