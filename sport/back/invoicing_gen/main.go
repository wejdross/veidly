package main

import (
	"log"
	"os"
	"sport/invoicing"
	"sport/rsv"
	"time"

	"github.com/google/uuid"
)

func main() {
	invoice := invoicing.Invoice{
		ID:           uuid.New(),
		Number:       "1/1023",
		InstructorID: uuid.New(),
		ObjID:        uuid.New(),
		ObjType:      rsv.AdyenSmKey,
		DateOfIssue:  time.Now(),
		YearOfIssue:  time.Now().Year(),
		DateOfSale:   time.Now(),
		SellerLines: []string{
			"line 1",
			"line 2",
			"line 3",
		},
		BuyerLines: []string{
			"line 1",
			"line 2",
			"line 3",
		},
		Paid:            4321414,
		MethodOfPayment: "transfer",
		Record: invoicing.InvoiceRecord{
			Qty:            1,
			UnitGrossPrice: 156,
			VatRate:        23,
		},
	}
	b, err := invoicing.PrintInvoice(&invoice)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("1.pdf", b, 0600)
}
