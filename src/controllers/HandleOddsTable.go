package controllers


import (
	"net/http"
	"fmt"
	// "io"
	// "mime/multipart"
	// "io/ioutil"
)

func HandleOddsTableUpdate (res http.ResponseWriter,request *http.Request){
	_, _, errinput := request.FormFile("csv")
	fmt.Println(request.FormValue("email"))
	fmt.Println(errinput)
	// fmt.Println(handler.Filename)
	// data, _ := ioutil.ReadAll(file) 
	// fmt.Println(data)
	fmt.Println("IM HERE")
}