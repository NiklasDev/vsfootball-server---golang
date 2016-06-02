package controllers

import (
	"encoding/json"
	"fmt"
	"jsonOutputs"
	"models"
	"net/http"
	"strings"
)

func HandleForgotPassword(res http.ResponseWriter, request *http.Request) {
	email := strings.ToLower(request.FormValue("email"))
	success, message, status := models.ForgotPassword(email)
	output := jsonOutputs.SuccessMessage{success, message, status}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
