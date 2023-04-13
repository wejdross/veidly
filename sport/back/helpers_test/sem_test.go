package helpers_test

import (
	"fmt"
	"math"
	"sport/helpers"
	"testing"
	"time"
)

// predict total semaphore duration
//
// itd is iteration duration
func TotalSemDur(itd time.Duration, threads, max int) time.Duration {
	noit := int(math.Ceil(float64(max) / float64(threads)))
	return time.Duration(noit) * itd
}

func testSem(th, mx int, fns []func(int) error, timeErr, itd time.Duration) error {
	s := time.Now()
	if err := helpers.Sem(th, mx, fns); err != nil {
		return err
	}
	df := time.Now().Sub(s)
	df = df.Round(time.Microsecond)
	pred := TotalSemDur(itd, th, mx)
	x := math.Abs(float64(df - pred))
	if time.Duration(x) > timeErr {
		return fmt.Errorf(
			"invalid duration. expected semaphore(%d, %d) to run %v, but it ran for %v",
			th, mx, pred, df)
	}
	fmt.Printf(
		"semaphore(%d, %d) ran %8v, expected %v, err: %3v / %3v\n",
		th, mx, df, pred, time.Duration(x), timeErr)
	return nil
}

func TestSem(t *testing.T) {
	itd := time.Millisecond * 20
	terr := time.Millisecond * 2
	fns := []func(int) error{func(int) error {
		time.Sleep(itd)
		return nil
	}}

	// this should still work if numcpu is smaller than th

	th := 3
	mx := 6
	if err := testSem(th, mx, fns, terr, itd); err != nil {
		t.Fatal(err)
	}

	th = 4
	mx = 7
	if err := testSem(th, mx, fns, terr, itd); err != nil {
		t.Fatal(err)
	}

	th = 4
	mx = 9
	if err := testSem(th, mx, fns, terr, itd); err != nil {
		t.Fatal(err)
	}

	th = 1
	mx = 3
	if err := testSem(th, mx, fns, terr, itd); err != nil {
		t.Fatal(err)
	}

	th = 2
	mx = 2
	if err := testSem(th, mx, fns, terr, itd); err != nil {
		t.Fatal(err)
	}

	th = 2
	mx = 4
	if err := testSem(th, mx, fns, terr, itd); err != nil {
		t.Fatal(err)
	}

	// itd = time.Millisecond
	// terr = time.Millisecond
	// fns = []func(int) error{func(int) error {
	// 	time.Sleep(itd)
	// 	return nil
	// }}
	// th = 40
	// mx = 100000
	// if err := testSem(th, mx, fns, terr, itd); err != nil {
	// 	t.Fatal(err)
	// }
}
