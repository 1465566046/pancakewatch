// Package pcwdb acts as a wrapper around pogreb.
// Its primary purpose is to allow multiple Subscriptions
// to be stored under a single key.
package pcwdb

import (
	"encoding/json"

	"github.com/akrylysov/pogreb"
)

// DB wraps pogreb, hides its internal functions,
// and adds several helper methods.
type DB struct {
	internalDB *pogreb.DB
}

// New creates a new DB at path.
func New(path string) (db DB, err error) {
	db.internalDB, err = pogreb.Open(path, nil)
	return
}

// Subscribe appends sub to address.
// There can be multiple subscriptions per address.
func (db DB) Subscribe(address string, sub Subscription) error {
	key := []byte(address)
	has, err := db.internalDB.Has(key)
	if err != nil {
		return err
	}
	var subs []Subscription
	if has {
		subs, err = db.Subscriptions(address)
		if err != nil {
			return err
		}
	}
	subs = append(subs, sub)
	subsJSON, err := json.Marshal(subs)
	if err != nil {
		return err
	}

	err = db.internalDB.Put(key, subsJSON)
	return err
}

// Subscriptions returns all of the subscriptions at address.
func (db DB) Subscriptions(address string) (subs []Subscription, err error) {
	subsJSON, err := db.internalDB.Get([]byte(address))
	if err != nil {
		return
	}
	err = json.Unmarshal(subsJSON, &subs)
	return
}

// Unsubscribe removes a given phone number from
func (db DB) Unsubscribe(address string, phone int) error {
	subs, err := db.Subscriptions(address)
	if err != nil {
		return err
	}
	for i, sub := range subs {
		if sub.PhoneNumber == phone {
			updatedSubs := make([]Subscription, 0)
			updatedSubs = append(updatedSubs, subs[:i]...)
			subs = append(updatedSubs, subs[i+1:]...)
			break
		}
	}
	subsJSON, err := json.Marshal(subs)
	if err != nil {
		return err
	}
	err = db.internalDB.Put([]byte(address), subsJSON)
	return err
}

// ForEach will iterate over each address in the database.
// It will call fn, passing it the address and subscriptions at the address.
func (db DB) ForEach(fn func(string, []Subscription) error) error {
	items := db.internalDB.Items()
	for {
		address, subsJSON, err := items.Next()
		if err == pogreb.ErrIterationDone {
			break
		}
		if err != nil {
			return err
		}
		var subs []Subscription
		if err := json.Unmarshal(subsJSON, &subs); err != nil {
			return err
		}
		if err := fn(string(address), subs); err != nil {
			return err
		}
	}
	return nil
}

// Close will close the database.
// This call should be deferred to after error checking New.
func (db DB) Close() {
	db.internalDB.Close()
}
