package instr_test

import (
	"encoding/json"
	"fmt"
	"sport/api"
	"sport/helpers"
	"sport/instr"
	"sport/rsv"
	"sport/user"
	"testing"
	"time"

	"github.com/google/uuid"
)

func validateVacationResponse(
	hdr string,
	got []instr.VacationInfo,
	expected []instr.VacationRequest,
	returnIDobj *instr.VacationRequest,
	returnID *uuid.UUID,
) error {
	if len(got) != len(expected) {
		return fmt.Errorf("%s: invalid vacation len, expected: %d, got %d", hdr, len(expected), len(got))
	}
	ridj := ""
	if returnIDobj != nil && returnID != nil {
		ridj = helpers.JsonMustSerializeStr(returnIDobj)
	}
	for i := range expected {
		foundOne := false
		for j := range got {
			gotj := helpers.JsonMustSerializeStr(got[j].VacationRequest)
			expj := helpers.JsonMustSerializeStr(expected[i])
			if got[j].ID == uuid.Nil {
				return fmt.Errorf("%s: invalid ID on the object: %s", hdr, helpers.JsonMustSerializeStr(got[j]))
			}
			if gotj == expj {
				if returnIDobj != nil && returnID != nil {
					if ridj == expj {
						*returnID = got[j].ID
					}
				}
				foundOne = true
				break
			}
		}
		if !foundOne {
			return fmt.Errorf("%s: couldnt find vacation record \n%s\n in response \n%s\n",
				hdr, helpers.JsonMustSerializeFormatStr(expected[i]),
				helpers.JsonMustSerializeFormatStr(got))
		}
	}
	return nil
}

func TestCrud(t *testing.T) {
	token, iid, err := instrCtx.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}
	token2, _, err := instrCtx.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}
	url := "/api/instructor/vacation"

	now := rsv.NowMin()

	req := instr.VacationRequest{
		DateStart: now,
		DateEnd:   now.Add(time.Hour * 24),
	}
	req2 := instr.VacationRequest{
		DateStart: now.Add(time.Hour * 24),
		DateEnd:   now.Add(2 * time.Hour * 24),
	}
	req3 := instr.VacationRequest{
		DateStart: now.Add(3 * time.Hour * 24),
		DateEnd:   now.Add(5 * time.Hour * 24),
	}

	method := "POST"

	tcs := []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 401,
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 400,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader:      helpers.JsonMustSerializeReader(instr.VacationRequest{}),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 400,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(instr.VacationRequest{
				DateStart: now,
			}),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 400,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(instr.VacationRequest{
				DateEnd: now,
			}),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedStatusCode: 400,
			RequestReader: helpers.JsonMustSerializeReader(instr.VacationRequest{
				DateEnd:   now,
				DateStart: now.Add(time.Hour),
			}),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedStatusCode: 204,
			RequestReader:      helpers.JsonMustSerializeReader(req),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 204,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader:      helpers.JsonMustSerializeReader(req2),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 204,
			RequestHeaders:     user.GetAuthHeader(token2),
			RequestReader:      helpers.JsonMustSerializeReader(req3),
		},
	}

	instrCtx.Api.TestAssertCases(t, tcs)

	var vid uuid.UUID

	method = "GET"
	tcs = []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 400,
		},
		{
			RequestMethod:      method,
			RequestUrl:         url + "?instructor_id=" + iid.String(),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var got []instr.VacationInfo
				if err := json.Unmarshal(b, &got); err != nil {
					return err
				}
				expected := []instr.VacationRequest{
					req,
					req2,
				}
				return validateVacationResponse("get0", got, expected, &req, &vid)
			},
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var got []instr.VacationInfo
				if err := json.Unmarshal(b, &got); err != nil {
					return err
				}
				expected := []instr.VacationRequest{
					req,
					req2,
				}
				return validateVacationResponse("get1", got, expected, &req, &vid)
			},
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token2),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var got []instr.VacationInfo
				if err := json.Unmarshal(b, &got); err != nil {
					return err
				}
				expected := []instr.VacationRequest{
					req3,
				}
				return validateVacationResponse("get1", got, expected, nil, nil)
			},
		},
	}

	instrCtx.Api.TestAssertCases(t, tcs)

	method = "PATCH"

	// modify 1st request
	patchReq := instr.UpdateVacationRequest{
		ID: vid,
		VacationRequest: instr.VacationRequest{
			DateStart: now.Add(time.Hour * 123),
			DateEnd:   now.Add(time.Hour * 231),
		},
	}

	tcs = []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 401,
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 400,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader:      helpers.JsonMustSerializeReader(instr.UpdateVacationRequest{}),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 400,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(instr.UpdateVacationRequest{
				VacationRequest: req2,
			}),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 400,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(instr.UpdateVacationRequest{
				ID: vid,
			}),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 404,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(instr.UpdateVacationRequest{
				ID:              uuid.New(),
				VacationRequest: patchReq.VacationRequest,
			}),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 404,
			RequestHeaders:     user.GetAuthHeader(token2),
			RequestReader:      helpers.JsonMustSerializeReader(patchReq),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 204,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader:      helpers.JsonMustSerializeReader(patchReq),
		},
		{
			RequestMethod:      "GET",
			RequestUrl:         url,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var got []instr.VacationInfo
				if err := json.Unmarshal(b, &got); err != nil {
					return err
				}
				expected := []instr.VacationRequest{
					patchReq.VacationRequest,
					req2,
				}
				return validateVacationResponse("get after patch", got, expected, nil, nil)
			},
		},
	}

	instrCtx.Api.TestAssertCases(t, tcs)

	method = "DELETE"

	tcs = []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 401,
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 400,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader:      helpers.JsonMustSerializeReader(instr.DeleteVacationRequest{}),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 404,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(instr.DeleteVacationRequest{
				ID: uuid.New(),
			}),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 404,
			RequestHeaders:     user.GetAuthHeader(token2),
			RequestReader: helpers.JsonMustSerializeReader(instr.DeleteVacationRequest{
				ID: vid,
			}),
		},
		{
			RequestMethod:      method,
			RequestUrl:         url,
			ExpectedStatusCode: 204,
			RequestHeaders:     user.GetAuthHeader(token),
			RequestReader: helpers.JsonMustSerializeReader(instr.DeleteVacationRequest{
				ID: vid,
			}),
		},
		{
			RequestMethod:      "GET",
			RequestUrl:         url,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var got []instr.VacationInfo
				if err := json.Unmarshal(b, &got); err != nil {
					return err
				}
				expected := []instr.VacationRequest{
					req2,
				}
				return validateVacationResponse("get after delete", got, expected, nil, nil)
			},
		},
	}

	instrCtx.Api.TestAssertCases(t, tcs)

}
