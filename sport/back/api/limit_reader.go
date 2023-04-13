package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	reader with maximum size limit,
	ready for use as gin middleware
*/

const (
	SizeK  int64 = 1024
	SizeM  int64 = 1024 * SizeK
	SizeG  int64 = 1024 * SizeM
	SIZE_T int64 = 1024 * SizeG
)

type LimitReader struct {
	ctx   *gin.Context
	rdr   io.ReadCloser
	limit int64
	rd    int64
}

func (mbr *LimitReader) TooLargeResult() (n int, err error) {
	n, err = 0, fmt.Errorf("Request exceeded allowed size")

	ctx := mbr.ctx
	_ = ctx.Error(err)
	ctx.Header("connection", "close")
	ctx.AbortWithStatus(http.StatusRequestEntityTooLarge)

	return
}

func (mbr *LimitReader) PeekEOF(buf []byte) (n int, err error) {
	buf = buf[:1]
	n, err = mbr.rdr.Read(buf)
	if err == io.EOF && n == 0 {
		return
	}
	return mbr.TooLargeResult()
}

func (mbr *LimitReader) Read(buf []byte) (n int, err error) {

	if len(buf) == 0 {
		return 0, nil
	}

	if mbr.limit == mbr.rd {
		return mbr.PeekEOF(buf)
	}

	if int64(len(buf)) > (mbr.limit - mbr.rd) {
		buf = buf[:mbr.limit-mbr.rd]
	}

	n, err = mbr.rdr.Read(buf)

	mbr.rd += int64(n)

	if mbr.rd == mbr.limit {
		return mbr.PeekEOF(buf)
	}

	return
}

func (mbr *LimitReader) Close() error {
	return mbr.rdr.Close()
}

func LimitReaderMiddleware(_api *Ctx) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.Body = &LimitReader{
			ctx: ctx,
			rdr: ctx.Request.Body,
			limit: _api.Request.BodyLimit.B +
				_api.Request.BodyLimit.K*SizeK +
				_api.Request.BodyLimit.M*SizeM +
				_api.Request.BodyLimit.G*SizeG,
			rd: 0,
		}
		ctx.Next()
	}
}
