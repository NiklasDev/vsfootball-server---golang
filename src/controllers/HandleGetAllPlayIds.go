package controllers

import (
	"encoding/json"
	"fmt"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleGetAllPlayIds(res http.ResponseWriter, request *http.Request) {
	success, message, status, playIds := models.GetPlaybookIds()
	outputJson, _ := json.Marshal(jsonOutputs.PlayIdsOutput{Success: success, Message: message, Playids: playIds, Statuscode: status})
	fmt.Fprintf(res, string(outputJson))
}
