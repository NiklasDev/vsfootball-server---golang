package controllers

import (
	"encoding/json"
	"fmt"
	"jsonOutputs"
	"models"
	"net/http"
)

func HandleViralCoefficient(res http.ResponseWriter, request *http.Request) {

	success, message, status := models.ViralCoefficient()
	output := jsonOutputs.SuccessMessage{success, message, status}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
