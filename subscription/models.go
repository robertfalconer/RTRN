package subscription

import (
	"appengine"
	"appengine/datastore"
)

type Subscription struct {
	Type      string
	Object    string
	ObjectId  string
	Aspect    string
	ChannelId string
}

func GetSubscriptionById(context appengine.Context, subscriptionId string) (*Subscription, error) {
	subscriptionKey := datastore.NewKey(context, "Subscription", subscriptionId, 0, nil)

	var subscription Subscription
	err := datastore.Get(context, subscriptionKey, &subscription)

	if err != nil {
		return nil, err
	}

	return &subscription, err
}

func CreateSubscriptionFromConfirmation(context appengine.Context, confirmation *InstagramSubscriptionConfirmation, channelId string) (*datastore.Key, error) {
	subscriptionKey := datastore.NewKey(context, "Subscription", confirmation.Id, 0, nil)
	subscription := &Subscription{
		Type:      confirmation.Type,
		Object:    confirmation.Object,
		ObjectId:  confirmation.ObjectId,
		Aspect:    confirmation.Aspect,
		ChannelId: channelId,
	}
	return datastore.Put(context, subscriptionKey, subscription)
}
