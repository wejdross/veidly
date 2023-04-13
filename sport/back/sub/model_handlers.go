package sub

import (
	"database/sql"
	"fmt"
	"sport/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
POST /api/sub/model

DESC
	create subscription model

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 204 on success
		- error status code
			- 400 when validation of request fails
			- 404 if instructor was not found
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerPostSubModel() gin.HandlerFunc {
	return func(g *gin.Context) {

		userID := g.MustGet("UserID").(uuid.UUID)

		var req SubModelRequest
		if err := helpers.ReadJsonBodyFromReader(
			g.Request.Body, &req, func() error {
				return req.Validate(ctx)
			},
		); err != nil {
			g.AbortWithError(400, err)
			return
		}

		iid, err := ctx.Instr.DalReadInstructorID(userID)
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		sm := req.NewSubModel(ctx, userID, iid)
		if err := ctx.DalCreateSubModel(&sm); err != nil {
			if e := helpers.HttpErr(err); e != nil {
				e.WriteAndAbort(g)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}
		g.AbortWithStatusJSON(200, struct{ ID uuid.UUID }{sm.ID})
	}
}

/*
PATCH /api/sub/model

DESC
	patch subscription model

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 204 on success
		- error status code
			- 400 when validation of request fails
			- 404 if update target was not found
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerPatchSubModel() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)

		var req UpdateSubModelRequest
		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &req, func() error {
			return req.Validate(ctx)
		}); err != nil {
			g.AbortWithError(400, err)
			return
		}

		if err := ctx.DalUpdateSubModel(&req, userID); err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		g.AbortWithStatus(204)
	}
}

/*
DELETE /api/sub/model

DESC
	delete subscription model

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 204 on success
		- error status code
			- 400 when validation of request fails
			- 404 if update target was not found
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerDeleteSubModel() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		var req DeleteSubModelRequest
		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}
		if err := ctx.DalDeleteSubModel(req.ID, userID); err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}
		g.AbortWithStatus(204)
	}
}

/*
GET /api/sub/model

DESC
	get sub models for instructor

RETURNS
	- status code:
		- 200 on success
		- error status code
			- 404 if user is not instructor
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerGetSubModel() gin.HandlerFunc {
	return func(g *gin.Context) {

		var res []SubModel
		var err error

		if g.GetHeader("Authorization") != "" {
			var userID uuid.UUID
			userID, err = ctx.Api.AuthorizeUserFromCtx(g)
			if err != nil {
				g.AbortWithError(401, err)
				return
			}
			res, err = ctx.DalReadSubModels(ReadSubModelRequest{
				InstrUserID: &userID,
			})
		} else {
			instrIDstr := g.Query("instructor_id")
			if instrIDstr == "" {
				g.AbortWithError(400, fmt.Errorf("invalid instructor_id"))
				return
			}
			iid, err := uuid.Parse(instrIDstr)
			if err != nil {
				g.AbortWithError(400, err)
				return
			}
			res, err = ctx.DalReadSubModels(ReadSubModelRequest{
				InstructorID: &iid,
			})
		}

		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		g.AbortWithStatusJSON(200, res)
	}
}
