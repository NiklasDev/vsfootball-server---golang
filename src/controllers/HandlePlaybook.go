package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandlePlaybook(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]

	success, message, playbook, status := models.GetPlaybook(guid)
	outputJson, _ := json.Marshal(jsonOutputs.PlaybookOutput{Success: success, Message: message, Playbook: playbook, Statuscode: status})
	fmt.Fprintf(res, string(outputJson))

}
