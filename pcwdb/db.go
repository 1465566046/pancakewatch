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

// ForEachAsync will iterate over each address in the database.
// Each interation will be run in its own Goroutine.
// It will call fn, passing it the address and subscriptions at the address.
func (db DB) ForEachAsync(fn func(string, []Subscription) error) chan error {
	items := db.internalDB.Items()
	errs := make(chan error)
	for {
		address, subsJSON, err := items.Next()
		if err == pogreb.ErrIterationDone {
			break
		}
		if err != nil {
			errs <- err
			return errs
		}
		go func() {
			var subs []Subscription
			if err := json.Unmarshal(subsJSON, &subs); err != nil {
				errs <- err
				return
			}
			if err := fn(string(address), subs); err != nil {
				errs <- err
				return
			}
		}()
	}
	return errs
}

// Close will close the database.
// This call should be deferred to after error checking New.
func (db DB) Close() {
	db.internalDB.Close()
}
