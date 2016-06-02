package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleRematch(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]
	inviteeGuid := request.FormValue("inviteeguid")
	possession := request.FormValue("possession")
	teamName := request.FormValue("teamName")
	playIdSelected := request.FormValue("playIdSelected")
	flippedPlay := request.FormValue("flippedPlay")
	//Add What is on the lines
	whatsOnTheLine:= request.FormValue("whatsOnTheLine")
	success, message, status, gameId := models.Rematch(guid, inviteeGuid, possession, teamName, playIdSelected, flippedPlay,whatsOnTheLine)

	output := jsonOutputs.GameOutput{Success: success, Message: message, Statuscode: status, Gameid: gameId}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))

}
