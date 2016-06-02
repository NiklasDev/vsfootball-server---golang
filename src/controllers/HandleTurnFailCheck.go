package controllers

import (
	// "encoding/json"
	"fmt"
	// "github.com/gorilla/mux"
	// "jsonOutputs"
	"models"
	"net/http"
)

func HandleTurnFailCheck(res http.ResponseWriter, request *http.Request) {
	// vars := mux.Vars(request)

	models.TurnFailCheck()

	// outputJson, _ := json.Marshal(jsonOutputs.GameIds{Success: success, Message: message})
	fmt.Fprintf(res, "Turns checked.")

}
