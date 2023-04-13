package sub_qr

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"sport/helpers"
	"sport/sub"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

/*
	create QR code
*/
func (ctx *Ctx) CreateQrHandler() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var req CreateQrCodeRequest

		userID := g.MustGet("UserID").(uuid.UUID)

		if err = helpers.ReadJsonBodyFromReader(g.Request.Body, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}

		_, err = ctx.Sub.DalReadSingleSub(sub.ReadSubRequest{
			ID:     &req.SubID,
			UserID: &userID,
		})

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
func (ctx *Ctx) EvalQrHandler() gin.HandlerFunc {
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

		s, err := ctx.Sub.DalReadSingleSubWithJoins(sub.ReadSubRequest{
			ID:          &qr.SubID,
			InstrUserID: &userID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		g.AbortWithStatusJSON(200, s)
	}
}

/*
	confirm QR code
*/
func (ctx *Ctx) ConfirmQrHandler() gin.HandlerFunc {
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

		if err = ctx.Sub.DalConfirmSub(qr.SubID, userID); err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(409, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		g.AbortWithStatus(204)
	}
}
