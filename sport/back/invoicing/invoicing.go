package invoicing

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"sport/helpers"
	"sport/instr"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/lib/pq"
)

type InvoiceRecord struct {
	Qty            int
	UnitGrossPrice int
	VatRate        int
}

type Invoice struct {
	// IMPORTANT:
	// invoice ID must be kept secret from other users
	ID              uuid.UUID
	Number          string
	InstructorID    uuid.UUID
	ObjID           uuid.UUID
	ObjType         string
	DateOfIssue     time.Time
	YearOfIssue     int
	DateOfSale      time.Time
	SellerLines     []string
	BuyerLines      []string
	Paid            int
	MethodOfPayment string
	Record          InvoiceRecord
}

func (ctx *Ctx) DalGetNextInvoiceNumber(
	tx *sql.Tx, yearOfIssue int) (string, error) {

	q := `select count(*)
		from invoices
		where year_of_issue = $1`

	row := tx.QueryRow(q, yearOfIssue)

	count := 0

	err := row.Scan(&count)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d/%d", count+1, yearOfIssue), nil
}

func (ctx *Ctx) DalCreateInvoice(tx *sql.Tx, invoice *Invoice) error {

	const q = `insert into invoices (
		id,
		number,
		instructor_id,
		obj_type,
		obj_id,
		date_of_issue,
		year_of_issue,
		date_of_sale,
		seller,
		buyer,
		paid,
		method_of_payment,
		rec_qty,
		rec_unit_gross_price,
		rec_vat_rate
	) values (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13,
		$14,
		$15
	)`

	args := []interface{}{
		invoice.ID,
		invoice.Number,
		invoice.InstructorID,
		invoice.ObjType,
		invoice.ObjID,
		invoice.DateOfIssue,
		invoice.YearOfIssue,
		invoice.DateOfSale,
		pq.StringArray(invoice.SellerLines),
		pq.StringArray(invoice.BuyerLines),
		invoice.Paid,
		invoice.MethodOfPayment,
		invoice.Record.Qty,
		invoice.Record.UnitGrossPrice,
		invoice.Record.VatRate,
	}

	var err error
	if tx != nil {
		_, err = tx.Exec(q, args...)
	} else {
		_, err = ctx.Dal.Db.Exec(q, args...)
	}
	return err
}

