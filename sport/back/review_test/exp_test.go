package review_test

import (
	"encoding/json"
	"fmt"
	"sport/adyen_sm"
	"sport/api"
	"sport/review"
	"sport/rsv"
	"sport/train"
	"sport/user"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestExpiry(t *testing.T) {
	token, _, err := instrCtx.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}
	start := rsv.NowMin().Add(time.Hour)
	tr := train.NewTestCreateTrainingRequest(start, start.Add(time.Hour))
	trainingID, err := rsvCtx.Train.ApiCreateTraining(token, &tr)
	if err != nil {
		t.Fatal(err)
	}

	clientToken, err := instrCtx.User.ApiCreateAndLoginUser(nil)
	if err != nil {
		t.Fatal(err)
	}
	clientToken2, err := instrCtx.User.ApiCreateAndLoginUser(nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := rsv.NewTestCreateReservationRequest(trainingID, start)
	r, err := rsvCtx.ApiCreateAndQueryRsv(&clientToken, &rr, 0)
	if err != nil {
		t.Fatal(err)
	}

	if err = rsvCtx.AdyenSm.MoveStateToPayout(
		rsvCtx.RsvResponseToSmPassPtr(r),
		adyen_sm.ManualSrc,
	); err != nil {
		t.Fatal(err)
	}

	if r, err = rsvCtx.RefreshRsvAndValidate(r.ID, rsv.RsvValidationOpts{
		State: adyen_sm.Payout,
	}); err != nil {
		t.Fatal(err)
	}

	method := "GET"
	path := "/api/review/user"

	accessToken := ""

	reviewCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         path,
			ExpectedStatusCode: 401,
		},
		{
			RequestMethod:      method,
			RequestUrl:         path,
			RequestHeaders:     user.GetAuthHeader(clientToken),
			ExpectedStatusCode: 400,
		},
		{
			RequestMethod:      method,
			RequestUrl:         path + "?rsv_id=" + uuid.NewString(),
			RequestHeaders:     user.GetAuthHeader(clientToken),
			ExpectedStatusCode: 404,
		},
		{
			RequestMethod:      method,
			RequestUrl:         path + "?rsv_id=" + r.ID.String(),
			RequestHeaders:     user.GetAuthHeader(clientToken2),
			ExpectedStatusCode: 404,
		},
		{
			RequestMethod:      method,
			RequestUrl:         path + "?rsv_id=" + r.ID.String(),
			RequestHeaders:     user.GetAuthHeader(clientToken),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				rres := review.ReviewResponse{}
				if err = json.Unmarshal(b, &rres); err != nil {
					return err
				}

				if rres.Type != review.TokenReviewType {
					return fmt.Errorf("unexpected review type, expected token")
				}

				res := rres.Token

				errCorrection := time.Minute
				expected := time.Now().In(time.UTC).Add(time.Duration(reviewCtx.Config.ReviewExp))
				df := res.ExpireOn.Sub(expected)
				if df < 0 {
					df = -df
				}
				if df >= errCorrection {
					return fmt.Errorf("unexpected review exp: expected: %v, got %v", expected, res.ExpireOn)
				}

				if res.AccessToken == "" {
					return fmt.Errorf("invalid review access token")
				}

				accessToken = res.AccessToken

				return nil
			},
		},
	})

	if err := reviewCtx.AgentDoOne(); err != nil {
		t.Fatal(err)
	}

	reviewCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         path + "?rsv_id=" + r.ID.String(),
			RequestHeaders:     user.GetAuthHeader(clientToken),
			ExpectedStatusCode: 200,
		},
	})

	exptime := time.Now().In(time.UTC).Add(-time.Duration(reviewCtx.Config.ReviewExp) - time.Minute)
	if err = reviewCtx.DalTestResetCreatedOn(accessToken, exptime); err != nil {
		t.Fatal(err)
	}

	if err := reviewCtx.AgentDoOne(); err != nil {
		t.Fatal(err)
	}

	reviewCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         path + "?rsv_id=" + r.ID.String(),
			RequestHeaders:     user.GetAuthHeader(clientToken),
			ExpectedStatusCode: 404,
		},
	})
}
