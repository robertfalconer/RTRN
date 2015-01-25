package web

import (
	"channels"
	"html/template"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/results", resultsHandler)
	http.HandleFunc("/test-channel", testChannelHandler)
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	params := map[string]string{"": ""}

	parseTemplate("search", params, writer, request)
}

func resultsHandler(writer http.ResponseWriter, request *http.Request) {
	if request.FormValue("lat") == "" || request.FormValue("lng") == "" {
		rootHandler(writer, request)
		return
	}
	token, channelIdentifer, err := channels.OpenNewChannel(request)
	log.Print(token)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	params := map[string]string{"token": token, "channelId": channelIdentifer}
	parseTemplate("results", params, writer, request)
}

func testChannelHandler(writer http.ResponseWriter, request *http.Request) {
	channelIdentifier := request.FormValue("cid")
	channels.SendToChannel(channelIdentifier, request)
}

func renderTemplate(templateName string, params map[string]string, writer http.ResponseWriter, request *http.Request) {
	htmlTemplate := loadTemplate(templateName)
	err := htmlTemplate.Execute(writer, params)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadTemplate(templateName string) *template.Template {
	filename := "web/templates/" + templateName + ".html"
	return template.Must(template.ParseFiles(filename))
}
