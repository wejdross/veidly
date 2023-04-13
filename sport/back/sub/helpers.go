package sub

import (
	"encoding/json"
	"fmt"
	"sport/adyen_sm"
	"sport/api"
	"sport/helpers"
	"sport/instr"
	"sport/user"
	"time"

	"github.com/google/uuid"
)

func NewTestSubModelRequest() SubModelRequest {
	v := helpers.CRNG_uint32Panic(100000)
	return SubModelRequest{
		Price: int((v%55)*100 + 1),
		SubModelEditable: SubModelEditable{
			Name:              helpers.CRNG_stringPanic(12),
			MaxEntrances:      int(v%100 + 1),
			Duration:          int(v%45 + 1),
			Currency:          "PLN",
			MaxActive:         int(v%76 + 1),
			IsFreeEntrance:    v%2 == 0,
			AllTrainingsByDef: v%2 == 1,
		},
	}
}

func (ctx *Ctx) ApiCreateSubModel(token string, sm *SubModelRequest) (uuid.UUID, error) {
	var id uuid.UUID
	err := ctx.Api.TestAssertCaseErr(&api.TestCase{
		RequestMethod:      "POST",
		RequestUrl:         "/api/sub/model",
		RequestReader:      helpers.JsonMustSerializeReader(sm),
		RequestHeaders:     user.GetAuthHeader(token),
		ExpectedStatusCode: 200,
		ExpectedBodyVal: func(b []byte, i interface{}) error {
			var res struct {
				ID uuid.UUID
			}
			if err := json.Unmarshal(b, &res); err != nil {
				return err
			}
			id = res.ID
			return nil
		},
	})
	return id, err
}

func (ctx *Ctx) ApiCreateSmBinding(token string, sm *SubModelBinding) error {
	err := ctx.Api.TestAssertCaseErr(&api.TestCase{
		RequestMethod:      "POST",
		RequestUrl:         "/api/sub/model/binding",
		RequestReader:      helpers.JsonMustSerializeReader(sm),
		RequestHeaders:     user.GetAuthHeader(token),
		ExpectedStatusCode: 204,
	})
	return err
}

func (ctx *Ctx) ApiRmSmBinding(token string, sm *SubModelBinding) error {
	err := ctx.Api.TestAssertCaseErr(&api.TestCase{
		RequestMethod:      "DELETE",
		RequestUrl:         "/api/sub/model/binding",
		RequestReader:      helpers.JsonMustSerializeReader(sm),
		RequestHeaders:     user.GetAuthHeader(token),
		ExpectedStatusCode: 204,
	})
	return err
}

func (ctx *Ctx) ApiCreateSub(token string, smID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := ctx.Api.TestAssertCaseErr(&api.TestCase{
		RequestMethod:      "POST",
		RequestUrl:         "/api/sub",
		ExpectedStatusCode: 200,
		RequestReader: helpers.JsonMustSerializeReader(ApiSubRequest{
			SubModelID: smID,
		}),
		RequestHeaders: user.GetAuthHeader(token),
		ExpectedBodyVal: func(b []byte, i interface{}) error {
			var res ApiSubResponse
			if err := json.Unmarshal(b, &res); err != nil {
				return err
			}
			if res.ID == uuid.Nil {
				return fmt.Errorf("invalid ID")
			}
			id = res.ID
			if res.Url == "" {
				return fmt.Errorf("invalid Url")
			}
			return nil
		},
	})
	return id, err
}

func (ctx *Ctx) ApiCreateAndQuerySub(token string, smID uuid.UUID) (*Sub, error) {
	sid, err := ctx.ApiCreateSub(token, smID)
	if err != nil {
		return nil, err
	}
	return ctx.DalReadSingleSub(ReadSubRequest{ID: &sid})
}

type SubTestOpts struct {
	// if != nil then will write user token into string @ provided address
	ReturnUserToken *string
	// if != nil then will write instr token into string @ provided address
	ReturnInstrToken *string

	UserIsInstr   bool
	UsersMayExist bool
}

var EmptySubTestOpts = SubTestOpts{}

