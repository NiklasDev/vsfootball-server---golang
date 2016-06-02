package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleRandomOpponent(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]
	success, message, status, matchedGuid, animationIds := models.RandomOpponentQueue(guid)
	output := jsonOutputs.RandomOpponentOutput{success, message, status, matchedGuid, animationIds}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
