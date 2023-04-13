package search_test

import (
	"fmt"
	"sport/helpers"
	"testing"
)

func TestBinSplit(t *testing.T) {

	type tc struct {
		arr            []float64
		el             float64
		expectedIx     int
		strictlyLesser bool
	}

	a1 := []float64{0, 1, 2, 3, 5, 6, 8, 10, 11}
	tcs := []tc{
		{ // 0
			arr:            a1,
			el:             6,
			expectedIx:     5,
			strictlyLesser: false,
		},
		{ // 1
			arr:            a1,
			el:             6,
			expectedIx:     4,
			strictlyLesser: true,
		},
		{ // 2
			arr:            a1,
			el:             20,
			expectedIx:     8,
			strictlyLesser: true,
		},
		{ // 3
			arr:            a1,
			el:             0,
			expectedIx:     0,
			strictlyLesser: false,
		},
		{ // 4
			arr:            a1,
			el:             -1,
			expectedIx:     -1,
			strictlyLesser: true,
		},
		{ // 5
			arr:            a1,
			el:             0,
			expectedIx:     -1,
			strictlyLesser: true,
		},
	}

	for i := range tcs {
		s := helpers.BinSplitAround(len(tcs[i].arr), func(ix int) bool {
			if tcs[i].strictlyLesser {
				return tcs[i].arr[ix] < tcs[i].el
			}
			return tcs[i].arr[ix] <= tcs[i].el
		})
		if s != tcs[i].expectedIx {
			t.Fatal(fmt.Errorf("Invalid split @ index %d. Expected %d, got %d",
				i, tcs[i].expectedIx, s))
		}
	}

}

// func TestSearchCache(t *testing.T) {

// 	itoken, iid, err := searchCtx.Instr.CreateTestInstructorWithUser()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	sc, err := searchCtx.GetSearchCache()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_ = iid
// 	_ = itoken

// 	//fmt.Println(helpers.JsonMustSerializeFormatStr(sc))
// }
