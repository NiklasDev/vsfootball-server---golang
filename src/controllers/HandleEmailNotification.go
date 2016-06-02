package controllers

import (
	"encoding/json"
	"fmt"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleEmailNotification(res http.ResponseWriter, request *http.Request) {
	email := request.FormValue("email")

	success, message, status := models.SendEmailVerification(email)
	output := jsonOutputs.SuccessMessage{success, message, status}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
