package rsv

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sport/adyen"
	"sport/helpers"
	"sport/instr"
	"sport/sub"
	"sport/train"
	"sport/user"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
POST /api/rsv/pricing

DESC
	mock rsv request to validate if and obtain pricing info

REQUEST
	- [optional] Authorization: Bearer <token> in header
	- application/json body of type PostReservationRequest

RETURNS
	- status code:
		- 200 With RsvPricingInfo
		- error status code
			- 400 when validation of request fails
			- 500 on unexpected error
			- 410 if discount code is not longer valid
			- 401 if authorization was provided and system could not authorize user
			- 404 if target training was not found
*/
func (ctx *Ctx) HandlerPostRsvPricing() gin.HandlerFunc {
	return func(g *gin.Context) {

		var err error
		status := 400
		var apiReq ApiReservationRequest
		var tres *train.TrainingWithJoins
		var p RsvPricingInfo

		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body, &apiReq,
			func() error { return apiReq.Validate(ctx) },
		); err != nil {
			goto end
		}

		if tres, err = ctx.Train.DalReadSingleTraining(
			train.DalReadTrainingsRequest{
				TrainingID: &apiReq.TrainingID,
				WithOccs:   true,
				WithGroups: true,
				WithDcs:    apiReq.DcID != nil,
				DcID:       apiReq.DcID,
			},
		); err != nil {
			if err == sql.ErrNoRows {
				status = 404
			} else {
				status = 500
			}
			goto end
		}

		p = ctx.GetRsvPricing(tres)

		// validate that discount code is available
		if err = ctx.ValidatePricingDc(&p); err != nil {
			status = 410
			goto end
		}

		g.AbortWithStatusJSON(200, p)
		return

	end:
		g.AbortWithError(status, err)
	}
}

// const confirmReservationPath = "/api/reservation/confirm"

type PostRsvResponse struct {
	ID  uuid.UUID
	AT  uuid.UUID
	Url string
}

