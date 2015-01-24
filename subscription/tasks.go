package subscription

import (
	"appengine"
	"appengine/delay"
	"appengine/urlfetch"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
)

const instagramSubscriptionsURL string = "https://api.instagram.com/v1/subscriptions/"
const instagramClientId string = "f0a3b8daa2944138816c1ed7cd91f666"
const instagramClientSecret string = "1e7869e58ae6463fbb6468eb6b9a7490"
const ngrokProxyURL string = "http://6173d145.ngrok.com"

var Subscribe = delay.Func("subscribe", func(c appengine.Context, channelId int, lat float64, lng float64) {
	params := url.Values{
		"client_id":     {instagramClientId},
		"client_secret": {instagramClientSecret},
		"object":        {"geography"},
		"aspect":        {"media"},
		"lat":           {fmt.Sprintf("%f", lat)},
		"lng":           {fmt.Sprintf("%f", lng)},
		"radius":        {"5000"},
		"callback_url":  {fmt.Sprintf("%s/webhook/", ngrokProxyURL)},
	}

	client := urlfetch.Client(c)
	resp, err := client.PostForm(instagramSubscriptionsURL, params)

	if err != nil {
		log.Println("subscription setup failed with error", err)
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("%s - %s", resp.Status, body)
})
