package rsv_qr_test

// import (
// 	"encoding/json"
// 	"fmt"
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

// func TestAnonQr(t *testing.T) {
// 	token, _, err := qrCtx.Rsv.Instr.CreateTestInstructorWithUser()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	start := rsv.NowMin().Add(time.Hour)
// 	tr := train.NewTestCreateTrainingRequest(start, start.Add(time.Hour))
// 	tid, err := trainCtx.ApiCreateTraining(token, &tr)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := rsv.NewTestCreateReservationRequest(tid, start)
// 	r, err := qrCtx.Rsv.ApiCreateAndQueryRsv(nil, &rr, 0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if err = qrCtx.Rsv.AdyenSm.MoveStateToCapture(
// 		qrCtx.Rsv.RsvResponseToSmPassPtr(r),
// 		adyen_sm.ManualSrc, nil,
// 	); err != nil {
// 		t.Fatal(err)
// 	}

// 	// try to create qr
// 	m := "POST"
// 	u := "/api/qr/rsv"

// 	var qrID uuid.UUID

// 	qrCtx.Api.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod: m,
// 			RequestUrl:    u,
// 			RequestReader: helpers.JsonMustSerializeReader(rsv_qr.CreateQrCodeRequest{
// 				QrCodeRequest: rsv_qr.QrCodeRequest{
// 					RsvID: r.ID,
// 				},
// 				AccessToken: &r.AccessToken,
// 			}),
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				qrID, err = qrCtx.TestVerifyQr(b, r.ID)
// 				return err
// 			},
// 			ExpectedStatusCode: 200,
// 		},
// 	})

// 	// verify qr

// 	m = "GET"
// 	u = "/api/qr/rsv/eval"

// 	qrCtx.Api.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod:      "GET",
// 			RequestUrl:         u + "?id=" + qrID.String(),
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
// 				return nil
// 			},
// 		},
// 	})

// }