/*
POST /api/rsv

DESC
	begin reservation process

REQUEST
	- [optional] Authorization: Bearer <token> in header
	- application/json body of type PostReservationRequest

RETURNS
	- status code:
		- 303 redirect with Location header set to payment page
		- error status code
			- 409 on conflict (requested reservation when instructor is not available)
			- 400 when validation of request fails
			- 500 on unexpected error
			- 410 if discount code is not longer valid
			- 501 if payments are disabled on server
			- 401 if authorization was provided and system could not authorize user
			- 404 if target training was not found
*/
func (ctx *Ctx) HandlerPostRsv() gin.HandlerFunc {
	return func(g *gin.Context) {

		var err error
		status := 400
		var apiReq ApiReservationRequest
		var rsv *DDLRsv
		var rsvID, rsvAT uuid.UUID
		var dcTx *sql.Tx
		id := ""
		t := ""

		var rsvReq ReservationRequest

		now := NowMin()
		rsvReq.DateOfRsv = now

		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body,
			&apiReq,
			func() error {
				return apiReq.Validate(ctx)
			},
		); err != nil {
			goto end
		}

		apiReq.Occurrence = apiReq.Occurrence.Round(time.Minute)

		if rsvReq.Tr, err = ctx.Train.DalReadSingleTraining(
			train.DalReadTrainingsRequest{
				TrainingID: &apiReq.TrainingID,
				WithOccs:   true,
				WithGroups: true,
				WithDcs:    apiReq.DcID != nil,
				DcID:       apiReq.DcID,
			},
		); err != nil {
			if err == sql.ErrNoRows {
				status = 404
			} else {
				status = 500
			}
			goto end
		}

		rsvReq.AvInfo, err = ctx.CanRegister(rsvReq.Tr, apiReq.Occurrence, &now, EmptyBuffers, false)
		if err != nil {
			goto end
		}

		if !rsvReq.AvInfo.IsAvailable {
			status = 409
			if rsvReq.AvInfo.Reason != 0 {
				err = fmt.Errorf("Rsv is not available, reason: %d", rsvReq.AvInfo.Reason)
				goto end
			}
			err = fmt.Errorf("rsv is not available")
			goto end
		}

		if g.GetHeader("Authorization") != "" {

			var uid uuid.UUID

			uid, err = ctx.Api.AuthorizeUserFromCtx(g)
			if err != nil {
				status = 401
				goto end
			}

			rsvReq.UserID = &uid
		}

		if apiReq.UseSavedData {

			if rsvReq.UserID == nil {
				err = fmt.Errorf("user must be logged in to use saved data")
				status = 400
				goto end
			}

			var u *user.User

			if u, err = ctx.User.DalReadUser(
				*rsvReq.UserID, user.KeyTypeID, true,
			); err != nil {
				status = 401
				goto end
			}

			rsvReq.UserInfo = u.PubUserInfo
			rsvReq.UserContactData = u.ContactData
			rsvReq.UseUserAcc = true
		} else {
			rsvReq.UserInfo.UserData = apiReq.UserData
			rsvReq.UserContactData = apiReq.ContactData
		}

		if err = ValidateUserDataForRsv(&rsvReq.UserInfo); err != nil {
			status = 400
			goto end
		}
		// make sure that we support lang
		// if we dont then fallback to default
		rsvReq.UserInfo.Language =
			ctx.User.LangCtx.ApiLangOrDefault(rsvReq.UserInfo.Language)

		if err = ValidateContactDataForRsv(&rsvReq.UserContactData); err != nil {
			status = 400
			goto end
		}

		rsvID = uuid.New()
		rsvAT = uuid.New()
		rsvReq.PricingInfo = ctx.GetRsvPricing(rsvReq.Tr)

		// validate that discount code is available
		if err = ctx.ValidatePricingDc(&rsvReq.PricingInfo); err != nil {
			status = 410
			goto end
		}

		if rsvReq.PricingInfo.Dc != nil {
			if dcTx, err = ctx.Dal.Db.Begin(); err != nil {
				status = 500
				goto end
			}
			if err = ctx.Dc.DalUpdateDcUsesTx(rsvReq.PricingInfo.Dc.ID, 1, dcTx); err != nil {
				dcTx.Rollback()
				if err == sql.ErrNoRows {
					status = 410
				} else {
					status = 500
				}
				goto end
			}
		}

		if rsvReq.UserID != nil {
			id = rsvID.String()
			t = "id"
		} else {
			id = rsvAT.String()
			t = "token"
		}

		if ctx.NoPaymentFlow {
			rsvReq.AdyenRes = &adyen.CreatePaymentLinkResponse{
				Url: fmt.Sprintf(ctx.Config.RsvDetailsUrlFmt, id, t),
				ID:  uuid.New().String(),
			}
		} else {

			preq := adyen.RegisterTransactionRequest{
				Amount:      rsvReq.PricingInfo.TotalPrice,
				SessionID:   adyenWhPrefix + rsvID.String(),
				Currency:    rsvReq.Tr.Training.Currency,
				Description: rsvReq.Tr.Training.Title,
				Email:       rsvReq.UserContactData.Email,
				Country:     rsvReq.UserInfo.Country,
				Language:    rsvReq.UserInfo.Language,
				UrlReturn:   fmt.Sprintf(ctx.Config.RsvDetailsUrlFmt, id, t),
			}

			rsvReq.AdyenRes, err = ctx.Adyen.RegisterTransaction(&preq)
			if err != nil {
				if dcTx != nil {
					dcTx.Rollback()
				}
				status = 500
				goto end
			}
		}

		rsv = rsvReq.NewReservation(ctx, &rsvID, &rsvAT)

		if err = ctx.CreateRsv(rsv, nil); err != nil {
			if dcTx != nil {
				dcTx.Rollback()
			}
			status = 500
			goto end
		}

		if dcTx != nil {
			if err = dcTx.Commit(); err != nil {
				status = 500
				goto end
			}
		}

		ctx.EmailUserAboutLinkAsync(rsv)

		{
			res := PostRsvResponse{
				ID:  rsv.ID,
				AT:  rsv.AccessToken,
				Url: rsvReq.AdyenRes.Url,
			}
			if apiReq.NoRedirect {
				g.AbortWithStatusJSON(200, res)
			} else {
				j, err := json.Marshal(res)
				if err != nil {
					if dcTx != nil {
						dcTx.Rollback()
					}
					g.AbortWithError(500, err)
					return
				}
				g.Redirect(303, rsvReq.AdyenRes.Url)
				fmt.Fprint(g.Writer, string(j))
			}

		}

		return

	end:
		g.AbortWithError(status, err)
	}
}

type CancelRsvRequest struct {
	ReservationID uuid.UUID
	AccessToken   uuid.UUID
}

func (req *CancelRsvRequest) Validate() error {
	if req.ReservationID == uuid.Nil {
		return fmt.Errorf("validate CancelReservationRequest: invalid rsv id")
	}
	return nil
}

