package subscription

import (
	"appengine"
	"appengine/delay"
	"log"
)

var Subscribe = delay.Func("subscribe", func(c appengine.Context, channelId int, coords string) {
	log.Printf("establishing subscription for channel=%d and coords=%s", channelId, coords)
})
