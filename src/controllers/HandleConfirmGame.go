package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleConfirmGame(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]

	gameid := request.FormValue("gameid")
	teamname := request.FormValue("teamname")

	success, message, status := models.ConfirmGame(guid, gameid, teamname)
	output := jsonOutputs.SuccessMessage{success, message, status}

	jsonOutput, _ := json.Marshal(output)
	fmt.Fprintf(res, string(jsonOutput))
}
