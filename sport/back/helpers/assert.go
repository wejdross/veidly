package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

/*
	Assert that a == b
*/
// func Assert(t *testing.T, a, b interface{}) {
// 	if err := AssertErr(a, b); err != nil {
// 		t.Fatal(err)
// 	}
// }

/*
	same as assert but with message header provided by user
*/
func AssertErr(hdr string, a, b interface{}) error {
	if _, ok := a.([]byte); ok {
		ab := a.([]byte)
		bb := b.([]byte)
		if bytes.Equal(ab, bb) {
			return nil
		}
	} else {
		if a == b {
			return nil
		}
	}
	hdr = strings.Replace(hdr, "%", "%%", -1)
	if ab, ok := a.([]byte); ok {
		as, bs := "", ""
		as = string(ab)
		if bb, ok := b.([]byte); ok {
			bs = string(bb)
		} else {
			return fmt.Errorf("attempted to compare []byte to not []byte")
		}
		return fmt.Errorf(hdr+":\n%s\n!=\n%s", as, bs)
	}
	return fmt.Errorf(hdr+":\n%v\n!=\n%v", a, b)
}

func Assert(t *testing.T, hdr string, a, b interface{}) {
	if err := AssertErr(hdr, a, b); err != nil {
		t.Fatal(err)
	}
}

func AssertJsonErr(hdr string, a, b interface{}) error {
	// if reflect.TypeOf(a) != reflect.TypeOf(b) {
	// 	return fmt.Errorf("Tried to assert ")
	// }
	aj, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		return err
	}
	bj, err := json.MarshalIndent(b, "", "    ")
	if err != nil {
		return err
	}
	return AssertErr(hdr, aj, bj)
}

func AssertJson(t *testing.T, hdr string, a, b interface{}) {
	if err := AssertJsonErr(hdr, a, b); err != nil {
		t.Fatal(err)
	}
}
