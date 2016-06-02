package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"models"
	"net/http"
)

func HandleVerify(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	guid := vars["guid"]
	_, url := models.VerifyAccount(guid)
	fmt.Println(url)
	http.Redirect(res, request, url, http.StatusFound)
	// fmt.Fprintf(res, "Account has been verified.")
}
