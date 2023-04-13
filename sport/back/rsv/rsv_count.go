package rsv

import (
	"fmt"
	"sport/helpers"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RsvCountKey struct {
	TrainID   uuid.UUID
	DateStart int64
	DateEnd   int64
}

// can register rsv count index
type RsvCount map[RsvCountKey]int

func RsvToCountKey(r *Rsv) RsvCountKey {
	var tid uuid.UUID
	if r.TrainingID != nil {
		tid = *r.TrainingID
	}
	return RsvCountKey{
		TrainID:   tid,
		DateStart: r.DateStart.UnixNano(),
		DateEnd:   r.DateEnd.UnixNano(),
	}
}

func ShortRsvToCountKey(r *ShortenedRsv) RsvCountKey {
	var tid uuid.UUID
	if r.TrainingID != nil {
		tid = *r.TrainingID
	}
	return RsvCountKey{
		TrainID:   tid,
		DateStart: r.DateStart.UnixNano(),
		DateEnd:   r.DateEnd.UnixNano(),
	}
}

func NewRsvCount(rr []RsvWithInstr) RsvCount {
	rcix := make(RsvCount)
	for i := 0; i < len(rr); i++ {
		if !rr[i].IsConfirmed {
			continue
		}
		key := RsvToCountKey(&rr[i].Rsv)
		rcix[key]++
	}
	return rcix
}

func (ctx *Ctx) NewRsvCountFromShortRsvs(rr []ShortenedRsv) RsvCount {
	rcix := make(RsvCount)
	for i := 0; i < len(rr); i++ {
		key := ShortRsvToCountKey(&rr[i])
		rcix[key]++
	}
	return rcix
}

func (ctx *Ctx) NewDalRsvCount(dr *helpers.DateRange) (RsvCount, error) {
	q := strings.Builder{}
	q.WriteString(`
		select r.training_id, r.date_start, r.date_end, count(1)
		from reservations r
		where is_confirmed = true
	`)
	args := make([]interface{}, 0, 2)

	if dr.IsNotZero() {
		next := len(args) + 1
		q.WriteString(fmt.Sprintf(` and r.date_start <= $%d `, next))
		next++
		q.WriteString(fmt.Sprintf(` and r.date_end >= $%d `, next))
		args = append(args,
			dr.End.In(time.UTC),
			dr.Start.In(time.UTC))
		next++
	}

	q.WriteString(" group by r.training_id, r.date_start, r.date_end ")

	r, err := ctx.Dal.Db.Query(q.String(), args...)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	res := make(RsvCount)

	var start, end time.Time
	var trainingID uuid.UUID
	var c int

	for r.Next() {
		if err := r.Scan(&trainingID, &start, &end, &c); err != nil {
			return nil, err
		}
		res[RsvCountKey{
			TrainID:   trainingID,
			DateStart: start.UnixNano(),
			DateEnd:   end.UnixNano(),
		}] += c
	}

	return res, nil
}
