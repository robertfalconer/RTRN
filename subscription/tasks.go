package subscription

import (
	"appengine"
	"appengine/delay"
	"appengine/urlfetch"
	"channels"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const instagramSubscriptionsURL string = "https://api.instagram.com/v1/subscriptions/"
const instagramRecentMediaURL string = "https://api.instagram.com/v1/geographies/%s/media/recent?client_id=%s"
const instagramClientId string = "f0a3b8daa2944138816c1ed7cd91f666"
const instagramClientSecret string = "1e7869e58ae6463fbb6468eb6b9a7490"

type InstagramSubscriptionConfirmationMessage struct {
	Data InstagramSubscriptionConfirmation `json:"data"`
}

type InstagramSubscriptionConfirmation struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Object      string `json:"object"`
	ObjectId    string `json:"object_id"`
	Aspect      string `json:"aspect"`
	CallbackUrl string `json:"callback_url"`
}

type InstagramRecentMediaMessage struct {
	Data []InstagramMediaObject `json:"data"`
}

type InstagramMediaObject struct {
	Type   string                    `json:"type"`
	Images map[string]InstagramImage `json:"images"`
}

type InstagramImage struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

var Subscribe = delay.Func("subscribe", func(context appengine.Context, hostname string, channelId string, lat string, lng string) {
	params := url.Values{
		"client_id":     {instagramClientId},
		"client_secret": {instagramClientSecret},
		"object":        {"geography"},
		"aspect":        {"media"},
		"lat":           {lat},
		"lng":           {lng},
		"radius":        {"5000"},
		"callback_url":  {fmt.Sprintf("http://%s/webhook", hostname)},
	}

	client := urlfetch.Client(context)
	resp, requestErr := client.PostForm(instagramSubscriptionsURL, params)

	if requestErr != nil {
		log.Println("subscription setup failed with error", requestErr)
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("subscription request returned non-200 response, %d - %s", resp.StatusCode, body)
		return
	}

	var message InstagramSubscriptionConfirmationMessage
	responseErr := json.Unmarshal(body, &message)

	if responseErr != nil {
		log.Println("json unmarshalling failed with error", responseErr)
		log.Println(string(body))
		return
	}

	confirmation := message.Data
	CreateSubscriptionFromConfirmation(context, &confirmation, channelId)

	// NOTE: maybe send "subscription_created" message to client channel?
})

var FetchRecentMedia = delay.Func("fetch-recent-media", func(context appengine.Context, subscriptionId string, geographyId string) {
	log.Printf("fetching recent media for geography=%s and subscription=%s", geographyId, subscriptionId)

	client := urlfetch.Client(context)
	resp, requestErr := client.Get(fmt.Sprintf(instagramRecentMediaURL, geographyId, instagramClientId))

	if requestErr != nil {
		log.Println("fetch recent media request failed with error", requestErr)
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("fetch recent media request returned non-200 response, %d - %s", resp.StatusCode, body)
		return
	}

	var message InstagramRecentMediaMessage
	responseErr := json.Unmarshal(body, &message)

	if responseErr != nil {
		log.Println("json unmarshalling failed with error", responseErr)
		return
	}

	subscription, err := GetSubscriptionById(context, subscriptionId)

	if err != nil {
		// TODO unsubscribe this geography from instagram API
		log.Println("unable to find subscription for media update")
		return
	}

	for _, mediaObject := range message.Data {
		if mediaObject.Type == "video" {
			continue
		}
		mediaObjectAsJSON, _ := json.Marshal(&mediaObject)
		log.Println(string(mediaObjectAsJSON))
		channels.SendToChannel(context, subscription.ChannelId, string(mediaObjectAsJSON))
	}

	// log.Println(string(body))
})
