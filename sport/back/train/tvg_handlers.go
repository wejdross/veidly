package train

import (
	"database/sql"
	"sport/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
GET /api/training/group

DESC
	get training groups

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 200 on success
		- error status code
			- 500 on unexpected error
	- on success application/json body of type []TrainingGroup
*/
func (ctx *Ctx) HandlerGetTrainingGroups() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		res, err := ctx.DalReadTrainingGroups(ReadTrainingGroupsRequest{UserID: &userID})
		if err != nil {
			g.AbortWithError(500, err)
			return
		}
		g.AbortWithStatusJSON(200, res)
	}
}

/*
PATCH /api/training/group

DESC
	update training group
	note that this will result in full structure merge

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
func (ctx *Ctx) HandlerPatchTrainingGroup() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		var req UpdateGroupRequest
		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}
		if err := ctx.DalUpdateTrainingGroup(userID, &req); err != nil {
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
DELETE /api/training/group

DESC
	delete training group

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
func (ctx *Ctx) HandlerDeleteTrainingGroup() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		var req DeleteGroupRequest
		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}
		if err := ctx.DalDeleteTrainingGroup(userID, &req); err != nil {
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
POST /api/training/group

DESC
	create training group

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
func (ctx *Ctx) HandlerPostTrainingGroup() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		var req GroupRequest
		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}
		grp := req.NewGroup(userID)
		if err := ctx.DalCreateTrainingGroup(&grp); err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}
		g.AbortWithStatusJSON(200, struct{ ID uuid.UUID }{grp.ID})
	}
}

/*
PUT /api/training/group/binding

DESC
	bind training to specified group
	can be reversed by calling DELETE

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 204 on success
		- error status code
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerPutGroupBinding() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		var req GroupBindingRequest
		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}
		err := ctx.DalAddTvgBinding(req.TrainingID, req.GroupID, userID)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}
		g.AbortWithStatus(204)
	}
}

/*
DELETE /api/training/group/binding

DESC
	REMOVE binding between training and specified group

REQUEST
	- Authorization: Bearer <token> in header

RETURNS
	- status code:
		- 204 on success
		- error status code
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerDeleteGroupBinding() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		var req GroupBindingRequest
		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}
		err := ctx.DalRemoveTvgBinding(req.TrainingID, req.GroupID, userID)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}
		g.AbortWithStatus(204)
	}
}
