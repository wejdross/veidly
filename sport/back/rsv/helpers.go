package rsv

import (
	"fmt"
	"sport/adyen_sm"
	"sport/api"
	"sport/helpers"
	"sport/instr"
	"sport/train"
	"sport/user"
	"time"

	"github.com/google/uuid"
)

func overlaps(start1, end1, start2, end2 time.Time) bool {
	return (start1.Before(end2)) && (end1.After(start2))
}

func NowMin() time.Time {
	now := time.Now().In(time.UTC)
	now = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		0, 0, time.UTC)
	return now
}

// returns ID of created reservation, and path to the payment link
func (ctx *Ctx) ApiGetRsvPricing(
	cr *ApiReservationRequest,
	expectedStatus int,
) (RsvPricingInfo, error) {
	if expectedStatus == 0 {
		expectedStatus = 200
	}
	var res RsvPricingInfo
	cases := []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/rsv/pricing",
			RequestReader:      helpers.JsonMustSerializeReader(cr),
			ExpectedStatusCode: expectedStatus,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				// if error requested dont validate
				if expectedStatus >= 400 {
					return nil
				}
				helpers.JsonMustDeserialize(b, &res)
				return nil
			},
		},
	}
	err := ctx.Api.TestAssertCasesErr(cases)
	return res, err
}

// returns ID of created reservation, and path to the payment link
func (ctx *Ctx) ApiCreateRsv(
	token *string,
	cr *ApiReservationRequest,
	expectedStatus int,
) (uuid.UUID, string, error) {
	var rh map[string]string = nil
	if token != nil {
		rh = api.JwtHeader(*token)
	}
	var ID uuid.UUID
	var location string
	if expectedStatus == 0 {
		expectedStatus = 303
	}
	cases := []api.TestCase{
		{
			RequestMethod:      "POST",
			RequestUrl:         "/api/rsv",
			RequestReader:      helpers.JsonMustSerializeReader(cr),
			ExpectedStatusCode: expectedStatus,
			RequestHeaders:     rh,
			ExpectedHeadersVal: func(m map[string][]string) error {
				// if error requested dont validate
				if expectedStatus/100 != 3 {
					return nil
				}
				c := m["Location"]
				if len(c) != 1 {
					return fmt.Errorf("invalid location hdr in response")
				}
				if c[0] == "" {
					return fmt.Errorf("location was empty")
				}
				location = c[0]
				return nil
			},
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				// if error requested dont validate
				if expectedStatus >= 400 {
					return nil
				}
				var res PostRsvResponse
				helpers.JsonMustDeserialize(b, &res)
				if res.ID == uuid.Nil {
					return fmt.Errorf("empty id returned")
				}
				ID = res.ID
				return nil
			},
		},
	}
	err := ctx.Api.TestAssertCasesErr(cases)
	return ID, location, err
}

func NewTestCreateReservationRequest(trainingID uuid.UUID, start time.Time) ApiReservationRequest {
	return ApiReservationRequest{
		TrainingID: trainingID,
		Occurrence: start,
		UserData: user.UserData{
			Name: "Jeanette",
			//LastName:  "Voerman",
			Language: "pl",
			Country:  "PL",
		},
	}
}

func (ctx *Ctx) ApiCreateAndQueryRsv(
	token *string,
	cr *ApiReservationRequest,
	status int,
) (*DDLRsvWithInstr, error) {
	id, _, err := ctx.ApiCreateRsv(token, cr, status)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, nil
	}
	rsv, err := ctx.ReadSingleRsv(ReadRsvsArgs{
		ID:             &id,
		WithInstructor: true,
	})
	if err != nil {
		return nil, err
	}
	return rsv, nil
}

type RsvTestOpts struct {
	RsvStartDateSinceNow time.Duration
	DisableLinkExpress   bool

	// if != nil then will write user token into string @ provided address
	ReturnUserToken *string
	// if != nil then will write instr token into string @ provided address
	ReturnInstrToken *string

	//
	AltInstructorUserRq *user.UserRequest
	AltTraineeUserRq    *user.UserRequest
	UserIsInstr         bool
	UsersMayExist       bool
	// note that trainingID and Occurrence will always be overriden
	AltRsvRequest *ApiReservationRequest
	AnonRsvUser   bool

	RsvPrice int
}

var EmptyRvsTestOpts = RsvTestOpts{}

