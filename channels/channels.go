package channels

import (
	"appengine"
	"appengine/channel"
	"net/http"
)

func OpenNewChannel(request *http.Request) (string, error) {
	context := appengine.NewContext(request)
	token, err := channel.Create(context, "Test")
	return token, err
}

func SendToChannel(request *http.Request) (error) {
	context := appengine.NewContext(request)
	err := channel.SendJSON(context, "Test", []string{"Stuff", "Things"})
	return err
}
