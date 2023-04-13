package sub_qr

import (
	"github.com/google/uuid"
)

func (ctx *Ctx) DalCanCreateQr(qr *QrCode) (bool, error) {
	const q = "select count(1) from sub_qr_codes where sub_id = $1"
	rw := ctx.Dal.Db.QueryRow(q, qr.SubID)
	var c int
	if err := rw.Scan(&c); err != nil {
		return false, err
	}
	if c > ctx.Config.MaxQrCodes {
		return false, nil
	}
	return true, nil
}

func (ctx *Ctx) DalCreateQr(qr *QrCode) error {
	const q = `insert into sub_qr_codes (id, sub_id) values ($1, $2)`
	_, err := ctx.Dal.Db.Query(q, qr.ID, qr.SubID)
	return err
}

func (ctx *Ctx) DalReadQrForRsv(rsvID uuid.UUID) ([]QrCode, error) {
	const q = `select id, sub_id from sub_qr_codes where sub_id = $1`
	dbr, err := ctx.Dal.Db.Query(q, rsvID)
	if err != nil {
		return nil, err
	}
	defer dbr.Close()
	l := 5
	if ctx.Config.MaxQrCodes < l {
		l = ctx.Config.MaxQrCodes
	}
	res := make([]QrCode, 0, ctx.Config.MaxQrCodes)
	var tmp QrCode
	for dbr.Next() {
		if err := dbr.Scan(&tmp.ID, &tmp.SubID); err != nil {
			return nil, err
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (ctx *Ctx) DalReadSingleQr(id uuid.UUID) (*QrCode, error) {
	const q = `select id, sub_id from sub_qr_codes where id = $1`
	dbr := ctx.Dal.Db.QueryRow(q, id)
	res := new(QrCode)
	if err := dbr.Scan(&res.ID, &res.SubID); err != nil {
		return nil, err
	}
	return res, nil
}
