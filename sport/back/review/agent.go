package review

import (
	"fmt"
	"os"
	"time"
)

func (ctx *Ctx) AgentDoOne() error {
	if err := ctx.DalDeleteExpiredReviewTokens(time.Duration(ctx.Config.ReviewExp)); err != nil {
		return err
	}
	return nil
}

func (ctx *Ctx) RunAgent() {
	for {
		err := ctx.AgentDoOne()
		if err != nil {
			fmt.Fprintf(os.Stderr, "review.AgentDoOne: %v\n", err)
			// TODO: notify admin
		}
		time.Sleep(time.Minute)
	}
}
