package instr

import (
	"database/sql"
	"fmt"
	"sport/helpers"
	"sport/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
GET /api/instructor

DESC
	get instructor information

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 200 on success
		- error status code
			- 404 if user is not an instructor
			- 500 on unexpected error
	- on success application/json body of type GetInstructorInfoResponse
*/
func (ctx *Ctx) HandlerGetInstructor() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		status := 500
		var key interface{}
		var keyType ReadInstructorInfoKeyType

		if g.GetHeader("Authorization") != "" {
			userID, err := ctx.Api.AuthorizeUserFromCtx(g)
			if err != nil {
				g.AbortWithError(401, err)
				return
			}
			key = userID
			keyType = UserID
		} else {
			instructorIDstr := g.Query("instructor_id")
			if instructorIDstr == "" {
				g.AbortWithError(404, fmt.Errorf("cant determine target instructor"))
				return
			}
			iid, err := uuid.Parse(instructorIDstr)
			if err != nil {
				g.AbortWithError(400, err)
				return
			}
			key = iid
			keyType = InstructorID
		}

		var ret *PubInstructorWithUser

		if ret, err = ctx.DalReadInstructorInfo(key, keyType); err != nil {
			if helpers.IsENF(err) {
				status = 404
			}
			goto end
		}

		g.AbortWithStatusJSON(200, ret)
		return
	end:
		g.AbortWithError(status, err)
	}
}

/*
POST /api/instructor

DESC
	create entry into instructor table

REQUEST
	- Authorization: Bearer <token> in header
	- application/json body of type InstructorRequest

RETURNS
	- status code:
		- 204 on success
		- error status code
			- 409 when user is already an instructor
			- 400 when validation of request fails
			- 500 on unexpected error
*/

func (ctx *Ctx) HandlerPostInstructor() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		status := 400
		var request InstructorRequest
		var instructor *Instructor
		userID := g.MustGet("UserID").(uuid.UUID)
		var u *user.User
		var di *DeletedInstructor

		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body, &request,
			func() error {
				return request.Validate(ctx)
			},
		); err != nil {
			goto end
		}

		instructor = request.NewInstructor(userID)

		u, err = ctx.User.DalReadUser(userID, user.KeyTypeID, true)
		if err != nil {
			if err == sql.ErrNoRows {
				status = 404
			} else {
				status = 500
			}
			goto end
		}

		if di, err = ctx.DalGetDeletedInstructor(u.Email); err != nil {
			if err != sql.ErrNoRows {
				status = 500
				goto end
			}
		} else {
			instructor.Refunds = di.Refunds
			instructor.QueuedPayoutCuts = di.QueuedPayoutCuts
		}

		if err = ctx.DalCreateInstructor(instructor); err != nil {
			if helpers.PgIsUqViolation(err) {
				status = 409
			} else {
				status = 500
			}
			goto end
		}

		g.AbortWithStatus(204)
		return

	end:
		g.AbortWithError(status, err)
	}
}

/*
DELETE /api/instructor

DESC
	delete record from instructor table along with associated data

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 204 on success
		- error status code
			- 404 if user is not an instructor
			- 400 when validation of request fails
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerDeleteInstructor() gin.HandlerFunc {
	return func(g *gin.Context) {

		var err error
		status := 500
		userID := g.MustGet("UserID").(uuid.UUID)
		var c bool
		var i *InstructorWithUser
		var tx *sql.Tx

		c, err = ctx.DalCanDeleteInstructor(userID)
		if err != nil {
			if helpers.IsENF(err) {
				status = 404
			} else {
				status = 500
			}
			goto end
		}

		/*
			there is a lot of DAL requests in here.
			should consider optimization...?
		*/

		i, err = ctx.DalReadInstructor(userID, UserID)
		if err != nil {
			goto end
		}

		if !c {
			status = 409
			err = fmt.Errorf("cant delete instructor since there are still active reservations")
			goto end
		}

		tx, err = ctx.Dal.Db.Begin()
		if err != nil {
			goto end
		}

		if len(i.Refunds) != 0 && i.QueuedPayoutCuts > 0 {
			var u *user.User
			u, err = ctx.User.DalReadUser(i.UserID, user.KeyTypeID, true)
			if err != nil {
				tx.Rollback()
				goto end
			}
			if err = ctx.DalInsertDeletedInstructor(&DeletedInstructor{
				Email:            u.Email,
				Refunds:          i.Refunds,
				QueuedPayoutCuts: i.QueuedPayoutCuts,
			}, tx); err != nil {
				tx.Rollback()
				goto end
			}
		}

		if err = ctx.DalDeleteInstructor(userID, tx); err != nil {
			goto end
		}

		if err := tx.Commit(); err != nil {
			goto end
		}

		g.AbortWithStatus(204)
		return

	end:
		g.AbortWithError(status, err)
	}
}

