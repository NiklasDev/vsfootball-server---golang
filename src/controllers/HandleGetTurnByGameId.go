package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleTurnsByGame(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]
	gameid := vars["gameid"]
	returnLimit := vars["returnlimit"]

	success, message, game, status := models.TurnsForGame(guid, gameid, returnLimit)

	outputJson, _ := json.Marshal(jsonOutputs.GameTurnOutput{Success: success, Message: message, Game: game, Statuscode: status})
	fmt.Fprintf(res, string(outputJson))

}
