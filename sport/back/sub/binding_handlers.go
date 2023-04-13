package sub

import (
	"database/sql"
	"sport/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
POST /api/sub/model/binding

DESC
	create subscription model binding [to training]

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
func (ctx *Ctx) HandlerPostSubModelBinding() gin.HandlerFunc {
	return func(g *gin.Context) {

		userID := g.MustGet("UserID").(uuid.UUID)

		var req SubModelBinding
		if err := helpers.ReadJsonBodyFromReader(
			g.Request.Body, &req, req.Validate,
		); err != nil {
			g.AbortWithError(400, err)
			return
		}

		if err := ctx.DalCreateSubModelBinding(&req, userID); err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else if helpers.PgIsUqViolation(err) {
				g.AbortWithError(409, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		g.AbortWithStatus(204)
	}
}

/*
DELETE /api/sub/model/binding

DESC
	delete subscription model binding [to training]

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
func (ctx *Ctx) HandlerDeleteSubModelBinding() gin.HandlerFunc {
	return func(g *gin.Context) {

		userID := g.MustGet("UserID").(uuid.UUID)

		var req SubModelBinding
		if err := helpers.ReadJsonBodyFromReader(
			g.Request.Body, &req, req.Validate,
		); err != nil {
			g.AbortWithError(400, err)
			return
		}

		if err := ctx.DalDeleteSubModelBinding(&req, userID); err != nil {
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
