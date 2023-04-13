package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func EnsuretTimeInRange(t1, t2 time.Time) error {
	df := t1.Sub(t2)
	const errorCorrection = time.Minute * 2
	if df <= -errorCorrection || df >= errorCorrection {
		return fmt.Errorf("EnsuretTimeInRange: out of range, expected %v got: %v\ndf was %v", t1, t2, df)
	}
	return nil
}

func NowMin() time.Time {
	now := time.Now().In(time.UTC)
	now = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		0, 0, time.UTC)
	return now
}

type DateRange struct {
	Start time.Time
	End   time.Time
}

func (dr DateRange) IsNotZero() bool {
	if dr.Start.IsZero() || dr.End.IsZero() {
		return false
	}
	return true
}

var ZeroDateRange = DateRange{}

func DateRangeFromQueryString(ctx *gin.Context) (DateRange, error) {
	start := ctx.Query("start")
	end := ctx.Query("end")
	var t int64
	var err error
	var ret DateRange
	if start != "" {
		if t, err = strconv.ParseInt(start, 10, 64); err != nil {
			return ZeroDateRange, err
		}
		ret.Start = time.Unix(t, 0)
	}
	if end != "" {
		if t, err = strconv.ParseInt(end, 10, 64); err != nil {
			return ZeroDateRange, err
		}
		ret.End = time.Unix(t, 0)
	}
	if (ret.Start.IsZero()) != (ret.End.IsZero()) {
		return ZeroDateRange, fmt.Errorf("Validate DateRange: start XOR end failed")
	}
	if ret.IsNotZero() {
		if ret.End.Before(ret.Start) {
			return ZeroDateRange, fmt.Errorf("Validate DateRange: end is before start")
		}
	}
	return ret, nil
}

func Overlaps(start1, end1, start2, end2 time.Time) bool {
	return (start1.Before(end2)) && (end1.After(start2))
}

func OverlapsInt(start1, end1, start2, end2 int) bool {
	return start1 < end2 && end1 > start2
}
