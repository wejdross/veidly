package sub

import (
	"database/sql"
	"fmt"
	"sport/adyen"
	"sport/helpers"
	"sport/user"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApiSubRequest struct {
	SubModelID uuid.UUID
	//DurationIx int
}

func (r *ApiSubRequest) Validate() error {
	if r.SubModelID == uuid.Nil {
		return fmt.Errorf("Validate ApiSubRequest: invalid SubModelID")
	}
	// if r.DurationIx < 0 {
	// 	return fmt.Errorf("Validate ApiSubRequest: invalid DurationIx")
	// }
	return nil
}

type ApiSubResponse struct {
	ID  uuid.UUID
	Url string
}

func (ctx *Ctx) HandlerPostSub() gin.HandlerFunc {
	return func(g *gin.Context) {

		var apiReq ApiSubRequest
		userID := g.MustGet("UserID").(uuid.UUID)

		if err := helpers.ReadJsonBodyFromReader(
			g.Request.Body,
			&apiReq,
			apiReq.Validate,
		); err != nil {
			g.AbortWithError(400, err)
			return
		}

		model, err := ctx.DalReadSingleSubModel(ReadSubModelRequest{
			ID: &apiReq.SubModelID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		u, err := ctx.User.DalReadUser(userID, user.KeyTypeID, true)
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(401, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		// if len(model.Durations) <= apiReq.DurationIx {
		// 	g.AbortWithError(400, fmt.Errorf("DurationIx doesnt exist"))
		// 	return
		// }

		dur := time.Duration(model.Duration) * time.Hour * 24

		now := time.Now().In(time.UTC)

		sr := SubRequest{
			ID:        uuid.New(),
			SubModel:  model,
			RefID:     uuid.New(),
			AdyenRes:  nil,
			UserID:    userID,
			DateStart: now,
			DateEnd:   now.Add(dur),
		}

		retUrl := fmt.Sprintf(ctx.Config.SubDetailsUrlFmt, sr.ID)

		ar := adyen.RegisterTransactionRequest{
			SessionID:   AdyenSmPrefix + sr.RefID.String(),
			Amount:      int(model.Total()),
			Currency:    model.Currency,
			Description: model.Name,
			Email:       u.Email,
			Country:     u.Country,
			Language:    u.Language,
			UrlReturn:   retUrl,
		}

		if ctx.Adyen.Mockup {
			sr.AdyenRes = &adyen.CreatePaymentLinkResponse{
				Url: retUrl,
				ID:  uuid.New().String(),
			}
		} else {
			sr.AdyenRes, err = ctx.Adyen.RegisterTransaction(&ar)
			if err != nil {
				g.AbortWithError(500, err)
				return
			}
		}

		s, err := sr.NewSub(ctx)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		if err := ctx.DalCreateSub(s); err != nil {
			g.AbortWithError(500, err)
			return
		}

		res := ApiSubResponse{
			ID:  s.ID,
			Url: sr.AdyenRes.Url,
		}

		g.AbortWithStatusJSON(200, res)
	}
}

func (ctx *Ctx) HandlerGetSubs() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)

		t := g.Param("type")

		req := ReadSubRequest{}

		switch t {
		case "instructor":
			req.InstrUserID = &userID
			break
		case "user":
			fallthrough
		default: // user
			req.UserID = &userID
		}

		if tmp := g.Query("id"); tmp != "" {
			id, err := uuid.Parse(tmp)
			if err != nil {
				g.AbortWithError(400, err)
				return
			}
			req.ID = &id
		}

		s, err := ctx.DalReadSubWithJoins(req)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		g.AbortWithStatusJSON(200, s)
	}
}