type DisputeRsvRequest struct {
	CancelRsvRequest
	Email string
	Msg   string
}

func (req *DisputeRsvRequest) Validate() error {

	if req.Email == "" || req.Msg == "" {
		return fmt.Errorf("invalid message")
	}

	if len(req.Email) > 64 || len(req.Msg) > 256 {
		return fmt.Errorf("invalid message")
	}

	return req.CancelRsvRequest.Validate()
}

/*
	this structure is parsed from query string
*/
type GetRsvQueryRequest struct {
	helpers.DateRange
	helpers.PaginationRequest
}

func (ctx *Ctx) GetRsvsHandlerContent(g *gin.Context) *RsvWithInstrPagination {

	//userID := ctx.MustGet("UserID").(uuid.UUID)
	var r GetRsvQueryRequest
	var err error
	status := 400
	var instructorID uuid.UUID
	var dbres *RsvWithInstrPagination
	var t string
	var userID uuid.UUID
	var id uuid.UUID
	var ids string
	var trainingIDs string
	var trainingID *uuid.UUID

	if err = r.PaginationRequest.FromQueryString(g, ctx.Config.MaxPageSize); err != nil {
		goto end
	}
	if r.DateRange, err = helpers.DateRangeFromQueryString(g); err != nil {
		goto end
	}

	t = g.Param("type")

	if t == "user" || t == "instructor" {
		userID, err = ctx.Api.AuthorizeUserFromCtx(g)
		if err != nil {
			status = 401
			goto end
		}
		ids = g.Query("id")
		if ids != "" {
			id, err = uuid.Parse(ids)
			if err != nil {
				status = 400
				goto end
			}
			if id == uuid.Nil {
				status = 400
				err = fmt.Errorf("nil id")
				goto end
			}
		}
	}

	if trainingIDs = g.Query("training_id"); trainingIDs != "" {
		if tid, err := uuid.Parse(trainingIDs); err != nil {
			goto end
		} else {
			trainingID = &tid
		}
	}

	switch t {
	case "user":
		if ids != "" {
			if dbres, err = ctx.ReadRsvWithInstr(
				ReadRsvsArgs{
					Pagination:     &r.PaginationRequest,
					DateRange:      r.DateRange,
					WithInstructor: true,
					TrainingID:     trainingID,
					UserID:         &userID,
					ID:             &id,
				},
			); err != nil {
				goto end
			}
		} else {
			if dbres, err = ctx.ReadRsvWithInstr(
				ReadRsvsArgs{
					Pagination:     &r.PaginationRequest,
					DateRange:      r.DateRange,
					WithInstructor: true,
					TrainingID:     trainingID,
					UserID:         &userID,
				},
			); err != nil {
				goto end
			}
		}
	case "instructor":
		if instructorID, err = ctx.Instr.DalReadInstructorID(userID); err != nil {
			status = 404
			goto end
		}
		if ids != "" {
			if dbres, err = ctx.ReadRsvWithInstr(
				ReadRsvsArgs{
					Pagination:     &r.PaginationRequest,
					DateRange:      r.DateRange,
					WithInstructor: true,
					TrainingID:     trainingID,
					ID:             &id,
					InstructorID:   &instructorID,
				},
			); err != nil {
				goto end
			}
		} else {
			if dbres, err = ctx.ReadRsvWithInstr(
				ReadRsvsArgs{
					Pagination:     &r.PaginationRequest,
					DateRange:      r.DateRange,
					WithInstructor: true,
					TrainingID:     trainingID,
					InstructorID:   &instructorID,
				},
			); err != nil {
				goto end
			}
		}
	case "token":
		ats := g.Query("access_token")
		if ats == "" {
			err = fmt.Errorf("no access_token specified")
			status = 404
			goto end
		}
		var at uuid.UUID
		at, err = uuid.Parse(ats)
		if err != nil {
			status = 404
			goto end
		}
		if at == uuid.Nil {
			status = 404
			err = fmt.Errorf("nil id")
			goto end
		}
		if dbres, err = ctx.ReadRsvWithInstr(
			ReadRsvsArgs{
				Pagination:     &r.PaginationRequest,
				DateRange:      r.DateRange,
				WithInstructor: true,
				TrainingID:     trainingID,
				AccessToken:    &at,
			},
		); err != nil {
			goto end
		}
	default:
		err = fmt.Errorf("invalid type: %s", g.Param("type"))
		goto end
	}

	return dbres

end:
	g.AbortWithError(status, err)
	return nil
}

