package review_test

import (
	"encoding/json"
	"fmt"
	"sport/adyen_sm"
	"sport/api"
	"sport/helpers"
	"sport/review"
	"sport/rsv"
	"sport/train"
	"sport/user"
	"testing"
	"time"
)

func TestContent(t *testing.T) {
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

	accessToken := ""

	reviewCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "GET",
			RequestUrl:         "/api/review/user?rsv_id=" + r.ID.String(),
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

	method := "POST"
	path := "/api/review"

	exptime := time.Now().In(time.UTC).Add(-time.Duration(reviewCtx.Config.ReviewExp) - time.Minute)
	if err = reviewCtx.DalTestResetCreatedOn(accessToken, exptime); err != nil {
		t.Fatal(err)
	}

	mark := 4
	rc := helpers.CRNG_stringPanic(10)

	reviewCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         path,
			ExpectedStatusCode: 400,
		},
		{
			RequestMethod:      method,
			RequestUrl:         path,
			RequestReader:      helpers.JsonMustSerializeReader(review.UpdateReviewRequest{}),
			ExpectedStatusCode: 400,
		},
		{
			RequestMethod: method,
			RequestUrl:    path,
			RequestReader: helpers.JsonMustSerializeReader(review.UpdateReviewRequest{
				ReviewContent: review.ReviewContent{
					Mark: 4,
				},
			}),
			ExpectedStatusCode: 400,
		},
		{
			RequestMethod: method,
			RequestUrl:    path,
			RequestReader: helpers.JsonMustSerializeReader(review.UpdateReviewRequest{
				ReviewContent: review.ReviewContent{
					Mark: review.MaxMark + 1,
				},
				AccessToken: accessToken,
			}),
			ExpectedStatusCode: 400,
		},
		{
			RequestMethod: method,
			RequestUrl:    path,
			RequestReader: helpers.JsonMustSerializeReader(review.UpdateReviewRequest{
				ReviewContent: review.ReviewContent{
					Mark:   mark,
					Review: rc,
				},
				AccessToken: accessToken,
			}),
			ExpectedStatusCode: 204,
		},
		{
			RequestMethod:      "GET",
			RequestUrl:         "/api/review/user?rsv_id=" + r.ID.String(),
			RequestHeaders:     user.GetAuthHeader(clientToken),
			ExpectedStatusCode: 200,
		},
	})

	// determine that agent cant affect this review

	if err := reviewCtx.AgentDoOne(); err != nil {
		t.Fatal(err)
	}

	// if rv, err := reviewCtx.DalReadSingleReview(
	// 	review.SingleRevKeyUserRsvKey{
	// 		UserID: clientID,
	// 		RsvID:  r.ID,
	// 	}, review.SingleRevKeyUserRsv, review.ContentReviewType, nil,
	// ); err != nil {
	// 	t.Fatal(err)
	// } else {
	// 	if rv.TrainingID != *r.TrainingID {
	// 		t.Fatal("invalid review training id")
	// 	}
	// 	if rv.AccessToken != "" {
	// 		t.Fatal("invalid review access token")
	// 	}
	// 	if rv.CreatedOn.Sub(exptime) >= time.Minute {
	// 		t.Fatal("invalid review created on")
	// 	}
	// 	if rv.Mark != mark {
	// 		t.Fatal("invalid review mark")
	// 	}
	// 	if rv.Review != rc {
	// 		t.Fatal("invalid review content")
	// 	}
	// }

	reviewCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "GET",
			RequestUrl:         "/api/review/user?rsv_id=" + r.ID.String(),
			RequestHeaders:     user.GetAuthHeader(clientToken),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				rres := review.ReviewResponse{}
				if err = json.Unmarshal(b, &rres); err != nil {
					return err
				}

				if rres.Type != review.ContentReviewType {
					return fmt.Errorf("unexpected review type, expected content")
				}

				res := rres.Content

				if rres.Token != nil {
					t.Fatal("invalid review token")
				}

				if res.Mark != mark {
					t.Fatal("invalid review mark")
				}
				if res.Review != rc {
					t.Fatal("invalid review content")
				}

				return nil
			},
		},
	})

	// public
	reviewCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      "GET",
			RequestUrl:         "/api/review/pub?training_id=" + r.TrainingID.String(),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var rres []review.PubReview
				if err = json.Unmarshal(b, &rres); err != nil {
					return err
				}

				if len(rres) != 1 {
					return fmt.Errorf("unexpected pub review length")
				}

				res := rres[0].ReviewContent

				if res.Mark != mark {
					t.Fatal("invalid pub review mark")
				}
				if res.Review != rc {
					t.Fatal("invalid pub review content")
				}

				return nil
			},
		},
	})

}
