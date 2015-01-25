package channels

import (
	"appengine"
	"appengine/channel"
	"code.google.com/p/go-uuid/uuid"
	"log"
)

func OpenChannel(context appengine.Context) (string, string, error) {
	channelId := uuid.New()
	log.Printf("creating new channel with id %s", channelId)
	token, err := channel.Create(context, channelId)
	return token, channelId, err
}

func SendToChannel(context appengine.Context, channelId string) error {
	err := channel.SendJSON(context, channelId, []string{"Stuff", "Things"})
	return err
}
