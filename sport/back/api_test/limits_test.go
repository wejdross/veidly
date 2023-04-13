package api_test

import (
	"sport/api"
	"sport/helpers"
	"strings"
	"testing"
)

func TestLimits(t *testing.T) {
	p := apiCtx.Request.BodyLimit
	defer func() {
		apiCtx.Request.BodyLimit = p
	}()
	apiCtx.Request.BodyLimit = struct {
		B int64
		K int64
		M int64
		G int64
	}{
		0, 1, 0, 0,
	}
	method := "POST"
	href := "/api/read_body"
	bodysize := apiCtx.Request.BodyLimit.B +
		apiCtx.Request.BodyLimit.K*api.SizeK +
		apiCtx.Request.BodyLimit.M*api.SizeM +
		apiCtx.Request.BodyLimit.G*api.SizeG
	var cases []api.TestCase
	/*
		normally this request would result in 401.
		we will later try to probe it with too large body and expect 400 instead
	*/
	cases = append(cases, api.TestCase{
		RequestMethod:      method,
		RequestUrl:         href,
		RequestReader:      strings.NewReader(helpers.CRNG_stringPanic(bodysize)),
		ExpectedStatusCode: 204,
	})
	cases = append(cases, api.TestCase{
		RequestMethod:      method,
		RequestUrl:         href,
		RequestReader:      strings.NewReader(helpers.CRNG_stringPanic(bodysize + 1)),
		ExpectedStatusCode: 413,
	})
	apiCtx.TestAssertCases(t, cases)
}
