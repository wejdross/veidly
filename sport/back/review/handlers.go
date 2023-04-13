package review

import (
	"database/sql"
	"fmt"
	"sport/helpers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctx *Ctx) HandlerPostReview() gin.HandlerFunc {
	return func(g *gin.Context) {
		var req UpdateReviewRequest
		var err error

		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body,
			&req,
			req.Validate,
		); err != nil {
			g.AbortWithError(400, err)
			return
		}

		if err = ctx.DalUpdateReview(&req); err != nil {
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

func (ctx *Ctx) HandlerDeleteReview() gin.HandlerFunc {
	return func(g *gin.Context) {
		var req DeleteReviewRequest
		var err error

		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body,
			&req,
			req.Validate,
		); err != nil {
			g.AbortWithError(400, err)
			return
		}

		userID := g.MustGet("UserID").(uuid.UUID)

		if err = ctx.DalDeleteReview(userID, &req); err != nil {
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

func (ctx *Ctx) HandlerGetUserReview() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		rsvIDs := g.Query("rsv_id")
		if rsvIDs == "" {
			g.AbortWithError(400, fmt.Errorf("invalid rsv_id"))
			return
		}
		rsvID, err := uuid.Parse(rsvIDs)
		if err != nil {
			g.AbortWithError(400, fmt.Errorf("parse rsv_id: %v", err))
			return
		}
		if rsvID == uuid.Nil {
			g.AbortWithError(400, fmt.Errorf("empty rsv_id"))
			return
		}
		rv, err := ctx.DalReadSingleReview(SingleRevKeyUserRsvKey{
			UserID: userID,
			RsvID:  rsvID,
		}, SingleRevKeyUserRsv, AnyReviewType, nil)
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		var res ReviewResponse

		if rv.AccessToken != "" {
			res.Type = TokenReviewType
			var rtr ReviewTokenResponse
			rtr.AccessToken = rv.AccessToken
			rtr.ExpireOn = rv.CreatedOn.Add(time.Duration(ctx.Config.ReviewExp))
			res.Token = &rtr
		} else {
			res.Type = ContentReviewType
			res.Content = new(ReviewContentResponse)
			res.Content.ID = rv.ID
			res.Content.UserInfo = rv.UserInfo
			res.Content.ReviewContent = rv.ReviewContent
		}

		g.AbortWithStatusJSON(200, res)
	}
}

func (ctx *Ctx) HandlerGetPubReviews() gin.HandlerFunc {
	return func(g *gin.Context) {
		trainingIDstr := g.Query("training_id")
		if trainingIDstr == "" {
			g.AbortWithError(400, fmt.Errorf("invalid training_id"))
			return
		}
		trainingID, err := uuid.Parse(trainingIDstr)
		if err != nil {
			g.AbortWithError(400, fmt.Errorf("parse training_id: %v", err))
			return
		}
		if trainingID == uuid.Nil {
			g.AbortWithError(400, fmt.Errorf("empty training_id"))
			return
		}
		rv, err := ctx.DalReadPubReviews(ReadReviewsOpts{TrainingID: &trainingID}, nil)
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		g.AbortWithStatusJSON(200, rv)
	}
}
