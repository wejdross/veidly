package sub_qr_test

// import (
// 	"encoding/json"
// 	_ "image/png"
// 	"sport/adyen_sm"
// 	"sport/api"
// 	"sport/helpers"
// 	"sport/sub"
// 	"sport/sub_qr"
// 	"sport/user"
// 	"testing"

// 	"github.com/google/uuid"
// )

// func testQr() error {
// 	token, _, err := qrCtx.Sub.Instr.CreateTestInstructorWithUser()
// 	if err != nil {
// 		return err
// 	}

// 	utoken, err := uctx.ApiCreateAndLoginUser(nil)
// 	if err != nil {
// 		return err
// 	}

// 	smReq := sub.NewTestSubModelRequest()
// 	smReq.MaxEntrances = 10
// 	smID, err := subCtx.ApiCreateSubModel(token, &smReq)
// 	if err != nil {
// 		return err
// 	}

// 	s, err := subCtx.ApiCreateAndQuerySub(utoken, smID)
// 	if err != nil {
// 		return err
// 	}

// 	if err := subCtx.AdyenSm.MoveStateToCapture(
// 		subCtx.SubToSmPassPtr(s), adyen_sm.ManualSrc, nil,
// 	); err != nil {
// 		return err
// 	}

// 	// try to create qr
// 	m := "POST"
// 	u := "/api/qr/sub"

// 	var qrID uuid.UUID

// 	err = qrCtx.Api.TestAssertCasesErr([]api.TestCase{
// 		{
// 			RequestMethod:      m,
// 			RequestUrl:         u,
// 			ExpectedStatusCode: 401,
// 		}, {
// 			RequestMethod: m,
// 			RequestUrl:    u,
// 			RequestReader: helpers.JsonMustSerializeReader(sub_qr.CreateQrCodeRequest{
// 				QrCodeRequest: sub_qr.QrCodeRequest{
// 					SubID: s.ID,
// 				},
// 			}),
// 			ExpectedStatusCode: 401,
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         u,
// 			RequestHeaders:     user.GetAuthHeader(token),
// 			ExpectedStatusCode: 400,
// 		}, {
// 			RequestMethod:  m,
// 			RequestUrl:     u,
// 			RequestHeaders: user.GetAuthHeader(token),
// 			RequestReader: helpers.JsonMustSerializeReader(sub_qr.CreateQrCodeRequest{
// 				QrCodeRequest: sub_qr.QrCodeRequest{
// 					SubID: s.ID,
// 				},
// 			}),
// 			ExpectedStatusCode: 404,
// 		}, {
// 			RequestMethod:  m,
// 			RequestUrl:     u,
// 			RequestHeaders: user.GetAuthHeader(utoken),
// 			RequestReader: helpers.JsonMustSerializeReader(sub_qr.CreateQrCodeRequest{
// 				QrCodeRequest: sub_qr.QrCodeRequest{
// 					SubID: s.ID,
// 				},
// 			}),
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				qrID, err = qrCtx.TestVerifyQr(b, s.ID)
// 				return err
// 			},
// 			ExpectedStatusCode: 200,
// 		},
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	// verify qr

// 	m = "GET"
// 	u = "/api/qr/sub/eval"

// 	uq := u + "?id=" + qrID.String()

// 	err = qrCtx.Api.TestAssertCasesErr([]api.TestCase{
// 		{
// 			RequestMethod:      m,
// 			RequestUrl:         u,
// 			ExpectedStatusCode: 401,
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         uq,
// 			ExpectedStatusCode: 401,
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         uq,
// 			RequestHeaders:     user.GetAuthHeader(utoken),
// 			ExpectedStatusCode: 404,
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         uq,
// 			RequestHeaders:     user.GetAuthHeader(token),
// 			ExpectedStatusCode: 200,
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				var res sub.Sub
// 				if err := json.Unmarshal(b, &res); err != nil {
// 					return err
// 				}
// 				return helpers.AssertJsonErr("get sub", &res, &s)
// 			},
// 		},
// 	})

// 	// confirm qr

// 	u = "/api/qr/sub/confirm"
// 	uq = u + "?id=" + qrID.String()

// 	cs := []api.TestCase{
// 		{
// 			RequestMethod:      m,
// 			RequestUrl:         u,
// 			ExpectedStatusCode: 401,
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         uq,
// 			ExpectedStatusCode: 401,
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         uq,
// 			RequestHeaders:     user.GetAuthHeader(utoken),
// 			ExpectedStatusCode: 409,
// 		},
// 	}

// 	for i := 0; i < smReq.MaxEntrances; i++ {
// 		cs = append(cs, api.TestCase{
// 			RequestMethod:      m,
// 			RequestUrl:         uq,
// 			RequestHeaders:     user.GetAuthHeader(token),
// 			ExpectedStatusCode: 204,
// 		})
// 	}

// 	cs = append(cs, api.TestCase{
// 		RequestMethod:      m,
// 		RequestUrl:         uq,
// 		RequestHeaders:     user.GetAuthHeader(token),
// 		ExpectedStatusCode: 409,
// 	})

// 	return qrCtx.Api.TestAssertCasesErr(cs)
// }

// func TestQr(t *testing.T) {
// 	if err := testQr(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestLotsOfQr(t *testing.T) {
// 	cqty := 10000
// 	waiter := make(chan error, cqty)
// 	semaphore := make(chan struct{}, 48)
// 	for i := 0; i < cqty; i++ {
// 		semaphore <- struct{}{}
// 		go func(_i int) {
// 			waiter <- testQr()
// 			<-semaphore
// 		}(i)
// 	}
// 	var e error
// 	for i := 0; i < cqty; i++ {
// 		err := <-waiter
// 		if err != nil && e == nil {
// 			e = err
// 		}
// 	}

// 	if e != nil {
// 		t.Fatal(e)
// 	}
// }
