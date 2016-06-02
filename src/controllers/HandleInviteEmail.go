package controllers

import (
	"net/http"
	// "models"
	// "github.com/gorilla/mux"
	"encoding/json"
	"fmt"
	"jsonOutputs"
)

func HandleInviteEmail(res http.ResponseWriter, request *http.Request) {
	// vars := mux.Vars(request)
	// guid := vars["guid"]

	// email := request.FormValue("EmailAddress")

	var output jsonOutputs.SuccessMessage

	output = jsonOutputs.SuccessMessage{Success: "true", Message: "Success", Statuscode: 200}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))

}
