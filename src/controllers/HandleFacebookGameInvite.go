package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleFacebookGameInvite(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]
	facebookId := request.FormValue("inviteeFacebookId")
	possession := request.FormValue("possession")
	playId := request.FormValue("playIdSelected")
	teamName := request.FormValue("teamName")
	flippedPlay := request.FormValue("flippedPlay")
	//Add What is on the lines
	whatsOnTheLine:= request.FormValue("whatsOnTheLine")

	success, message, status, gameId := models.FacebookInviteToGame(guid, facebookId, possession, teamName, playId, flippedPlay,whatsOnTheLine)
	output := jsonOutputs.GameOutput{Success: success, Message: message, Statuscode: status, Gameid: gameId}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))

}
