package rsv

import (
	"fmt"
	"sport/dc"
	"sport/train"
	"time"
)

// Following must always be true:
//
// TotalPrice = ProcessingFee + SplitPayout + SplitIncomeFee
//
// TrainingPrice = SplitPayout + SplitIncomeFee
type RsvPricingInfo struct {
	// this is how much user will pay
	TotalPrice int
	// estimated transaction processing costs
	ProcessingFee int
	// this is amount we will transfer for the instructor
	SplitPayout int
	// our fee for this reservation
	SplitIncomeFee int
	// refund value in case _user_ requests refund
	// (if instructor does this we have no choice but to refund 100%)
	RefundAmount int
	// discount code used - if any
	Dc *dc.Dc
	// price defined by instructor
	InstrPrice int
}

func (ctx *Ctx) GetRsvPricing(tr *train.TrainingWithJoins) RsvPricingInfo {
	tam := tr.Training.Price
	var d *dc.Dc
	if len(tr.Dcs) == 1 {
		d = &tr.Dcs[0]
		// reduce training price by discount value
		tam = tam - ((tam * tr.Dcs[0].Discount) / 100)
	}
	// processing fee
	processingFee := (ctx.Config.ProcessingFee * tam) / 100
	// // calculate our fee
	incomeFee := (ctx.Config.ServiceFee * tam) / 100
	// calculate $$$ for instructor
	payoutVal := tam - incomeFee
	// this will be refunded to user if user requests refund
	refundValue := (ctx.Config.RefundAmount * tam) / 100
	return RsvPricingInfo{
		TotalPrice:     processingFee + incomeFee + payoutVal,
		ProcessingFee:  processingFee,
		SplitPayout:    payoutVal,
		SplitIncomeFee: incomeFee,
		RefundAmount:   refundValue,
		InstrPrice:     tam,
		Dc:             d,
	}
}

func (ctx *Ctx) ValidatePricingDc(pi *RsvPricingInfo) error {
	if pi.Dc == nil {
		return nil
	}
	d := pi.Dc
	if d.RedeemedQuantity >= d.Quantity {
		return fmt.Errorf("invalid quantity")
	}
	now := time.Now().In(time.UTC)
	if time.Time(d.ValidStart).After(now) || time.Time(d.ValidEnd).Before(now) {
		return fmt.Errorf("too soon/late to redeem this discount code")
	}
	return nil
}
