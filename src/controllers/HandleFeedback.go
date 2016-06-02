package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
	"time"
)

func HandleFeedback(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]

	comment := request.FormValue("comment")
	gameid := request.FormValue("gameid")
	screen := request.FormValue("screen")
	os := request.FormValue("os")
	fmt.Println("Comment:" + comment)
	fmt.Println("gameid:" + gameid)
	fmt.Println("screen:" + screen)

	feedback := &models.Feedback{
		Created: time.Now().Unix(),
		Comment: comment,
		Userid:  guid,
		Gameid:  gameid,
		Screen:  screen,
		Os:      os}
	success, message, status := models.AddFeedback(feedback)
	output := jsonOutputs.SuccessMessage{success, message, status}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
