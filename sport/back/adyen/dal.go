package adyen

import (
	"sport/helpers"
	"time"
)

func (ctx *Ctx) LogDonation(email string, val int) error {
	const q = `insert into donations 
		(email,donation_amount, timestamp) 
		values ($1, $2, $3)`
	res, err := ctx.Dal.Db.Exec(q, email, val, time.Now())
	if err != nil {
		return err
	}
	if err := helpers.PgMustBeOneRow(res); err != nil {
		return err
	}
	return nil
}
