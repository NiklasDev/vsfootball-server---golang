package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"jsonOutputs"
	"models"
	"net/http"
	"strings"
	"time"
)

func HandleSignUp(res http.ResponseWriter, request *http.Request) {
	emailAddress := strings.ToLower(request.FormValue("email"))
	password := request.FormValue("password")
	firstname := request.FormValue("firstname")
	lastname := request.FormValue("lastname")
	platform := request.FormValue("platform")
	guid, guidErr := uuid.NewV4()
	if guidErr != nil {
		fmt.Println(guidErr.Error())
	}
	user := &models.User{
		Created:         time.Now().Unix(),
		Updated:         time.Now().Unix(),
		Firstname:       firstname,
		Lastname:        lastname,
		Guid:            strings.Replace(guid.String(), "-", "", -1),
		Password:        password,
		Accesstoken:     "",
		Accounttype:     "",
		Tokenexpiration: "",
		Username:        emailAddress,
		Verified:        0,
		Playsowned:      "",
		Platform:        platform}
	success, message, status := models.CreateAccount(user)
	output := jsonOutputs.SuccessMessage{success, message, status}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
