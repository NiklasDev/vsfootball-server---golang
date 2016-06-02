package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

//Put chat submit functionality here.
// Send a message
// Parms: guid. gameid, text
// url : /{GUID}/game/{GAMEID}/chat
// Request method :POST
// request body :
//  Message <STRING>
func HandleChatSubmit(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]
	gameid := vars["gameid"]
	message := request.FormValue("message")
	success, response, status := models.SubmitChat(guid, gameid, message)
	fmt.Println(response)
	output := jsonOutputs.SuccessMessage{success, response, status}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
