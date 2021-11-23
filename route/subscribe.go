package route

import (
	"net/http"
	"pancakewatch/pancakeswap"
	"pancakewatch/pcwdb"

	"github.com/gin-gonic/gin"
)

type SubscribeForm struct {
	TokenAddress string  `form:"token-address"`
	PhoneNumber  int     `form:"phone-number"`
	TargetPrice  float32 `form:"target-price"`
}

func Subscribe(db pcwdb.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var subForm SubscribeForm
		if err := c.ShouldBind(&subForm); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		token, err := pancakeswap.TokenInfo(subForm.TokenAddress)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		if token.Data.Name == "" {
			c.String(http.StatusBadRequest, "invalid token address")
			return
		}
		sub := pcwdb.Subscription{PhoneNumber: subForm.PhoneNumber, TargetPrice: subForm.TargetPrice}
		if err := db.Subscribe(subForm.TokenAddress, sub); err != nil {
			c.Error(err)
			return
		}
		msg := "Subscribed. Subscription will be removed after text is sent."
		c.String(http.StatusAccepted, msg)
	}
}
