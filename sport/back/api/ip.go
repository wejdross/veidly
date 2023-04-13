package api

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
)

func GetIP(ctx *gin.Context, tryHeaders bool) (string, error) {

	if tryHeaders {
		raddr := ctx.GetHeader("X-Real-IP")
		if raddr != "" {
			return raddr, nil
		}

		raddr = ctx.GetHeader("X-Forwarded-For")
		if raddr != "" {
			return raddr, nil
		}
	}

	ip, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", fmt.Errorf("empty ip")
	}
	return ip, nil
}
