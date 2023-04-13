package user

import (
	"fmt"
	"os"
	"time"
)

func (ctx *Ctx) RunRetentionLoop() {
	if !ctx.Config.EnableRetentionLoop {
		return
	}
	const maxErrCount = 5
	errCount := 0
	for {
		if errCount > maxErrCount {
			fmt.Println("RunRetentionLoop: terminating due to errors")
			return
		}
		err := ctx.DalDeleteInactiveUsers(
			time.Second *
				time.Duration(ctx.Config.RetentionLoopRecordTTLseconds))
		if err != nil {
			errCount++
			fmt.Fprintf(
				os.Stderr,
				"RunRetentionLoop. error %d/%d: %s\n",
				errCount,
				maxErrCount,
				err.Error())
		} else {
			errCount = 0
		}
		time.Sleep(time.Second *
			time.Duration(ctx.Config.RetentionLoopIntervalSeconds))
	}
}
