package channels

import (
	"appengine"
	"appengine/channel"
	"code.google.com/p/go-uuid/uuid"
	"log"
	"net/http"
)

func OpenChannel(context appengine.Context) (string, string, error) {
	channelId := uuid.New()
	log.Printf("creating new channel with id %s", channelId)
	token, err := channel.Create(context, channelId)
	return token, channelId, err
}

func SendToChannel(context appengine.Context, channelId string, data string) error {
	err := channel.SendJSON(context, channelId, data)
	return err
}

func ChannelClosed(request *http.Request) {
	// context := appengine.NewContext(request)
	// token := request.FormValue("from")
	// send token to subscription package to unsubscribe
	//
}
