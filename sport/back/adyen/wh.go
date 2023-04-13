package adyen

import (
	"fmt"
	"sport/helpers"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctx *Ctx) HandlerAdyenWH() gin.HandlerFunc {
	return func(g *gin.Context) {

		// raw as fuck, but will do for the moment
		token := g.Query("apikey")
		if token != ctx.Config.Auth.Apikey {
			g.AbortWithError(401, fmt.Errorf("invalid apikey provided"))
			return
		}

		var ni NotificationRequestItem
		if err := helpers.ReadJsonBodyFromReader(
			g.Request.Body, &ni, nil,
		); err != nil {
			go ctx.MailSupportAboutWHfail("JSON", err, nil)
			g.AbortWithError(400, err)
			return
		}

		id := ni.SessionID
		parts := strings.Split(id, "_")
		if len(parts) != 2 {
			e := fmt.Errorf("Invalid MerchantReference: " + id)
			go ctx.MailSupportAboutWHfail("MerchantReference", e, ni)
			g.AbortWithError(422, e)
			return
		}
		key := parts[0]
		val := parts[1]

		if key == "DONATE" {
			if ctx.Dal != nil {
				// raw method as fuck, but it's late
				someUuid := uuid.NewString()
				if len(val) < len(someUuid) {
					err := fmt.Errorf("invalid merchant id")
					go ctx.MailSupportAboutWHfail("InvalidMerchantID",
						err, ni)
					g.AbortWithError(400, err)
					return
				}
				value := ni.Amount
				if err := ctx.LogDonation(val[len(someUuid):], value); err != nil {
					go ctx.MailSupportAboutWHfail("Failed to log donation",
						err, ni)
					g.AbortWithError(500, err)
					return
				}
				if _, err := ctx.VerifyTransaction(&VerifyTransactionRequest{
					SessionID: ni.SessionID,
					Amount:    value,
					Currency:  ni.Currency,
					OrderID:   ni.OrderID,
				}); err != nil {
					go ctx.MailSupportAboutWHfail("Failed to verify transaction",
						err, ni)
					g.AbortWithError(500, err)
					return
				}
			}
			g.Status(204)
			g.Writer.WriteHeaderNow()
			g.Abort()
			return
		}

		if _, err := uuid.Parse(val); err != nil {
			go ctx.MailSupportAboutWHfail("InvalidMerchantReferenceUUID", err, ni)
			g.AbortWithError(422, err)
			return
		}

		if f, ok := ctx.handlers[key]; ok {
			ni.SessionID = parts[1]
			if e := f(&ni); e != nil {
				go ctx.MailSupportAboutWHfail("ABRT", e, ni)
				g.AbortWithError(422, e)
				return
			}
		} else {
			e := fmt.Errorf("unknown Handler key: %s", key)
			go ctx.MailSupportAboutWHfail("MerchantReference", e, ni)
			g.AbortWithError(422, e)
			return
		}
		g.Status(204)
		g.Writer.WriteHeaderNow()
		g.Abort()
	}
}
