package web

import (
	"appengine"
	"channels"
	"html/template"
	"log"
	"net/http"
	"subscription"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/results", resultsHandler)
	http.HandleFunc("/test-channel", testChannelHandler)
	http.HandleFunc("/_ah/channel/disconnected/", disconnectChannelHandler)
	http.HandleFunc("/webhook", subscription.InstagramWebhookHandler)
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	params := map[string]string{"": ""}
	renderTemplate("search", params, writer, request)
}

func resultsHandler(writer http.ResponseWriter, request *http.Request) {

	lat, lng := request.FormValue("lat"), request.FormValue("lng")

	if lat == "" || lng == "" {
		rootHandler(writer, request)
		return
	}

	context := appengine.NewContext(request)
	token, channelId, err := channels.OpenChannel(context)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("channel created with token %s and id %s", token, channelId)
	subscription.Subscribe.Call(context, request.Host, channelId, lat, lng)
	params := map[string]string{"token": token, "channelId": channelId}
	renderTemplate("results", params, writer, request)
}

func testChannelHandler(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	channelId := request.FormValue("cid")
	responseMap := map[string]string{"URL":"http://lorempixel.com/800/600/"}
	channels.SendToChannel(context, channelId, responseMap)
}

func disconnectChannelHandler(writer http.ResponseWriter, request *http.Request) {
	channels.ChannelClosed(request)
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
