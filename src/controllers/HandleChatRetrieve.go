package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleChatRetrieve(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]
	gameid := vars["gameid"]
	success, message, chatList, status := models.RetrieveChat(gameid, guid)

	outputJson, _ := json.Marshal(jsonOutputs.ChatListOutput{Success: success, Message: message, GameId: gameid, Chatlist: chatList, StatusCode: status})
	fmt.Fprintf(res, string(outputJson))
}
