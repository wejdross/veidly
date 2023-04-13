package helpers

import (
	"fmt"
	"math"
	"reflect"
	"runtime"

	"github.com/google/uuid"
)

// -1 on not found
func BinSplitAround(
	len int,
	lesserCmp func(ix int) bool,
) (splitAround int) {
	l := 0
	r := len - 1
	resIX := -1
	for l <= r {
		m := (l + r) / 2
		if lesserCmp(m) { // can be < or <=
			l = m + 1
			resIX = m
		} else { // > or >=
			r = m - 1
		}
	}
	return resIX
}

func DgToRad(dg float64) float64 {
	return (dg * math.Pi) / 180
}

// dont overdo it with the number of functions
func RunFunctionsInParallel(fns []func() error) error {
	erc := make(chan error, len(fns))
	for i := range fns {
		go func(i int) {
			erc <- fns[i]()
		}(i)
	}
	muerr := NewMuxError(fmt.Sprintf("RunFunctionsInParallel(%d):", len(fns)))
	for i := 0; i < len(fns); i++ {
		err := <-erc
		if err != nil {
			fn := runtime.FuncForPC(reflect.ValueOf(fns[i]).Pointer()).Name()
			muerr.Add(fmt.Errorf("%s: %v", fn, err))
		}
	}
	if muerr.HasErrors() {
		return muerr
	}
	return nil
}

func IDMapToStringArr(idm map[uuid.UUID]struct{}) []string {
	iids := make([]string, len(idm))
	var i int
	for x := range idm {
		iids[i] = x.String()
		i++
	}
	return iids
}