/*
GET /api/rsv/t/:type

DESC
	get reservations for user (calendar, schedule)

REQUEST
	- Authorization: Bearer <token> in header
	- query string
		- fields for GetRsvQueryRequest (such as: start, end, page, size)
			containing pagination and filtering info
		- schedule:
			if provided with any value, then this endpoint will return data formatted
			in []TrainingSchedule format
		- training_id:
			if provided will filter by provided training
	- :type represents source of calendar.
		can be either:
			user (/api/rsv/t/user):
				UserID is treated as client
			instructor (/api/rsv/t/instructor):
				UserID is treated as instructor
			token (/api/rsv/t/token)
				UserID is treated as client
				but token id is required

RETURNS
	- on success application/json on type ReadReservationsResponse
	- status code:
		- 200 on success
		- error status code
			- 404 (only /api/rsv/t/instructor) if user is not an instructor
				and we cant fetch the calendar
			- 400 when validation of request params fails
			- 500 on unexpected error
			- 401 on authorization failure
*/
func (ctx *Ctx) HandlerGetRsv() gin.HandlerFunc {
	return func(g *gin.Context) {
		rsvs := ctx.GetRsvsHandlerContent(g)
		if rsvs != nil {
			g.AbortWithStatusJSON(200, rsvs)
		}
	}
}

type GetInstrContactRequest struct {
	InstructorID   uuid.UUID
	RsvAccessToken uuid.UUID
}

func ParseGetInstrContactRequest(g *gin.Context) (*GetInstrContactRequest, error) {

	s := g.Query("instructor_id")
	if s == "" {
		return nil, fmt.Errorf("parse GetInstrContactRequest: invalid instructor_id")
	}

	var err error
	res := new(GetInstrContactRequest)

	res.InstructorID, err = uuid.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("parse GetInstrContactRequest: %v", err)
	}

	s = g.Query("access_token")
	if s != "" {
		res.RsvAccessToken, err = uuid.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("parse GetInstrContactRequest: %v", err)
		}
	}

	return res, nil
}

func (ir *GetInstrContactRequest) Validate() error {
	if ir.InstructorID == uuid.Nil {
		return fmt.Errorf("validate GetInstrContactRequest: invalid instructor_id")
	}
	return nil
}

func (ctx *Ctx) HandlerGetInstrContact() gin.HandlerFunc {
	return func(g *gin.Context) {

		var err error
		req, err := ParseGetInstrContactRequest(g)
		if err != nil {
			g.AbortWithError(400, err)
			return
		}

		rr, err := ctx.Instr.DalReadInstructor(req.InstructorID, instr.InstructorID)
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		var isActive = true
		args := ReadRsvsArgs{
			IsActive:     &isActive,
			InstructorID: &req.InstructorID,
		}

		if req.RsvAccessToken == uuid.Nil {

			if g.GetHeader("Authorization") == "" {
				g.AbortWithError(401, fmt.Errorf("no auth provided"))
				return
			}

			var userID uuid.UUID
			userID, err = ctx.Api.AuthorizeUserFromCtx(g)
			if err != nil {
				g.AbortWithError(401, err)
				return
			}

			if rr.ContactData.Share {
				g.AbortWithStatusJSON(200, rr.ContactData)
				return
			}

			args.UserID = &userID
		} else {
			args.AccessToken = &req.RsvAccessToken
		}

		rs, err := ctx.ReadDDLRsvs(args)
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		if len(rs.Rsv) == 0 {

			for {
				if args.UserID == nil {
					break
				}

				if ctx.Sub == nil {
					break
				}

				// if no rsv found, try to find sub
				sb, err := ctx.Sub.DalReadSubWithJoins(sub.ReadSubRequest{
					InstructorID: &req.InstructorID,
					UserID:       args.UserID,
				})
				if err != nil {
					g.AbortWithError(500, err)
					return
				}

				if len(sb) != 0 {
					for i := range sb {
						if // ? sb[i].IsConfirmed &&
						sb[i].IsActive {
							g.AbortWithStatusJSON(200, rr.ContactData)
							return
						}
					}
				}
			}

			g.AbortWithError(401, fmt.Errorf("no active reservation found"))
			return
		}

		g.AbortWithStatusJSON(200, rr.ContactData)
	}
}
