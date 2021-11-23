package main

import (
	"fmt"
	"os"
	"pancakewatch/pancakeswap"
	"pancakewatch/pcwdb"
	"strconv"
	"time"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func renderFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func check(client *twilio.RestClient, db pcwdb.DB) error {
	from := os.Getenv("PHONE_NUMBER")
	return <-db.ForEachAsync(func(address string, subs []pcwdb.Subscription) error {
		token, err := pancakeswap.TokenInfo(address)
		if err != nil {
			return err
		}
		price, err := strconv.ParseFloat(token.Data.Price, 32)
		if err != nil {
			return err
		}
		for _, sub := range subs {
			if (sub.IsTargetUnder && float32(price) <= sub.TargetPrice) || float32(price) >= sub.TargetPrice {
				db.Unsubscribe(address, sub.PhoneNumber)
				params := &openapi.CreateMessageParams{}
				params.SetTo(fmt.Sprint(sub.PhoneNumber))
				params.SetFrom(from)
				msg := fmt.Sprintf(
					"Target price of $%v on %v has been reached. Current price: %v",
					renderFloat(float64(sub.TargetPrice)),
					token.Data.Name,
					renderFloat(price),
				)
				params.SetBody(msg)

				_, err := client.ApiV2010.CreateMessage(params)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func runChecker(client *twilio.RestClient, db pcwdb.DB) {
	for {
		if err := check(client, db); err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Second)
	}
}
