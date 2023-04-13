package dc_test

import (
	"database/sql"
	"fmt"
	"sport/api"
	"sport/dc"
	"sport/helpers"
	"sport/user"
	"testing"
	"time"
)

func TestUse(t *testing.T) {
	token, _, err := instrCtx.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}
	s := time.Now().In(time.UTC).Add(-time.Hour).Round(time.Second)
	const q = 30
	m := "POST"
	u := "/api/dc"
	req := &dc.DcRequest{
		Name:       "foo",
		Quantity:   q,
		ValidStart: dc.Time(s),
		ValidEnd:   dc.Time(s.Add(time.Hour)),
		Discount:   50,
	}
	req2 := &dc.DcRequest{
		Name:       "bar",
		Quantity:   q,
		ValidStart: dc.Time(s),
		ValidEnd:   dc.Time(s.Add(time.Hour)),
		Discount:   50,
	}
	dcCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(req),
			RequestHeaders:     user.GetAuthHeader(token),
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(req2),
			RequestHeaders:     user.GetAuthHeader(token),
		},
	})

	m = "GET"
	var dcs []dc.Dc

	dcCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				helpers.JsonMustDeserialize(b, &dcs)
				if len(dcs) != 2 {
					return fmt.Errorf("invalid dcs")
				}
				return nil
			},
		},
	})

	c := q
	w := make(chan error, c)

	for i := 0; i < c; i++ {
		go func() {
			e := dcCtx.DalUpdateDcUses(dcs[0].ID, 1)
			w <- e
		}()
	}

	err = nil
	for i := 0; i < c; i++ {
		e := <-w
		if e != nil && err == nil {
			err = e
		}
	}

	if err != nil {
		t.Fatal(err)
	}

	//

	c = q + 1
	w = make(chan error, c)

	for i := 0; i < c; i++ {
		go func() {
			e := dcCtx.DalUpdateDcUses(dcs[1].ID, 1)
			w <- e
		}()
	}

	err = nil
	for i := 0; i < c; i++ {
		e := <-w
		if e != nil && err == nil {
			err = e
		}
	}

	if err != sql.ErrNoRows {
		t.Fatal("expected error got nil")
	}

	c = q
	w = make(chan error, c)

	for i := 0; i < c; i++ {
		go func() {
			e := dcCtx.DalUpdateDcUses(dcs[1].ID, -1)
			w <- e
		}()
	}

	err = nil
	for i := 0; i < c; i++ {
		e := <-w
		if e != nil && err == nil {
			err = e
		}
	}

	if err != nil {
		t.Fatal(err)
	}

	e := dcCtx.DalUpdateDcUses(dcs[1].ID, -1)
	if e == nil {
		t.Fatal("expected error got nil")
	}
}
