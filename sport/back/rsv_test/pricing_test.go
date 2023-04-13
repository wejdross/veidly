package rsv_test

import (
	"fmt"
	"sport/helpers"
	"sport/rsv"
	"sport/train"
	"testing"
	"time"
)

func testPricing(token string, price int, expectedPricingInfo rsv.RsvPricingInfo) error {
	s := rsv.NowMin().Add(time.Hour)
	tr := train.NewTestCreateTrainingRequest(s, s.Add(time.Hour))
	tr.Training.AllowExpress = true
	tr.Training.Price = price
	tid, err := rsvCtx.Train.ApiCreateTraining(token, &tr)
	if err != nil {
		return err
	}
	rr := rsv.NewTestCreateReservationRequest(tid, s)
	pi, err := rsvCtx.ApiGetRsvPricing(&rr, 0)
	if err != nil {
		return err
	}
	if err = helpers.AssertJsonErr("invalid RsvPricing", &pi, &expectedPricingInfo); err != nil {
		return err
	}
	rsv, err := rsvCtx.ApiCreateAndQueryRsv(&token, &rr, 0)
	if err != nil {
		return err
	}
	if rsv.ProcessingFee != expectedPricingInfo.ProcessingFee {
		return fmt.Errorf("invalid ProcessingFee: expected %d got %d",
			rsv.ProcessingFee, expectedPricingInfo.ProcessingFee)
	}
	if rsv.RefundAmount != expectedPricingInfo.RefundAmount {
		return fmt.Errorf("invalid RefundAmount: expected %d got %d",
			rsv.RefundAmount, expectedPricingInfo.RefundAmount)
	}
	if rsv.SplitIncomeFee != expectedPricingInfo.SplitIncomeFee {
		return fmt.Errorf("invalid SplitIncomeFee: expected %d got %d",
			rsv.SplitIncomeFee, expectedPricingInfo.SplitIncomeFee)
	}
	return nil
}

// this test wont take discount codes into account
// there is separate test in dc_test module which takes care of that
func TestRsvPricing(t *testing.T) {
	token, _, err := rsvCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	sf, pf, ra := 0, 0, 0
	sf = rsvCtx.Config.ServiceFee
	pf = rsvCtx.Config.ProcessingFee
	ra = rsvCtx.Config.RefundAmount
	rsvCtx.Config.ServiceFee = 10
	rsvCtx.Config.ProcessingFee = 3
	rsvCtx.Config.RefundAmount = 30
	defer func() {
		rsvCtx.Config.ServiceFee = sf
		rsvCtx.Config.ProcessingFee = pf
		rsvCtx.Config.RefundAmount = ra
	}()

	if err := testPricing(token, 100, rsv.RsvPricingInfo{
		TotalPrice:     103,
		ProcessingFee:  3,
		SplitPayout:    90,
		SplitIncomeFee: 10,
		RefundAmount:   30,
		Dc:             nil,
		InstrPrice:     100,
	}); err != nil {
		t.Fatal(err)
	}

	if err := testPricing(token, 150, rsv.RsvPricingInfo{
		TotalPrice:     154,
		ProcessingFee:  4,
		SplitPayout:    135,
		SplitIncomeFee: 15,
		RefundAmount:   45,
		Dc:             nil,
		InstrPrice:     150,
	}); err != nil {
		t.Fatal(err)
	}

	if err := testPricing(token, 1234, rsv.RsvPricingInfo{
		TotalPrice:     1271,
		ProcessingFee:  37,
		SplitPayout:    1111,
		SplitIncomeFee: 123,
		RefundAmount:   370,
		Dc:             nil,
		InstrPrice:     1234,
	}); err != nil {
		t.Fatal(err)
	}

	if err := testPricing(token, 24785, rsv.RsvPricingInfo{
		TotalPrice:     25528,
		ProcessingFee:  743,
		SplitPayout:    22307,
		SplitIncomeFee: 2478,
		RefundAmount:   7435,
		Dc:             nil,
		InstrPrice:     24785,
	}); err != nil {
		t.Fatal(err)
	}
}