// create basic objects in database used to test different rsv flows
func (ctx *Ctx) PrepSubTesting(
	opts SubTestOpts,
) (*Sub, error) {
	ur := &user.UserRequest{
		Email:    helpers.CRNG_stringPanic(10) + "@foo.bar",
		Password: "!@@#!sadASDASDAdasdas3452",
		UserData: user.UserData{
			Name:     "Velvet Velour",
			Language: "pl",
			Country:  "PL",
		},
	}
	instructorToken, err := ctx.User.ApiCreateAndLoginUser(ur)
	var userExists bool
	if err != nil {
		if opts.UsersMayExist {
			instructorToken, err = ctx.User.ApiLoginUser(ur)
			if err != nil {
				return nil, err
			}
			userExists = true
		} else {
			return nil, err
		}
	}

	if opts.ReturnInstrToken != nil {
		*opts.ReturnInstrToken = instructorToken
	}

	instructorUserID, err := ctx.Api.AuthorizeUserFromToken(instructorToken)
	if err != nil {
		return nil, err
	}

	err = ctx.Instr.ApiCreateInstructor(instructorToken, nil)
	if err != nil {
		if !userExists {
			return nil, err
		}
	}

	if err = ctx.Instr.DalUpdateCardInfo(instructorUserID, &instr.ProcessedCardInfo{
		CardRefID:      uuid.NewString(),
		CardBrand:      "test",
		CardHolderName: "test",
		CardSummary:    "1111",
	}); err != nil {
		return nil, err
	}

	var userToken string
	if opts.UserIsInstr {
		userToken = instructorToken
	} else {
		userToken, err = ctx.User.ApiCreateAndLoginUser(&user.UserRequest{
			Email:    helpers.CRNG_stringPanic(10) + "@foo.bar",
			Password: "!@@#!sadASDASDAdasdas3452",
			UserData: user.UserData{
				Name:     "na na na",
				Language: "pl",
				Country:  "PL",
			},
		})
		if err != nil {
			if opts.UsersMayExist {
				userToken, err = ctx.User.ApiLoginUser(ur)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	if opts.ReturnUserToken != nil {
		*opts.ReturnUserToken = userToken
	}

	smReq := NewTestSubModelRequest()
	smID, err := ctx.ApiCreateSubModel(instructorToken, &smReq)
	if err != nil {
		return nil, err
	}

	sid, err := ctx.ApiCreateSub(userToken, smID)
	if err != nil {
		return nil, err
	}

	sb, err := ctx.DalReadSingleSubWithJoins(ReadSubRequest{
		ID: &sid,
	})
	if err != nil {
		return nil, err
	}

	return &sb.Sub, nil
}

func EnsureSubIsInState(s *Sub, state adyen_sm.State) error {
	if s.State != state {
		return fmt.Errorf("invalid state, expected %s got %s", state, s.State)
	}
	return nil
}

type SubValidationOpts struct {
	State   adyen_sm.State
	Timeout time.Time
	// if true will reset sub timeout to now (after validation)
	ResetTimeout bool
	// no initial db request will be made if provided
	Sub *Sub
	//
	ExpectedCancelCount int
	//
	SmRetries      int
	ForceSmRetries bool
}

var EmptySubValidationOpts = SubValidationOpts{}

// ensure that timeout is right
func EnsureSubTimeoutInRange(
	s *Sub,
	to time.Time,
) error {
	df := s.SmTimeout.Sub(to)
	errorCorrection := time.Minute
	if df <= -errorCorrection || df >= errorCorrection {
		return fmt.Errorf("invalid sm_timeout, expected %v got: %v\ndf was %v", to, s.SmTimeout, df)
	}
	return nil
}

/*
	this is common function which may be used to validate rsv after (webhook / daemon / ... )
	action has been made

*/
func (ctx *Ctx) RefreshSubAndValidate(
	id uuid.UUID,
	vo SubValidationOpts,
) (*Sub, error) {
	var s *Sub
	var err error
	if vo.Sub == nil {
		if _s, err := ctx.DalReadSingleSubWithJoins(ReadSubRequest{
			ID: &id,
		}); err != nil {
			return nil, err
		} else {
			s = &_s.Sub
		}
	} else {
		s = vo.Sub
	}

	if vo.State != "" {
		if err := EnsureSubIsInState(s, vo.State); err != nil {
			return nil, err
		}
	}

	if !vo.Timeout.IsZero() {
		if err := EnsureSubTimeoutInRange(s, vo.Timeout); err != nil {
			return nil, err
		}
	}

	if vo.ExpectedCancelCount != 0 {
		if s.SmCache == nil {
			return nil, fmt.Errorf("null sm_cache")
		}
		if s.SmCache[adyen_sm.WaitCancelOrRefund].Retries != vo.ExpectedCancelCount {
			return nil, fmt.Errorf(
				"invalid cancel_or_refund retry count expected %d got %d",
				vo.ExpectedCancelCount,
				s.SmCache[adyen_sm.WaitCancelOrRefund].Retries)
		}
	}

	if vo.SmRetries != 0 || vo.ForceSmRetries {
		if s.SmRetries != vo.SmRetries {
			return nil, fmt.Errorf(
				"invalid sm_retries, expected %d got %d",
				vo.SmRetries,
				s.SmRetries)
		}
	}

	if vo.ResetTimeout {
		if err = ctx.UpdateSubSmTimeout(s.ID, helpers.NowMin()); err != nil {
			return nil, err
		}

		if _s, err := ctx.DalReadSingleSubWithJoins(ReadSubRequest{ID: &id}); err != nil {
			return nil, err
		} else {
			s = &_s.Sub
		}
	}

	return s, nil
}

func (ctx *Ctx) ApiCancelSub(
	subID uuid.UUID,
	userToken string,
	statusCode int,
) error {
	if statusCode == 0 {
		statusCode = 204
	}
	return ctx.Api.TestAssertCaseErr(&api.TestCase{
		RequestMethod:      "POST",
		RequestUrl:         "/api/sub/cancel",
		ExpectedStatusCode: statusCode,
		RequestHeaders:     api.JwtHeader(userToken),
		RequestReader: helpers.JsonMustSerializeReader(SmApiRequest{
			SubID: subID,
		}),
	})
}

func (ctx *Ctx) ApiExpireSub(
	subID uuid.UUID,
	userToken string,
	statusCode int,
) error {
	if statusCode == 0 {
		statusCode = 204
	}
	return ctx.Api.TestAssertCaseErr(&api.TestCase{
		RequestMethod:      "POST",
		RequestUrl:         "/api/sub/expire",
		ExpectedStatusCode: statusCode,
		RequestHeaders:     api.JwtHeader(userToken),
		RequestReader: helpers.JsonMustSerializeReader(SmApiRequest{
			SubID: subID,
		}),
	})
}

func (ctx *Ctx) ApiDisputeSub(
	subID uuid.UUID,
	userToken string,
	statusCode int,
	isInstructor bool,
) error {
	if statusCode == 0 {
		statusCode = 204
	}
	t := "user"
	if isInstructor {
		t = "instructor"
	}
	return ctx.Api.TestAssertCaseErr(&api.TestCase{
		RequestMethod:      "POST",
		RequestUrl:         "/api/sub/dispute/" + t,
		ExpectedStatusCode: statusCode,
		RequestHeaders:     api.JwtHeader(userToken),
		RequestReader: helpers.JsonMustSerializeReader(DisputeSubRequest{
			SmApiRequest: SmApiRequest{
				SubID: subID,
			},
			DisputeRequest: adyen_sm.DisputeRequest{
				Email: "foo@foo.foo",
				Msg:   "foo",
			},
		}),
	})
}

func (ctx *Ctx) ApiRefundSub(
	subID uuid.UUID,
	token string,
	statusCode int,
	isInstructor bool,
) error {
	if statusCode == 0 {
		statusCode = 204
	}
	var url = ""
	if isInstructor {
		url = "/api/sub/refund/instructor"
	} else {
		url = "/api/sub/refund/user"
	}
	return ctx.Api.TestAssertCaseErr(&api.TestCase{
		RequestMethod:      "POST",
		RequestUrl:         url,
		ExpectedStatusCode: statusCode,
		RequestHeaders:     api.JwtHeader(token),
		RequestReader: helpers.JsonMustSerializeReader(SmApiRequest{
			SubID: subID,
		}),
	})
}
