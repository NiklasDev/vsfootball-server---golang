package templateUtil

import (
	"html/template"
	"net/http"
)

//Compile templates on start
var Templates = template.Must(template.ParseFiles("tmpl/login.html", "tmpl/push.html"))

//Display the named template
func Display(w http.ResponseWriter, tmpl string, data interface{}) {
	Templates.ExecuteTemplate(w, tmpl+".html", data)
}
