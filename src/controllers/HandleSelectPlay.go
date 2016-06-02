package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleSelectPlay(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]
	turnid := vars["turnid"]
	flippedPlay := request.FormValue("flip")
	playid := request.FormValue("playid")

	fmt.Println(guid + turnid + playid)

	success, message, status, animation := models.SelectPlay(guid, turnid, playid, flippedPlay)

	output := jsonOutputs.SelectPlayOutput{Success: success, Message: message, Statuscode: status, AnimationIds: animation}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))

}
