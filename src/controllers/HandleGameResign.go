package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleGameResign(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]

	gameid := request.FormValue("gameid")
	success, message, status := models.ResignGame(guid, gameid)
	output := jsonOutputs.SuccessMessage{success, message, status}

	outputJson, _ := json.Marshal(output)

	fmt.Fprintf(res, string(outputJson))

}
