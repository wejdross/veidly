package dc_test

import (
	"sport/api"
	"sport/dc"
	"sport/helpers"
	"sport/rsv"
	"sport/train"
	"sport/user"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRsv(t *testing.T) {

	sf, pf, rv := 0, 0, 0
	sf = rsvCtx.Config.ServiceFee
	pf = rsvCtx.Config.ProcessingFee
	rv = rsvCtx.Config.RefundAmount
	rsvCtx.Config.ServiceFee = 10
	rsvCtx.Config.ProcessingFee = 3
	rsvCtx.Config.RefundAmount = 30
	defer func() {
		rsvCtx.Config.ServiceFee = sf
		rsvCtx.Config.ProcessingFee = pf
		rsvCtx.Config.RefundAmount = rv
	}()

	token, iid, err := instrCtx.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}
	// POST
	m := "POST"
	u := "/api/dc"

	s := time.Now().In(time.UTC).Add(-time.Hour).Round(time.Second)

	dcs := []dc.DcRequest{
		// if u change this disount code remember to update expected prices below
		{
			Name:       "foo",
			Quantity:   1,
			ValidStart: dc.Time(s),
			ValidEnd:   dc.Time(s.Add(2 * time.Hour)),
			Discount:   50,
		},
		{
			Name:       "bar",
			Quantity:   2,
			ValidStart: dc.Time(s),
			ValidEnd:   dc.Time(s.Add(2 * time.Hour)),
			Discount:   10,
		},
	}

	dcCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(dcs[0]),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(dcs[1]),
			RequestHeaders:     user.GetAuthHeader(token),
		},
	})

	m = "GET"
	ids := make(map[string]uuid.UUID)

	dcCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return AssertGetRes(b, dcs, iid, ids)
			},
		},
	})

	// if u change it change expected paymnents below
	const price = 1000

	m = "POST"
	u = "/api/dc/binding"
	s = rsv.NowMin().Add(time.Hour)
	tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Hour))
	tr.Training.AllowExpress = true
	tr.Training.Price = price
	tid, err := trainCtx.ApiCreateTraining(token, &tr)
	if err != nil {
		t.Fatal(err)
	}

	dcCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 204,
			RequestReader: helpers.JsonMustSerializeReader(dc.DcBindingRequest{
				TrainingID: tid,
				DcID:       ids[dcs[0].Name],
			}),
			RequestHeaders: user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 204,
			RequestReader: helpers.JsonMustSerializeReader(dc.DcBindingRequest{
				TrainingID: tid,
				DcID:       ids[dcs[1].Name],
			}),
			RequestHeaders: user.GetAuthHeader(token),
		},
	})

	targetDc := dc.Dc{
		ID:               ids[dcs[0].Name],
		RedeemedQuantity: 0,
		InstrID:          iid,
		DcRequest:        dcs[0],
	}
	targetDcID := ids[dcs[0].Name]

	rr := rsv.NewTestCreateReservationRequest(tid, s)
	rr.DcID = &targetDcID
	pi, err := rsvCtx.ApiGetRsvPricing(&rr, 0)
	if err != nil {
		t.Fatal(err)
	}

	expectedPi := rsv.RsvPricingInfo{
		TotalPrice:     515,
		ProcessingFee:  15,
		SplitPayout:    450,
		SplitIncomeFee: 50,
		RefundAmount:   150,
		InstrPrice:     500,
		Dc:             pi.Dc,
	}
	if err := helpers.AssertJsonErr("rsv DC", &pi, &expectedPi); err != nil {
		t.Fatal(err)
	}
	rsv, err := rsvCtx.ApiCreateAndQueryRsv(&token, &rr, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := helpers.AssertJsonErr("rsv DC", rsv.Dc, &targetDc); err != nil {
		t.Fatal(err)
	}
	if rsv.ProcessingFee != expectedPi.ProcessingFee {
		t.Fatalf("invalid ProcessingFee: expected %d got %d",
			rsv.ProcessingFee, expectedPi.ProcessingFee)
	}
	if rsv.RefundAmount != expectedPi.RefundAmount {
		t.Fatalf("invalid RefundAmount: expected %d got %d",
			rsv.RefundAmount, expectedPi.RefundAmount)
	}
	if rsv.SplitIncomeFee != expectedPi.SplitIncomeFee {
		t.Fatalf("invalid SplitIncomeFee: expected %d got %d",
			rsv.SplitIncomeFee, expectedPi.SplitIncomeFee)
	}
	if ds, err := dcCtx.DalReadDc(dc.ReadDcRequest{
		ID:      rr.DcID,
		InstrID: &iid,
	}, nil); err != nil {
		t.Fatal(err)
	} else {
		if len(ds) != 1 {
			t.Fatal("invalid ds")
		}
		if ds[0].RedeemedQuantity != 1 {
			t.Fatal("invalid redeemed quantity")
		}
	}
}
