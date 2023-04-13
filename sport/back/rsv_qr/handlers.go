package rsv_qr

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"sport/helpers"
	"sport/rsv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

/*
	create QR code
*/
func (ctx *Ctx) HandlersPostQr() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var req CreateQrCodeRequest

		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body, &req, req.Validate); err != nil {

			g.AbortWithError(400, err)
			return
		}

		if g.GetHeader("Authorization") != "" {
			var userID uuid.UUID
			userID, err = ctx.Api.AuthorizeUserFromCtx(g)
			if err != nil {
				g.AbortWithError(401, err)
				return
			}
			_, err = ctx.Rsv.ReadSingleRsv(rsv.ReadRsvsArgs{
				ID:     &req.RsvID,
				UserID: &userID,
			})
		} else {
			if req.AccessToken == nil {
				g.AbortWithError(401, fmt.Errorf("no access_token provided"))
				return
			}
			_, err = ctx.Rsv.ReadSingleRsv(rsv.ReadRsvsArgs{
				AccessToken: req.AccessToken,
			})
		}

		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		qr := req.NewQrCode()

		if c, err := ctx.DalCanCreateQr(qr); err != nil {
			g.AbortWithError(500, err)
			return
		} else {
			if !c {
				g.AbortWithError(425, fmt.Errorf("user exceeded allowed number of qr codes"))
				return
			}
		}

		if err := ctx.DalCreateQr(qr); err != nil {
			g.AbortWithError(500, err)
			return
		}

		u := qr.ToInterfaceUrl(ctx)

		png, err := qrcode.Encode(u, qrcode.High, req.Size)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}
		if req.DataUrl {
			g.AbortWithStatus(200)
			b64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
			g.Writer.WriteString(b64)
		} else {
			g.Header("Content-Type", "image/png")
			g.AbortWithStatus(200)
			g.Writer.Write(png)
		}
	}
}

/*
	evaluate QR code
*/
func (ctx *Ctx) HandlersEvalQr() gin.HandlerFunc {
	return func(g *gin.Context) {

		ids := g.Query("id")
		if ids == "" {
			g.AbortWithError(400, fmt.Errorf("empty id"))
			return
		}
		id, err := uuid.Parse(ids)
		if err != nil {
			g.AbortWithError(400, err)
			return
		}

		qr, err := ctx.DalReadSingleQr(id)
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		userID := g.MustGet("UserID").(uuid.UUID)

		r, err := ctx.Rsv.ReadSingleRsv(rsv.ReadRsvsArgs{
			ID:             &qr.RsvID,
			InstrUserID:    &userID,
			WithInstructor: true,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		res := EvalQrResponse{
			RsvID:       qr.RsvID,
			ConfirmCode: 0,
		}

		if r.InstructorID == nil {
			g.AbortWithError(500, fmt.Errorf("instructor doesnt exist"))
			return
		}

		if !r.IsConfirmed {
			res.ConfirmCode |= NotCaptured
		}

		if r.QrConfirmed {
			res.ConfirmCode |= AlreadyConfirmed
		}

		if res.ConfirmCode == 0 {
			if err := ctx.Rsv.DalConfirmQrCode(qr.RsvID, *r.InstructorID); err != nil {
				g.AbortWithError(500, err)
				return
			}
		}

		g.AbortWithStatusJSON(200, res)
	}
}
