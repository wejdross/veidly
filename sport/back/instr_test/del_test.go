package instr_test

import (
	"sport/api"
	"sport/helpers"
	"sport/instr"
	"testing"
	"time"
)

func TestDeletedInstructor(t *testing.T) {

	token, err := instrCtx.User.ApiCreateAndLoginUser(nil)
	if err != nil {
		t.Fatal(err)
	}

	id, err := instrCtx.Api.AuthorizeUserFromToken(token)
	if err != nil {
		t.Fatal(err)
	}

	ir := instr.InstructorRequest{}
	i := ir.NewInstructor(id)
	i.Refunds = make(instr.InstructorRefunds)
	i.Refunds[time.Now().Unix()] = 1
	i.QueuedPayoutCuts = 2

	if err := instrCtx.DalCreateInstructor(i); err != nil {
		t.Fatal(err)
	}

	instrCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod: "DELETE",
			RequestUrl:    "/api/instructor",
			RequestHeaders: map[string]string{
				"Authorization": "Bearer " + token,
			},
			ExpectedStatusCode: 204,
		},
	})

	_, err = instrCtx.DalReadInstructor(id, instr.UserID)
	if err == nil {
		t.Fatal("instructor should be gone now")
	}

	if !helpers.IsENF(err) {
		t.Fatal("instructor should be gone now")
	}

	if err := instrCtx.ApiCreateInstructor(token, nil); err != nil {
		t.Fatal(err)
	}

	ires, err := instrCtx.DalReadInstructor(id, instr.UserID)
	if err != nil {
		t.Fatal(err)
	}

	if ires.QueuedPayoutCuts != i.QueuedPayoutCuts {
		t.Fatalf("invalid QueuedPayoutCuts, expected %d got %d",
			i.QueuedPayoutCuts, ires.QueuedPayoutCuts)
	}

	if len(ires.Refunds) != len(i.Refunds) {
		t.Fatalf("invalid Refunds, expected len %d got %d",
			len(i.Refunds), len(ires.Refunds))
	}

	for x := range ires.Refunds {
		if ires.Refunds[x] != i.Refunds[x] {
			t.Fatalf("invalid Refunds, expected record %v to be %d got %d",
				x, i.Refunds[x], ires.Refunds[x])
		}
	}

}
