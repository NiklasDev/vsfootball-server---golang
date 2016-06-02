package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleGetPlayForPlayId(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	playId := vars["playid"]

	success, message, status, play := models.GetPlayForPlayId(playId)

	outputJson, _ := json.Marshal(jsonOutputs.PlaybookOutput{Success: success, Message: message, Playbook: play, Statuscode: status})
	fmt.Fprintf(res, string(outputJson))
}
