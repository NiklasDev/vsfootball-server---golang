package controllers

import (
	"encoding/json"
	"fmt"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleGetAnimationIds(res http.ResponseWriter, request *http.Request) {
	success, message, status, animationIds := models.GetAnimationIds()
	output := jsonOutputs.AnimationIdsOutput{Success: success, Message: message, Statuscode: status, Animationids: animationIds}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
