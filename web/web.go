package web

import (
	"net/http"
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
	htmlTemplate := loadTemplate(templateName, writer)
	err := htmlTemplate.Execute(writer, nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadTemplate(templateName string, writer http.ResponseWriter) (*template.Template) {
	fileName := "web/templates/" + templateName + ".html"
	return template.Must(template.ParseFiles(fileName))
}
