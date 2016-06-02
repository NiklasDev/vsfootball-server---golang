package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"jsonOutputs"
	"models"
	"net/http"
	"templateUtil"
)

var store = sessions.NewCookieStore([]byte("Tempsalt"))

func HandlePush(res http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	platform := vars["type"]
	message := request.FormValue("message")
	deviceid := request.FormValue("deviceid")
	if platform == "ios" {
		go models.SendAPNSMessageInner(deviceid, message, "", "")
	} else if platform == "android" {
		go models.SendGCMMessage(deviceid, message, "")
	}
	output := jsonOutputs.SuccessMessage{
		Success:    "true",
		Message:    "Message sent.",
		Statuscode: 200}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))

}
func HandleGetAllDevices(res http.ResponseWriter, request *http.Request) {
	iosDevices := models.GetIosDevices()
	androiddevices := models.GetAndroidDevices()
	output := jsonOutputs.DeviceOutput{
		Success:        "True",
		Message:        "Gathered ios Devices",
		Status:         200,
		Iosdevices:     iosDevices,
		Androiddevices: androiddevices}
	outputJson, _ := json.Marshal(output)
	fmt.Fprintf(res, string(outputJson))
}

func HandlePushHome(res http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		templateUtil.Display(res, "login", nil)
	} else {
		username := request.FormValue("username")
		password := request.FormValue("password")
		context := make(map[string]interface{})
		if username != "admin" || password != "engagepush" {
			context["Messages"] = []string{"Authentication Failed."}
			templateUtil.Display(res, "login", context)
		} else {
			session, _ := store.Get(request, "session-name")
			session.Values["login"] = "true"
			session.Save(request, res)
			templateUtil.Display(res, "push", nil)
		}

	}

}
