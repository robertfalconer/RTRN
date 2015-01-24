package subscription

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type InstagramSubscriptionUpdate struct {
	SubscriptionId int    `json:"subscription_id"`
	Object         string `json:"object"`
	ObjectId       string `json:"object_id"`
	ChangedAspect  string `json:"changed_aspect"`
	Time           int    `json:"time"`
}

func InstagramWebhook(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		log.Println("received instagram subscription challenge")
		fmt.Fprint(w, r.FormValue("hub.challenge"))
		return
	}

	if r.Method == "POST" {
		log.Println("received instagram subscription update")

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		message := make([]InstagramSubscriptionUpdate, 0)
		err := json.Unmarshal(body, &message)

		if err != nil {
			log.Println("json unmarshalling failed with error", err)
			log.Println(string(body))
		} else {
			log.Println("message:", message)
		}
	}

}
