package instr

import (
	"sport/api"
	"sport/helpers"

	"github.com/google/uuid"
)

func (ctx *Ctx) ApiEditPayoutInfo(token string, ci *CardInfo) error {
	tc := []api.TestCase{
		{
			RequestMethod:      "PATCH",
			RequestUrl:         "/api/instructor/payout",
			RequestReader:      helpers.JsonMustSerializeReader(ci),
			ExpectedStatusCode: 204,
			RequestHeaders:     api.JwtHeader(token),
		},
	}
	return ctx.Api.TestAssertCasesErr(tc)
}

func (ctx *Ctx) CreateTestInstructor(userID uuid.UUID, ir *InstructorRequest) (uuid.UUID, error) {
	if ir == nil {
		ir = &InstructorRequest{
			/* creating empty array instead of nil to prevent object validation from complaining */
			ProfileSections: make(ProfileSectionArr, 0),
		}
	}
	i := ir.NewInstructor(userID)
	return i.ID, ctx.DalCreateInstructor(i)
}

// returns token, instructor_id, error
func (ctx *Ctx) CreateTestInstructorWithUser() (string, uuid.UUID, error) {
	token, err := ctx.User.ApiCreateAndLoginUser(nil)
	if err != nil {
		return "", uuid.Nil, err
	}
	uid, err := ctx.Api.AuthorizeUserFromToken(token)
	if err != nil {
		return "", uuid.Nil, err
	}
	iid, err := ctx.CreateTestInstructor(uid, nil)
	if err != nil {
		return "", uuid.Nil, err
	}
	return token, iid, nil
}

/*
	FOLLOWING FUNCTIONS ARE MEANT TO BE USED IN TESTING ONLY
	those functions will call api endpoints to create resources
	thus they require modules to be initialized
*/

// create instructor via api request
func (ctx *Ctx) ApiCreateInstructor(token string, ir *InstructorRequest) error {
	if ir == nil {
		ir = &InstructorRequest{}
	}
	tc := []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/instructor",
			RequestReader:      helpers.JsonMustSerializeReader(ir),
			ExpectedStatusCode: 204,
			RequestHeaders:     api.JwtHeader(token),
		},
	}
	return ctx.Api.TestAssertCasesErr(tc)
}
