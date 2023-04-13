package adyen_sm

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type WhLog struct {
	ID            uuid.UUID
	ReservationID uuid.UUID
	Timestamp     time.Time
	Message       []byte
	Error         string
	IsSuccess     bool
}

func (ctx *Ctx) LogWhNotification(l *WhLog) {

	if ctx.noReplyCtx != nil && !l.IsSuccess {
		go ctx.MailSupportAboutWHNotification(l)
	}

	_, err := ctx.dal.Db.Exec(
		`insert into wh_log (
			id, reservation_id, timestamp, message, error, is_success
		) values ($1, $2, $3, $4, $5, $6)`,
		l.ID, l.ReservationID, l.Timestamp, l.Message, l.Error, l.IsSuccess)

	if err == nil {
		return
	}

	if _, err := fmt.Fprint(os.Stderr, err); err != nil {
		panic(err)
	}

	ctx.LogEventToFs(l, l.ID, l.Timestamp)
}
