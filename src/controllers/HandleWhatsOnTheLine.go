package controllers

import (
	"encoding/json"
	"fmt"
	"models"
	// "github.com/gorilla/mux"
	"jsonOutputs"
	"net/http"
)

func HandleWhatsOnTheLinePreset(res http.ResponseWriter, request *http.Request) {
	// vars := mux.Vars(request)
	// guid := vars["guid"]
	response.Header().Set("Content-Type", "application/json")
	status,message,statusCode,preset:=models.GetPreset()

	output := jsonOutputs.WhatsOnTheLinePresetList{
		Success: status, Message: message, Statuscode: statusCode, PresetList: preset}
    
    outputJson, _ := json.Marshal(output)

    fmt.Fprint(response, string(outputJson))

}

func HandleWhatsOnTheLine(res http.ResponseWriter, request *http.Request) {
	// vars := mux.Vars(request)
	// guid := vars["guid"]
	gameId := vars["gameid"]
	whatsOnTheLine := request.FormValue("whatsontheline")
	output := jsonOutputs.WhatsOnTheLineForAGame{
		Success: "true", Message: "Successfully added what's on the line items.", Statuscode: 200, Confirmed: true, WhatsOnTheLine: "whatsOnTheLine", WhatsOnTheLineId: 1}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}

func HandleWhatsOnTheLineForAGame(res http.ResponseWriter, request *http.Request) {
	// vars := mux.Vars(request)
	// guid := vars["guid"]
	// gameId := vars["gameid"]
	output := jsonOutputs.WhatsOnTheLineForAGame{
		Success: "true", Message: "Successfully gathered what's on the line items.", Statuscode: 200, Confirmed: false, WhatsOnTheLine: "Winner gets 5 dollars.", WhatsOnTheLineId: 1}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}

func HandleWhatsOnTheLineConfirm(res http.ResponseWriter, request *http.Request) {
	// vars := mux.Vars(request)
	// guid := vars["guid"]
	// gameId := vars["gameid"]
	// whatsOnTheLineId := vars["whatsonthelineid"]
	output := jsonOutputs.WhatsOnTheLineConfirmation{
		Success: "true", Message: "Successfully confirmed what's on the line.", Statuscode: 200, WhatsOnTheLineId: 1}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
