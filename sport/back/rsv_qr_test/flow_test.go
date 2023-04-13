package rsv_qr_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	_ "image/png"
// 	"sport/adyen_sm"
// 	"sport/api"
// 	"sport/helpers"
// 	"sport/rsv"
// 	"sport/rsv_qr"
// 	"sport/train"
// 	"sport/user"
// 	"testing"
// 	"time"

// 	"github.com/google/uuid"
// )

// func testQr() error {
// 	token, _, err := qrCtx.Rsv.Instr.CreateTestInstructorWithUser()
// 	if err != nil {
// 		return err
// 	}

// 	utoken, err := trainCtx.User.ApiCreateAndLoginUser(nil)
// 	if err != nil {
// 		return err
// 	}

// 	utoken2, err := trainCtx.User.ApiCreateAndLoginUser(nil)
// 	if err != nil {
// 		return err
// 	}

// 	start := rsv.NowMin().Add(time.Hour)
// 	tr := train.NewTestCreateTrainingRequest(start, start.Add(time.Hour))
// 	tid, err := trainCtx.ApiCreateTraining(token, &tr)
// 	if err != nil {
// 		return err
// 	}

// 	rr := rsv.NewTestCreateReservationRequest(tid, start)
// 	r, err := qrCtx.Rsv.ApiCreateAndQueryRsv(&utoken, &rr, 0)
// 	if err != nil {
// 		return err
// 	}

// 	if err = qrCtx.Rsv.AdyenSm.MoveStateToCapture(
// 		qrCtx.Rsv.RsvResponseToSmPassPtr(r),
// 		adyen_sm.ManualSrc, nil,
// 	); err != nil {
// 		return err
// 	}

// 	// try to create qr
// 	m := "POST"
// 	u := "/api/qr/rsv"

// 	var qrID uuid.UUID

// 	err = qrCtx.Api.TestAssertCasesErr([]api.TestCase{
// 		{
// 			RequestMethod:      m,
// 			RequestUrl:         u,
// 			ExpectedStatusCode: 400,
// 		}, {
// 			RequestMethod: m,
// 			RequestUrl:    u,
// 			RequestReader: helpers.JsonMustSerializeReader(rsv_qr.CreateQrCodeRequest{
// 				QrCodeRequest: rsv_qr.QrCodeRequest{
// 					RsvID: r.ID,
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
// 			RequestReader: helpers.JsonMustSerializeReader(rsv_qr.CreateQrCodeRequest{
// 				QrCodeRequest: rsv_qr.QrCodeRequest{
// 					RsvID: r.ID,
// 				},
// 			}),
// 			ExpectedStatusCode: 404,
// 		}, {
// 			RequestMethod:  m,
// 			RequestUrl:     u,
// 			RequestHeaders: user.GetAuthHeader(utoken2),
// 			RequestReader: helpers.JsonMustSerializeReader(rsv_qr.CreateQrCodeRequest{
// 				QrCodeRequest: rsv_qr.QrCodeRequest{
// 					RsvID: r.ID,
// 				},
// 			}),
// 			ExpectedStatusCode: 404,
// 		}, {
// 			RequestMethod:  m,
// 			RequestUrl:     u,
// 			RequestHeaders: user.GetAuthHeader(utoken),
// 			RequestReader: helpers.JsonMustSerializeReader(rsv_qr.CreateQrCodeRequest{
// 				QrCodeRequest: rsv_qr.QrCodeRequest{
// 					RsvID: r.ID,
// 				},
// 			}),
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				qrID, err = qrCtx.TestVerifyQr(b, r.ID)
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
// 	u = "/api/qr/rsv/eval"

// 	uq := u + "?id=" + qrID.String()

// 	err = qrCtx.Api.TestAssertCasesErr([]api.TestCase{
// 		{
// 			RequestMethod:      "GET",
// 			RequestUrl:         u,
// 			ExpectedStatusCode: 401,
// 		}, {
// 			RequestMethod:      "GET",
// 			RequestUrl:         uq,
// 			ExpectedStatusCode: 401,
// 		}, {
// 			RequestMethod:      "GET",
// 			RequestUrl:         uq,
// 			RequestHeaders:     user.GetAuthHeader(utoken),
// 			ExpectedStatusCode: 404,
// 		}, {
// 			RequestMethod:      "GET",
// 			RequestUrl:         uq,
// 			RequestHeaders:     user.GetAuthHeader(token),
// 			ExpectedStatusCode: 200,
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				var res rsv_qr.EvalQrResponse
// 				if err := json.Unmarshal(b, &res); err != nil {
// 					return err
// 				}
// 				if res.RsvID != r.ID {
// 					return fmt.Errorf("Invalid rsv id")
// 				}
// 				if res.ConfirmCode != 0 {
// 					return fmt.Errorf("Invalid confirm code, expected 0 got %d", res.ConfirmCode)
// 				}
// 				return nil
// 			},
// 		}, {
// 			RequestMethod:      "GET",
// 			RequestUrl:         uq,
// 			RequestHeaders:     user.GetAuthHeader(token),
// 			ExpectedStatusCode: 200,
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				var res rsv_qr.EvalQrResponse
// 				if err := json.Unmarshal(b, &res); err != nil {
// 					return err
// 				}
// 				if res.RsvID != r.ID {
// 					return fmt.Errorf("Invalid rsv id")
// 				}
// 				if res.ConfirmCode != rsv_qr.AlreadyConfirmed {
// 					return fmt.Errorf("Invalid confirm code")
// 				}
// 				return nil
// 			},
// 		},
// 	})
// 	return err
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
