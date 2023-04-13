package dc

import (
	"database/sql"
	"fmt"
	"sport/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
POST /api/dc

DESC
	create discount code

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 200 on success
		- error status code
			- 400 when validation of request fails
			- 404 if update target was not found
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerPostDc() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		var req DcRequest
		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &req, req.Validate); err != nil {
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
		dc := req.NewDc(iid)
		if err := ctx.DalCreateDc(nil, dc); err != nil {
			if e := helpers.HttpErr(err); e != nil {
				e.WriteAndAbort(g)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}
		g.AbortWithStatusJSON(200, struct{ ID uuid.UUID }{dc.ID})
	}
}

/*
PATCH /api/dc

DESC
	patch discount code

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
func (ctx *Ctx) HandlerPatchDc() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		var req UpdateDcRequest
		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}
		if err := ctx.DalUpdateDc(&req, userID); err != nil {
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
DELETE /api/dc

DESC
	delete discount code

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
func (ctx *Ctx) HandlerDeleteDc() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		var req DeleteDcRequest
		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}
		if err := ctx.DalDeleteDc(req.ID, userID); err != nil {
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
GET /api/dc

DESC
	get discount code for specified name

REQUEST
	- query string:
		name: name of the code
		training_id: id of the training

RETURNS
	- status code:
		- 200 on success
		- error status code
			- 404 if user is not instructor
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerRedeemDc() gin.HandlerFunc {
	return func(g *gin.Context) {

		name := g.Query("name")
		tidstr := g.Query("training_id")

		if name == "" {
			g.AbortWithError(500, fmt.Errorf("invalid name"))
			return
		}
		if tidstr == "" {
			g.AbortWithError(500, fmt.Errorf("invalid training_id"))
			return
		}

		tid, err := uuid.Parse(tidstr)
		if err != nil {
			g.AbortWithError(500, fmt.Errorf("invalid training_id"))
			return
		}

		res, err := ctx.DalReadDc(ReadDcRequest{
			Name:       name,
			TrainingID: &tid,
		}, nil)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		if len(res) != 1 {
			g.AbortWithError(404, fmt.Errorf("couldnt find dc"))
			return
		}

		g.AbortWithStatusJSON(200, res[0])
	}
}

/*
GET /api/dc/redeem

DESC
	find discount code

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 200 on success
		- error status code
			- 404 if user is not instructor
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerGetDc() gin.HandlerFunc {
	return func(g *gin.Context) {

		userID := g.MustGet("UserID").(uuid.UUID)
		iid, err := ctx.Instr.DalReadInstructorID(userID)
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		res, err := ctx.DalReadDc(ReadDcRequest{
			InstrID: &iid,
		}, nil)
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		g.AbortWithStatusJSON(200, res)
	}
}

/*
POST /api/dc/binding

DESC
	create discount code binding with training

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 200 on success
		- error status code
			- 400 when validation of request fails
			- 404 if update target was not found
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerPostDcBinding() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		var req DcBindingRequest
		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}
		if err := ctx.DalCreateDcBinding(req.TrainingID, req.DcID, userID); err != nil {
			if err == sql.ErrNoRows || helpers.PgIsFkViolation(err) {
				g.AbortWithError(404, err)
			} else if helpers.PgIsUqViolation(err) {
				g.AbortWithError(400, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}
		g.AbortWithStatus(204)
	}
}

/*
DELETE /api/dc/binding

DESC
	remove discount code binding with training

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 200 on success
		- error status code
			- 400 when validation of request fails
			- 404 if update target was not found
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerDeleteDcBinding() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		var req DcBindingRequest
		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}
		if err := ctx.DalDeleteDcBinding(req.TrainingID, req.DcID, userID); err != nil {
			if err == sql.ErrNoRows || helpers.PgIsFkViolation(err) {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}
		g.AbortWithStatus(204)
	}
}
