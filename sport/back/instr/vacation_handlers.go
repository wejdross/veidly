package instr

import (
	"fmt"
	"sport/helpers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctx *Ctx) HandlerPostVacation() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var req VacationRequest
		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body, &req, req.Validate,
		); err != nil {
			g.AbortWithError(400, err)
			return
		}
		userID := g.MustGet("UserID").(uuid.UUID)
		iid, err := ctx.DalReadInstructorID(userID)
		if err != nil {
			if helpers.IsENF(err) {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}
		v := req.NewVacation(iid)
		if err := ctx.DalCreateVacation(v); err != nil {
			if helpers.PgIsUqViolation(err) {
				g.AbortWithError(409, err)
			} else {
				g.AbortWithError(500, err)
			}
		} else {
			g.AbortWithStatus(204)
		}
	}
}

func (ctx *Ctx) HandlerPatchVacation() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var req UpdateVacationRequest
		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body, &req, req.Validate,
		); err != nil {
			g.AbortWithError(400, err)
			return
		}
		userID := g.MustGet("UserID").(uuid.UUID)
		iid, err := ctx.DalReadInstructorID(userID)
		if err != nil {
			if helpers.IsENF(err) {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}
		if err := ctx.DalUpdateVacation(iid, &req); err != nil {
			if helpers.IsENF(err) {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
		} else {
			g.AbortWithStatus(204)
		}
	}
}

func (ctx *Ctx) HandlerDeleteVacation() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var req DeleteVacationRequest
		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body, &req, req.Validate,
		); err != nil {
			g.AbortWithError(400, err)
			return
		}
		userID := g.MustGet("UserID").(uuid.UUID)
		iid, err := ctx.DalReadInstructorID(userID)
		if err != nil {
			if helpers.IsENF(err) {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}
		if err := ctx.DalDeleteVacation(iid, &req); err != nil {
			if helpers.IsENF(err) {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
		} else {
			g.AbortWithStatus(204)
		}
	}
}

func (ctx *Ctx) HandlerGetVacations() gin.HandlerFunc {
	return func(g *gin.Context) {
		var iid uuid.UUID
		if g.GetHeader("Authorization") != "" {
			userID, err := ctx.Api.AuthorizeUserFromCtx(g)
			if err != nil {
				g.AbortWithError(401, err)
				return
			}
			iid, err = ctx.DalReadInstructorID(userID)
			if err != nil {
				if helpers.IsENF(err) {
					g.AbortWithError(404, err)
				} else {
					g.AbortWithError(500, err)
				}
				return
			}
		} else {
			iids := g.Query("instructor_id")
			if iids == "" {
				g.AbortWithError(400, fmt.Errorf("empty instructor_id"))
				return
			}
			var err error
			iid, err = uuid.Parse(iids)
			if err != nil {
				g.AbortWithError(400, err)
				return
			}
		}
		if res, err := ctx.DalReadVacations(iid, time.Now().In(time.UTC)); err != nil {
			g.AbortWithError(500, err)
		} else {
			g.AbortWithStatusJSON(200, res)
		}
	}
}
