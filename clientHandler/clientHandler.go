package clientHandler

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"html/template"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/results", handler)
}

func handler(writer http.ResponseWriter, request *http.Request) {
	var templateName string
	if request.URL.Path == "/" {
		templateName = "search";
	} else {
		templateName = request.URL.Path[len("/"):];
	}
	htmlTemplate, err := loadTemplate(templateName, writer)
	if err != nil {
		fmt.Fprint(writer, "Here")
		http.Error(writer, err.Error(), http.StatusNotFound)
	}
	err = htmlTemplate.Execute(writer, "template")
	if err != nil {
		fmt.Fprint(writer, "There")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func loadTemplate(templateName string, writer http.ResponseWriter) (*template.Template, error) {
	filename := "clientHandler/templates/" + templateName + ".html"
	html, err := ioutil.ReadFile(filename)
	fmt.Fprint(writer, "::" + string(html) + "::")
	if err != nil {
		fmt.Fprint(writer, "Over Here")
		return nil, err
	}
	return template.Must(template.New("template").Parse(string(html))), nil
}
