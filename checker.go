package main

import (
	"fmt"
	"pancakewatch/pcwdb"
)

func check(db pcwdb.DB) {
	if err := db.ForEach(func(address string, subs []pcwdb.Subscription) error {
		fmt.Println("address", address)
		for _, sub := range subs {
			fmt.Println("phone number", sub.PhoneNumber, "target price", sub.TargetPrice)
		}
		return nil
	}); err != nil {
		panic(err)
	}
}