func (ctx *Ctx) DalGetInvoices(id, instructorID uuid.UUID) ([]Invoice, error) {
	qb := strings.Builder{}
	qb.WriteString(`select
		id,
		number,
		instructor_id,
		obj_type,
		obj_id,
		date_of_issue,
		year_of_issue,
		date_of_sale,
		seller,
		buyer,
		paid,
		method_of_payment,
		rec_qty,
		rec_unit_gross_price,
		rec_vat_rate
	from invoices
	where `)

	args := make([]interface{}, 0, 2)

	if instructorID != uuid.Nil {
		qb.WriteString("instructor_id = $1")
		args = append(args, instructorID)
	}

	if id != uuid.Nil {
		if len(args) > 0 {
			qb.WriteString("and id = $2")
		} else {
			qb.WriteString("id = $1")
		}
		args = append(args, id)
	}

	dbres, err := ctx.Dal.Db.Query(qb.String(), args...)
	if err != nil {
		return nil, err
	}

	defer dbres.Close()

	res := make([]Invoice, 0, 2)

	for dbres.Next() {
		res = append(res, Invoice{})
		r := &(res[len(res)-1])
		if err := dbres.Scan(
			&r.ID,
			&r.Number,
			&r.InstructorID,
			&r.ObjType,
			&r.ObjID,
			&r.DateOfIssue,
			&r.YearOfIssue,
			&r.DateOfSale,
			(*pq.StringArray)(&r.SellerLines),
			(*pq.StringArray)(&r.BuyerLines),
			&r.Paid,
			&r.MethodOfPayment,
			&r.Record.Qty,
			&r.Record.UnitGrossPrice,
			&r.Record.VatRate,
		); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (ctx *Ctx) CreateInvoice(
	invoicePrice int,
	instructorID, objectID uuid.UUID,
	objType string) error {

	instr, err := ctx.Instr.DalReadInstructor(instructorID, instr.InstructorID)
	if err != nil {
		return err
	}

	// serializable may be too much in here... but screw it
	tx, err := ctx.Dal.Db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return err
	}

	now := time.Now()

	invoice := Invoice{
		ID:              uuid.New(),
		DateOfIssue:     now,
		YearOfIssue:     now.Year(),
		DateOfSale:      now,
		SellerLines:     ctx.Config.CompanyLines,
		Paid:            invoicePrice,
		BuyerLines:      instr.InvoiceLines,
		InstructorID:    instructorID,
		ObjType:         objType,
		ObjID:           objectID,
		MethodOfPayment: "transfer",
		Record: InvoiceRecord{
			UnitGrossPrice: invoicePrice,
			Qty:            1,
			VatRate:        23,
		},
	}

	for i := 0; i < 50; i++ {
		number, err := ctx.DalGetNextInvoiceNumber(tx, invoice.YearOfIssue)
		if err != nil {
			tx.Rollback()
			if helpers.PgIsConcurrentUpdate(err) {
				continue
			}
			return err
		}

		invoice.Number = number

		if err := ctx.DalCreateInvoice(tx, &invoice); err != nil {
			tx.Rollback()
			if helpers.PgIsConcurrentUpdate(err) {
				continue
			}
			return err
		}

		if err = tx.Commit(); err != nil {
			tx.Rollback()
			if helpers.PgIsConcurrentUpdate(err) {
				continue
			}
			return err
		}

		return nil
	}

	return fmt.Errorf("DalGetNextPenaltySerial: cannot serialize transaction -- too many requests")
}

func max(x ...int) int {
	m := 0
	for i := range x {
		if x[i] > m {
			m = x[i]
		}
	}
	return m
}

func spaces(len int) string {
	rb := strings.Builder{}
	for i := 0; i < len; i++ {
		rb.WriteString(" ")
	}
	return rb.String()
}

func PrintInvoice(invoice *Invoice) ([]byte, error) {

	const resourceBasePath = "../invoicing/"

	m := pdf.NewMaroto(consts.Portrait, consts.Letter)
	const ffamily = "Roboto"
	m.AddUTF8Font(ffamily, consts.Normal, path.Join(resourceBasePath, "Roboto-Regular.ttf"))
	m.AddUTF8Font(ffamily, consts.Italic, path.Join(resourceBasePath, "Roboto-Italic.ttf"))
	m.AddUTF8Font(ffamily, consts.Bold, path.Join(resourceBasePath, "Roboto-Bold.ttf"))
	m.AddUTF8Font(ffamily, consts.BoldItalic, path.Join(resourceBasePath, "Roboto-BoldItalic.ttf"))

	m.Row(10, func() {
		m.Text("Faktura VAT nr. "+invoice.Number, props.Text{
			Size:   20,
			Align:  consts.Center,
			Family: ffamily,
		})
	})

	m.Row(5, func() {})

	m.Row(20, func() {
		m.Text(
			"Data wystawienia: "+invoice.DateOfIssue.Format("2006-01-02"),
			props.Text{
				Top:         0,
				Align:       consts.Right,
				Size:        12,
				Extrapolate: true,
				Family:      ffamily,
			})
		m.Text(
			"Data spzedaży: "+invoice.DateOfSale.Format("2006-01-02"),
			props.Text{
				Top:         6,
				Align:       consts.Right,
				Size:        12,
				Extrapolate: true,
				Family:      ffamily,
			})
	})

	m.Row(25, func() {
		m.Col(6, func() {
			m.Text("Sprzedawca", props.Text{
				Size:   18,
				Family: ffamily,
				Align:  consts.Center,
			})
			for i := range invoice.SellerLines {
				m.Text(
					invoice.SellerLines[i],
					props.Text{
						Top:         float64(6*(i) + 10),
						Size:        12,
						Extrapolate: true,
						Family:      ffamily,
						Align:       consts.Center,
					})
			}
		})
		m.Col(6, func() {
			m.Text("Nabywca", props.Text{
				Size:   18,
				Family: ffamily,
				Align:  consts.Center,
			})
			for i := range invoice.BuyerLines {
				m.Text(
					invoice.BuyerLines[i],
					props.Text{
						Top:         float64(6*(i) + 10),
						Align:       consts.Center,
						Size:        12,
						Extrapolate: true,
						Family:      ffamily,
					})
			}
		})
	})

	m.Row(5, func() {})
	m.Line(10)

	price := invoice.Record.UnitGrossPrice

	gross := float64(price) / 100

	grossString := strconv.FormatFloat(
		gross,
		'f',
		2,
		64,
	)

	net := float64(price) / float64(100+invoice.Record.VatRate)

	netString := strconv.FormatFloat(
		net,
		'f',
		2,
		64,
	)

	vatAmountString := strconv.FormatFloat(
		gross-net,
		'f',
		2,
		64,
	)

	maxOffset := max(len(netString), len(grossString), len(vatAmountString))

	defOff := 20

	if maxOffset > defOff {
		defOff = maxOffset + 1
	}

	m.Row(30, func() {
		m.Col(6, func() {
			m.Text("Usługi marketingowe", props.Text{
				Size:   12,
				Align:  consts.Left,
				Family: ffamily,
			})
		})
		m.Col(6, func() {
			m.Text(
				fmt.Sprintf("Netto %s %s zł",
					spaces(defOff+(maxOffset-len(netString))),
					netString),
				props.Text{
					Size:   12,
					Align:  consts.Right,
					Family: ffamily,
				})
			vatStr := fmt.Sprintf("%d%%", invoice.Record.VatRate)
			spacesLen := defOff + 2*(maxOffset-len(vatAmountString)) - 2*len(vatStr)
			left := spacesLen / 2
			right := spacesLen - left
			m.Text(
				fmt.Sprintf("VAT %s%s%s %s zł",
					spaces(left), vatStr, spaces(right), vatAmountString),
				props.Text{
					Top:    6,
					Size:   12,
					Align:  consts.Right,
					Family: ffamily,
				})
			m.Text(
				fmt.Sprintf("Brutto %s %s zł",
					spaces(defOff+(maxOffset-len(grossString))),
					grossString),
				props.Text{
					Top:    12,
					Size:   12,
					Align:  consts.Right,
					Family: ffamily,
				})
			m.Text(
				fmt.Sprintf("Razem %s %s zł",
					spaces(defOff+(maxOffset-len(grossString))),
					grossString),
				props.Text{
					Top:    18,
					Size:   12,
					Align:  consts.Right,
					Family: ffamily,
				})
			total := price / 100
			rem := price % 100
			m.Text(
				fmt.Sprintf(
					"%s zł. %d/100",
					PLNToWords(total),
					rem,
				),
				props.Text{
					Top:    24,
					Size:   12,
					Align:  consts.Right,
					Family: ffamily,
					Style:  consts.Italic,
				})
		})
	})

	m.Col(12, func() {
		m.Row(10, func() {
			m.Text("FAKTURA OPŁACONA", props.Text{
				Size:   18,
				Top:    40,
				Align:  consts.Center,
				Family: ffamily,
			})
		})
	})

	m.Row(10, func() {})

	m.Col(12, func() {
		m.Row(40, func() {
			_ = m.FileImage(path.Join(resourceBasePath, "veidly.png"), props.Rect{
				Center:  true,
				Percent: 80,
			})
		})
	})

	res, err := m.Output()

	if err != nil {
		return nil, err
	}

	return res.Bytes(), err
}
