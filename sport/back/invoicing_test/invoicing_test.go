package invoicing_test

import (
	"reflect"
	"sport/instr"
	"sport/invoicing"
	"sport/rsv"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestIntToTriplets(t *testing.T) {
	var testCases = []struct {
		name string
		arg  int
		want []int
	}{
		{
			name: "1",
			arg:  1,
			want: []int{1},
		},
		{
			name: "2",
			arg:  123344445555,
			want: []int{555, 445, 344, 123},
		},
	}
	for i := range testCases {
		tt := &testCases[i]
		t.Run(tt.name, func(t *testing.T) {
			got := invoicing.IntToTriplets(tt.arg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntToTriplets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvoice(t *testing.T) {

	_, iid, err := invoicingCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().Round(time.Hour).In(time.FixedZone("", 0))

	invoice := invoicing.Invoice{
		ID:           uuid.New(),
		Number:       uuid.NewString(),
		InstructorID: iid,
		ObjID:        uuid.New(),
		DateOfIssue:  now,
		YearOfIssue:  now.Year(),
		DateOfSale:   now,
		SellerLines:  invoicingCtx.Config.CompanyLines,
		BuyerLines: []string{
			"1",
			"2",
			"3",
		},
		Paid:            1234,
		MethodOfPayment: "transfer",
		Record: invoicing.InvoiceRecord{
			Qty:            1,
			UnitGrossPrice: 12340,
			VatRate:        23,
		},
	}

	if err := invoicingCtx.DalCreateInvoice(nil, &invoice); err != nil {
		t.Fatal(err)
	}

	invoices, err := invoicingCtx.DalGetInvoices(invoice.ID, invoice.InstructorID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = invoicingCtx.DalGetInvoices(uuid.Nil, invoice.InstructorID)
	if err != nil {
		t.Fatal(err)
	}

	if len(invoices) != 1 {
		t.Fatal("obtained invalid number of invoices")
	}

	gotInvoice := invoices[0]

	if !reflect.DeepEqual(&invoice, &gotInvoice) {
		t.Errorf("DalGetInvoices() = \n%v\nwant\n%v", gotInvoice, invoice)
	}

	_, err = invoicing.PrintInvoice(&invoice)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateInvoice(t *testing.T) {
	price := 1234
	token, err := invoicingCtx.Instr.User.ApiCreateAndLoginUser(nil)
	if err != nil {
		t.Fatal(err)
	}
	uid, err := invoicingCtx.Instr.User.Api.AuthorizeUserFromToken(token)
	if err != nil {
		t.Fatal(err)
	}

	instrLines := []string{
		uuid.NewString(),
		uuid.NewString(),
		uuid.NewString(),
	}

	iid, err := invoicingCtx.Instr.CreateTestInstructor(uid, &instr.InstructorRequest{
		InvoiceLines: instrLines,
	})
	if err != nil {
		t.Fatal(err)
	}
	objectID := uuid.New()

	err = invoicingCtx.CreateInvoice(price, iid, objectID, rsv.AdyenSmKey)
	if err != nil {
		t.Fatal(err)
	}

	invoices, err := invoicingCtx.DalGetInvoices(uuid.Nil, iid)
	if err != nil {
		t.Fatal(err)
	}

	if len(invoices) != 1 {
		t.Fatal("obtained invalid number of invoices")
	}

	invoice := invoices[0]

	expectedInvoice := invoicing.Invoice{
		InstructorID:    iid,
		ObjID:           objectID,
		BuyerLines:      instrLines,
		Paid:            price,
		MethodOfPayment: "transfer",
		Record: invoicing.InvoiceRecord{
			Qty:            1,
			UnitGrossPrice: price,
			VatRate:        23,
		},
		ObjType: rsv.AdyenSmKey,
		// i don't know how to predict those fields
		ID:          invoice.ID,
		Number:      invoice.Number,
		DateOfIssue: invoice.DateOfIssue,
		YearOfIssue: invoice.DateOfIssue.Year(),
		DateOfSale:  invoice.DateOfSale,
		SellerLines: invoice.SellerLines,
	}

	if !reflect.DeepEqual(&invoice, &expectedInvoice) {
		t.Errorf("DalGetInvoices() = \n%v\nwant\n%v", invoice, expectedInvoice)
	}
}
