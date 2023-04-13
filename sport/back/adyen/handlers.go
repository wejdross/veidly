package adyen

import (
	"fmt"
	"net/mail"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctx *Ctx) DonateHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		amountStr := g.Query("amount")
		email := g.Query("email")
		_, err := mail.ParseAddress(email)
		if err != nil {
			g.AbortWithError(400, fmt.Errorf("invalid email: %v", err))
			return
		}
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			g.AbortWithError(409, err)
			return
		}
		if amount < 500 {
			g.AbortWithError(400, fmt.Errorf("minimum amount to donate is 500"))
			return
		}
		req := RegisterTransactionRequest{
			SessionID:   "DONATE_" + uuid.NewString() + email,
			Amount:      amount,
			Currency:    "PLN",
			Description: "donacja dla Veidly",
			Email:       email,
			Country:     "Pl",
			Language:    "pl",
			UrlReturn:   ctx.Config.DonateReturnUrl,
		}
		res, err := ctx.RegisterTransaction(&req)
		if err != nil {
			g.AbortWithError(400, fmt.Errorf("failed to register transaction: %v", err))
			return
		}
		g.AbortWithStatusJSON(200, struct {
			Url string
		}{
			res.Url,
		})
	}
}
