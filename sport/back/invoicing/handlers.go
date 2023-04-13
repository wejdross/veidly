package invoicing

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctx *Ctx) GetInvoiceHandler() gin.HandlerFunc {
	return func(g *gin.Context) {

		idStr := g.Query("id")
		var id uuid.UUID
		var err error

		if idStr != "" {
			id, err = uuid.Parse(idStr)
			if err != nil {
				g.AbortWithError(400, err)
				return
			}
		}

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

		invoices, err := ctx.DalGetInvoices(id, iid)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		g.AbortWithStatusJSON(200, invoices)
	}
}

func (ctx *Ctx) PrintInvoiceHandler() gin.HandlerFunc {
	return func(g *gin.Context) {

		idStr := g.Query("id")
		var id uuid.UUID
		var err error

		if idStr != "" {
			id, err = uuid.Parse(idStr)
			if err != nil {
				g.AbortWithError(400, err)
				return
			}
		}

		// IMPORTANT:
		// invoice ID must be kept secret from other users
		invoices, err := ctx.DalGetInvoices(id, uuid.Nil)
		if err != nil {
			g.AbortWithError(500, err)
			return
		}
		if len(invoices) != 1 {
			g.AbortWithError(404, fmt.Errorf("target invoice not found"))
			return
		}

		invoice := &invoices[0]

		bytes, err := PrintInvoice(invoice)

		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		g.Header("Content-Disposition",
			fmt.Sprintf("attachment; filename=invoice-%s.pdf", invoice.Number))
		g.Data(200, "application/pdf", bytes)
	}
}
