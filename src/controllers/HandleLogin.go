package controllers

import (
	"encoding/json"
	"fmt"
	"jsonOutputs"
	"models"
	"net/http"
	"strings"
)

func HandleLogin(res http.ResponseWriter, request *http.Request) {
	username := strings.ToLower(request.FormValue("username"))
	password := request.FormValue("password")

	guid, success, message, status := models.AccountLogin(username, password)
	output := jsonOutputs.LoginOutput{success, message, guid, status}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