// create basic objects in database used to test different rsv flows
func (ctx *Ctx) PrepRsvTesting(
	opts RsvTestOpts,
) (*DDLRsvWithInstr, error) {
	var rr *user.UserRequest
	if opts.AltInstructorUserRq != nil {
		rr = opts.AltInstructorUserRq
	} else {
		rr = &user.UserRequest{
			Email:    helpers.CRNG_stringPanic(10) + "@foo.bar",
			Password: "!@#!@3zasdas45345DFSDFSD#@!$",
			UserData: user.UserData{
				Name:     "Velvet Velour",
				Language: "pl",
				Country:  "PL",
			},
		}
	}
	instructorToken, err := ctx.User.ApiCreateAndLoginUser(rr)
	var userExists bool
	if err != nil {
		if opts.UsersMayExist {
			instructorToken, err = ctx.User.ApiLoginUser(rr)
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

	if opts.AltTraineeUserRq != nil {
		rr = opts.AltTraineeUserRq
	} else {
		rr = nil
	}

	var userToken string
	if opts.UserIsInstr {
		userToken = instructorToken
	} else {
		userToken, err = ctx.User.ApiCreateAndLoginUser(rr)
		if err != nil {
			if opts.UsersMayExist {
				userToken, err = ctx.User.ApiLoginUser(rr)
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

	now := NowMin()

	start := now.Add(opts.RsvStartDateSinceNow)
	end := start.Add(time.Hour)

	price := 10 * 100
	if opts.RsvPrice != 0 {
		price = opts.RsvPrice
	}

	trainingID, err := ctx.Train.ApiCreateTraining(instructorToken, &train.CreateTrainingRequest{
		Training: train.TrainingRequest{
			Title:         "dancing in vesuvius",
			Capacity:      1,
			Currency:      "PLN",
			Price:         price,
			ManualConfirm: false,
			AllowExpress:  !opts.DisableLinkExpress,
		},
		Occurrences: []train.CreateOccRequest{
			{
				OccRequest: train.OccRequest{
					DateStart:  start,
					DateEnd:    end,
					RepeatDays: 1,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var ud user.UserData
	var cd user.ContactData
	if opts.AltRsvRequest != nil {
		ud = opts.AltRsvRequest.UserData
		cd = opts.AltRsvRequest.ContactData
	} else {
		ud = user.UserData{
			Name:     "foo bar",
			Language: "pl",
			Country:  "PL",
			AboutMe:  "i like cookies",
		}
	}

	crr := &ApiReservationRequest{
		TrainingID:  trainingID,
		Occurrence:  start,
		UserData:    ud,
		ContactData: cd,
	}

	if opts.AnonRsvUser {
		return ctx.ApiCreateAndQueryRsv(nil, crr, 0)
	} else {
		return ctx.ApiCreateAndQueryRsv(&userToken, crr, 0)
	}
}

// ensure that timeout is right
func EnsureRsvTimeoutInRange(
	rsv *DDLRsv,
	to time.Time,
) error {
	df := rsv.SmTimeout.Sub(to)
	errorCorrection := time.Minute * 10
	if df <= -errorCorrection || df >= errorCorrection {
		return fmt.Errorf("invalid sm_timeout, expected %v got: %v\ndf was %v", to, rsv.SmTimeout, df)
	}
	return nil
}

func EnsureRsvIsInState(rsv *DDLRsvWithInstr, state adyen_sm.State) error {
	if rsv.State != state {
		return fmt.Errorf("invalid state, expected %s got %s", state, rsv.State)
	}
	return nil
}

type RsvValidationOpts struct {
	State   adyen_sm.State
	Timeout time.Time
	// if true will reset rsv timeout to now (after validation)
	ResetTimeout bool
	// no initial db request will be made if provided
	Rsv *DDLRsvWithInstr
	//
	ExpectedCancelCount int
	//
	SmRetries      int
	ForceSmRetries bool
	//
	ResetDateStart time.Time
}

/*
	this is common function which may be used to validate rsv after (webhook / daemon / ... )
	action has been made

*/
func (ctx *Ctx) RefreshRsvAndValidate(
	id uuid.UUID,
	vo RsvValidationOpts,
) (*DDLRsvWithInstr, error) {
	var rr *DDLRsvWithInstr
	var err error
	if vo.Rsv == nil {
		if rr, err = ctx.ReadRsvByID(id); err != nil {
			return nil, err
		}
	} else {
		rr = vo.Rsv
	}

	r := &rr.DDLRsv

	if vo.State != "" {
		if err := EnsureRsvIsInState(rr, vo.State); err != nil {
			return nil, err
		}
	}

	if !vo.Timeout.IsZero() {
		if err := EnsureRsvTimeoutInRange(
			r, vo.Timeout,
		); err != nil {
			return nil, err
		}
	}

	if vo.ExpectedCancelCount != 0 {
		if r.SmCache == nil {
			return nil, fmt.Errorf("null sm_cache")
		}
		if r.SmCache[adyen_sm.WaitCancelOrRefund].Retries != vo.ExpectedCancelCount {
			return nil, fmt.Errorf(
				"invalid cancel_or_refund retry count expected %d got %d",
				vo.ExpectedCancelCount,
				r.SmCache[adyen_sm.WaitCancelOrRefund].Retries)
		}
	}

	if vo.SmRetries != 0 || vo.ForceSmRetries {
		if r.SmRetries != vo.SmRetries {
			return nil, fmt.Errorf(
				"invalid sm_retries, expected %d got %d",
				vo.SmRetries,
				r.SmRetries)
		}
	}

	if !vo.ResetDateStart.IsZero() {
		len := r.DateEnd.Sub(r.DateStart)
		end := vo.ResetDateStart.Add(len)
		if err = ctx.UpdateRsvSmDateStart(r.ID, vo.ResetDateStart, end); err != nil {
			return nil, err
		}
	}

	if vo.ResetTimeout {
		if err = ctx.UpdateRsvSmTimeout(r.ID, NowMin()); err != nil {
			return nil, err
		}

		if rr, err = ctx.ReadRsvByID(id); err != nil {
			return nil, err
		}
		// r = &rr.DDLRsv
	}

	return rr, nil
}

func (ctx *Ctx) ApiCancelRsv(
	rsvID uuid.UUID,
	userToken string,
	statusCode int,
) error {
	if statusCode == 0 {
		statusCode = 204
	}
	return ctx.Api.TestAssertCaseErr(&api.TestCase{
		RequestMethod:      "POST",
		RequestUrl:         "/api/rsv/cancel",
		ExpectedStatusCode: statusCode,
		RequestHeaders:     api.JwtHeader(userToken),
		RequestReader: helpers.JsonMustSerializeReader(CancelRsvRequest{
			ReservationID: rsvID,
		}),
	})
}

func (ctx *Ctx) ApiExpireRsv(
	rsvID uuid.UUID,
	userToken string,
	statusCode int,
) error {
	if statusCode == 0 {
		statusCode = 204
	}
	return ctx.Api.TestAssertCaseErr(&api.TestCase{
		RequestMethod:      "POST",
		RequestUrl:         "/api/rsv/expire",
		ExpectedStatusCode: statusCode,
		RequestHeaders:     api.JwtHeader(userToken),
		RequestReader: helpers.JsonMustSerializeReader(CancelRsvRequest{
			ReservationID: rsvID,
		}),
	})
}

func (ctx *Ctx) ApiDisputeRsv(
	rsvID uuid.UUID,
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
		RequestUrl:         "/api/rsv/dispute/" + t,
		ExpectedStatusCode: statusCode,
		RequestHeaders:     api.JwtHeader(userToken),
		RequestReader: helpers.JsonMustSerializeReader(DisputeRsvRequest{
			CancelRsvRequest: CancelRsvRequest{
				ReservationID: rsvID,
			},
			Email: "foo@foo.foo",
			Msg:   "foo",
		}),
	})
}

func (ctx *Ctx) ApiRefundRsv(
	rsvID uuid.UUID,
	token string,
	statusCode int,
	isInstructor bool,
) error {
	if statusCode == 0 {
		statusCode = 204
	}
	var url = ""
	if isInstructor {
		url = "/api/rsv/refund/instructor"
	} else {
		url = "/api/rsv/refund/user"
	}
	return ctx.Api.TestAssertCaseErr(&api.TestCase{
		RequestMethod:      "POST",
		RequestUrl:         url,
		ExpectedStatusCode: statusCode,
		RequestHeaders:     api.JwtHeader(token),
		RequestReader: helpers.JsonMustSerializeReader(CancelRsvRequest{
			ReservationID: rsvID,
		}),
	})
}

func (ctx *Ctx) AssertRsvResponse(
	res *RsvWithInstrPagination,
	req []*ApiReservationRequest,
	userID *uuid.UUID,
	npages int,
	page int,
	size int,
) error {

	var err error

	if err = helpers.AssertErr("pagination:NPages", res.Pagination.NPages, npages); err != nil {
		return err
	}
	if err = helpers.AssertErr("pagination:Page", res.Pagination.Page, page); err != nil {
		return err
	}
	if err = helpers.AssertErr("pagination:Size", res.Pagination.Size, size); err != nil {
		return err
	}
	if err = helpers.AssertErr("len:Rsv", len(req), len(res.Rsv)); err != nil {
		return err
	}

	for i := 0; i < len(req); i++ {
		if userID == nil {
			return fmt.Errorf("test not yet implemented")
		} else {
			var u *user.User
			if u, err = ctx.User.DalReadUser(userID, user.KeyTypeID, true); err != nil {
				return err
			}
			if err = helpers.AssertJsonErr(
				"invalid user data",
				u.UserData,
				res.Rsv[i].UserInfo.UserData,
			); err != nil {
				return err
			}
		}
		if err = helpers.AssertErr(
			"invalid date start",
			req[i].Occurrence,
			res.Rsv[i].DateStart,
		); err != nil {
			return err
		}
		if err = helpers.AssertErr(
			"invalid trainingID",
			req[i].TrainingID,
			res.Rsv[i].Training.ID,
		); err != nil {
			return err
		}
	}
	return nil
}
