package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleGetAnimationForAnimationId(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	animationId := vars["animationid"]
	success, message, status, animation := models.GetAnimationForAnimationId(animationId)
	output := jsonOutputs.AnimationOutput{Success: success, Message: message, Statuscode: status, AnimationObj: animation}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
