package web

import (
	"net/http"
	"io/ioutil"
	"html/template"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/results", resultsHandler)
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	parseTemplate("search", writer, request)
}

func resultsHandler(writer http.ResponseWriter, request *http.Request) {
	parseTemplate("results", writer, request)
}

func parseTemplate(templateName string, writer http.ResponseWriter, request *http.Request) {
	htmlTemplate, err := loadTemplate(templateName, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	err = htmlTemplate.Execute(writer, "template")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadTemplate(templateName string, writer http.ResponseWriter) (*template.Template, error) {
	filename := "web/templates/" + templateName + ".html"
	html, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return template.Must(template.New("template").Parse(string(html))), nil
}
