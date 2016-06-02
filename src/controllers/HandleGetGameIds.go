package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleGetGameIds(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]

	success, message, gameIds, status := models.GameIdsForUser(guid)

	outputJson, _ := json.Marshal(jsonOutputs.GameIds{Success: success, Message: message, GameIds: gameIds, Statuscode: status})
	fmt.Fprintf(res, string(outputJson))

}
