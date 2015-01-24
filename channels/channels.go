package channels

import (
	"log"
	"strconv"
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"net/http"
)

type Channel struct {
	Location appengine.GeoPoint
	SubscriptionID int
}

func OpenNewChannel(request *http.Request) (string, string, error) {
	newChannel := ParseRequestToChannel(request)
	context := appengine.NewContext(request)
	tempKey := datastore.NewIncompleteKey(context, "Channel", nil)
	savedKey, err := datastore.Put(context, tempKey, newChannel)
	if err != nil {
		return "", "", err
	}
	savedKeyString := strconv.FormatInt(savedKey.IntID(), 10)
	log.Print(savedKeyString)
	token, err := channel.Create(context, savedKeyString)
	StartChannelRefreshTask(savedKeyString)
	return token, savedKeyString, err
}

func SendToChannel(channelIdentifier string, request *http.Request) (error) {
	context := appengine.NewContext(request)
	err := channel.SendJSON(context, channelIdentifier, []string{"Stuff", "Things"})
	return err
}

func ParseRequestToChannel(request *http.Request) (*Channel) {
	lat, _ := strconv.ParseFloat(request.FormValue("lat"), 32)
	lng, _ := strconv.ParseFloat(request.FormValue("lng"), 32)

	location := appengine.GeoPoint{
		Lat: lat,
		Lng: lng }
	channel := &Channel {
		Location: location }
	return channel
}

func AddSubscriptionIDToChannel(channelIdentifier string, request *http.Request) {
	context := appengine.NewContext(request)
	channelId, err := strconv.ParseInt(channelIdentifier, 10, 32)
	if err != nil {
		return
	}
	existingKey := datastore.NewKey(context, "Channel", "", channelId, nil)
	var channel Channel
	err = datastore.Get(context, existingKey, &channel)
	if err != nil {
		return
	}
}

func LoadChannelFromSubscriptionID(subscriptionId int, request *http.Request) (Channel, error) {
	context := appengine.NewContext(request)
	query := datastore.NewQuery("Channel").Filter("SubscriptionID = ", subscriptionId)
	var channels []Channel
	_, err := query.GetAll(context, &channels)
	return channels[0], err
}

func StartChannelRefreshTask(channelId string) {

}
