package helpers

import (
	"errors"
	"fmt"
	"strings"
)

func FnRpt(fn func() error, c int) []func() error {
	fns := make([]func() error, c)
	for i := 0; i < c; i++ {
		fns[i] = fn
	}
	return fns
}

// run specified functions in semaphore
//
// maxth, max must be > 0
//
// len(fn) can be any value but it must be > 0
func Sem(maxth, max int, fn []func(i int) error) error {

	semChan := make(chan struct{}, maxth)

	type res struct {
		err error
		ix  int
	}
	waiterChan := make(chan res, maxth)

	go func() {
		for i := 0; i < max; i++ {
			semChan <- struct{}{}
			go func(_i int) {
				err := fn[_i%len(fn)](_i)
				<-semChan
				waiterChan <- res{
					err, _i,
				}
			}(i)
		}
	}()

	errs := make([]res, 0, 4)

	for i := 0; i < max; i++ {
		r := <-waiterChan
		if r.err != nil {
			errs = append(errs, r)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Sem: %d errors occurred:", len(errs)))
	for i := range errs {
		sb.WriteString(fmt.Sprintf("\n\t%d: %v", errs[i].ix, errs[i].err))
	}
	return errors.New(sb.String())
}
