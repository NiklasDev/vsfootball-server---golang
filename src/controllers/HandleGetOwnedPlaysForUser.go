package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
	"strconv"
	"strings"
)

func HandleGetOwnedPlays(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]

	success, message, status, playIdCsv := models.GetOwnedPlays(guid)

	playIdStrings := strings.Split(playIdCsv, ",")
	if len(playIdCsv) > 0 {
		var playIdInts []int64

		for _, playid := range playIdStrings {
			intVal, _ := strconv.ParseInt(playid, 10, 64)
			if intVal > 0 {
				playIdInts = append(playIdInts, intVal)
			}
		}
		outputJson, _ := json.Marshal(jsonOutputs.PlayIdsOutput{Success: success, Message: message, Playids: playIdInts, Statuscode: status})
		fmt.Fprintf(res, string(outputJson))
	} else {
		outputJson, _ := json.Marshal(jsonOutputs.PlayIdsOutput{Success: success, Message: message, Playids: nil, Statuscode: status})
		fmt.Fprintf(res, string(outputJson))
	}

}
