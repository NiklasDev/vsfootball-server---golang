package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleStatsTotal(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]
	totalOrOppGuid := vars["opponent_guid"]
	success, message, status, wins, losses, ties := models.GetStats(guid, totalOrOppGuid)
	output := jsonOutputs.StatsOutput{
		Success:    success,
		Message:    message,
		Statuscode: status,
		Wins:       wins,
		Losses:     losses,
		Ties:       ties}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}

func HandleStatsForOpponent(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]
	success, message, statuscode, guidEmails := models.GetOpponentsForGuid(guid)
	var guidEmailJsonArray []jsonOutputs.PlayerGuidEmailOutput
	for _, guidEmail := range guidEmails {
		opp := jsonOutputs.PlayerGuidEmailOutput{
			Email: guidEmail["Email"],
			Guid:  guidEmail["Guid"]}
		guidEmailJsonArray = append(guidEmailJsonArray, opp)
	}
	output := jsonOutputs.OpponentsOutput{
		Success:    success,
		Message:    message,
		Statuscode: statuscode,
		Opponents:  guidEmailJsonArray}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
