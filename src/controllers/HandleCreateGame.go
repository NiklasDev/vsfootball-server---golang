package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
	"strings"
)

func HandleCreateGame(res http.ResponseWriter, request *http.Request) {
	// func CreateGame(guid,inviteEmail, possession, teamName, playIdSelected string) (string,bool)
	//func CreateGame(guid,inviteEmail, possession, teamName, playIdSelected string) error,bool{
	vars := mux.Vars(request)
	guid := vars["guid"]
	inviteEmail := strings.ToLower(request.FormValue("inviteEmail"))
	possession := request.FormValue("possession")
	teamName := request.FormValue("teamName")
	playIdSelected := request.FormValue("playIdSelected")
	flippedPlay := request.FormValue("flippedPlay")
	//Add What is on the lines
	whatsOnTheLine:= request.FormValue("whatsOnTheLine")
	results, success, gameid, status := models.CreateGame(guid, inviteEmail, possession, teamName, playIdSelected, flippedPlay,whatsOnTheLine)
	output := jsonOutputs.GameOutput{Success: success, Message: results, Gameid: gameid, Statuscode: status}

	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