/*
GET /api/instructor/can_delete

DESC
	determine if instructor can be removed

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 204 on success (can delete)
		- 409 conflict (cannot delete)
		- error status code
			- 404 if user is not an instructor
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerCanDeleteInstructor() gin.HandlerFunc {
	return func(g *gin.Context) {

		status := 500
		userID := g.MustGet("UserID").(uuid.UUID)
		if c, err := ctx.DalCanDeleteInstructor(userID); err != nil {
			if helpers.IsENF(err) {
				status = 404
			} else {
				status = 500
			}
			g.AbortWithError(status, err)
		} else {
			if c {
				g.AbortWithStatus(204)
			} else {
				g.AbortWithStatus(409)
			}
		}
	}
}

/*
PATCH /api/instructor

DESC
	update instructor data
	<<<note that this is full merge of the structure into table>>>

REQUEST
	- Authorization: Bearer <token> in header
	- application/json body of type InstructorRequest

RETURNS
	- status code:
		- 204 on success
		- error status code
			- 404 when user is not an instructor
			- 400 when validation of request fails
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerPatchInstructor() gin.HandlerFunc {
	return func(g *gin.Context) {

		var err error
		status := 400
		var request InstructorRequest
		userID := g.MustGet("UserID").(uuid.UUID)

		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body, &request,
			func() error {
				return request.Validate(ctx)
			},
		); err != nil {
			g.AbortWithError(status, err)
			return
		}

		if err = ctx.DalPatchInstructor(userID, &request); err != nil {
			if helpers.IsENF(err) {
				status = 404
			} else {
				status = 500
			}
			g.AbortWithError(status, err)
			return
		}

		g.AbortWithStatus(204)
	}
}

func (ctx *Ctx) HandlerPatchInstructorPayments() gin.HandlerFunc {
	return func(g *gin.Context) {

		var err error
		status := 400
		var request CardInfo
		userID := g.MustGet("UserID").(uuid.UUID)
		var pci *ProcessedCardInfo

		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body, &request, nil,
		); err != nil {
			goto end
		}

		if pci, err = ctx.ProcessCardInfo(&request); err != nil {
			goto end
		}

		if err = ctx.DalUpdateCardInfo(userID, pci); err != nil {
			if err != nil {
				status = 403
			}
			goto end
		}

		g.AbortWithStatus(204)
		return

	end:
		g.AbortWithError(status, err)
	}
}

func (ctx *Ctx) HandlerGetInstructorPayments() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		status := 400
		var data *ProcessedCardInfo
		userID := g.MustGet("UserID").(uuid.UUID)
		iid, err := ctx.DalReadInstructorID(userID)
		if err != nil {
			if err == sql.ErrNoRows {
				status = 404
			} else {
				status = 500
			}
		}
		if data, err = ctx.DalReadCardInfo(iid); err != nil {
			if err == sql.ErrNoRows {
				status = 404
			} else {
				status = 500
			}
			goto end
		}

		g.AbortWithStatusJSON(200, data)

		return

	end:
		g.AbortWithError(status, err)
	}
}

func (ctx *Ctx) HandlerDeletePayoutInfo() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		status := 400

		userID := g.MustGet("UserID").(uuid.UUID)

		if err = ctx.DalDeletePayoutInfo(userID); err != nil {
			if helpers.IsENF(err) {
				status = 404
			} else {
				status = 500
			}
			goto end
		}

		g.AbortWithStatus(204)
		return

	end:
		g.AbortWithError(status, err)
	}
}

type InfoResponse struct {
	ConfiguredState ConfigState
}

// this is general information endpont which we may use to communicate with instructor
func (ctx *Ctx) HandlerInstructorInfo() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		status := 400
		userID := g.MustGet("UserID").(uuid.UUID)
		var res InfoResponse

		if res.ConfiguredState, err = ctx.GetInstrConfigByID(userID); err != nil {
			if err == sql.ErrNoRows {
				status = 404
			} else {
				status = 500
			}
			goto end
		}

		g.AbortWithStatusJSON(200, res)

		return

	end:
		g.AbortWithError(status, err)
	}
}
