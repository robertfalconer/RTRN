package web

import (
	"log"
	"net/http"
	"html/template"
	"channels"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/results", resultsHandler)
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	params := map[string]string{"":""}

	parseTemplate("search", params, writer, request)
}

func resultsHandler(writer http.ResponseWriter, request *http.Request) {
	token, err := channels.OpenNewChannel(request);
	log.Print(token)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	params := map[string]string{"token":token}
	parseTemplate("results", params, writer, request)
}

func parseTemplate(templateName string, params map[string]string, writer http.ResponseWriter, request *http.Request) {
	htmlTemplate, err := loadTemplate(templateName)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	err = htmlTemplate.Execute(writer, params)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadTemplate(templateName string) (*template.Template, error) {
	filename := "web/templates/" + templateName + ".html"
	return template.Must(template.ParseFiles(filename)), nil
}
