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

func HandleFacebookSignup(res http.ResponseWriter, request *http.Request) {
	email := strings.ToLower(request.FormValue("email"))
	accounttype := request.FormValue("accounttype")
	accesstoken := request.FormValue("accesstoken")
	tokenexpiration := request.FormValue("tokenexpiration")
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
		Password:        "",
		Accesstoken:     accesstoken,
		Accounttype:     accounttype,
		Tokenexpiration: tokenexpiration,
		Username:        email,
		Verified:        1,
		Playsowned:      "",
		Platform:        platform}
	success, message, guidReturn, status := models.FacebookSignup(user)
	output := jsonOutputs.LoginOutput{success, message, guidReturn, status}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}
