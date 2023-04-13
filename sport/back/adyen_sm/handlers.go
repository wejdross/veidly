package adyen_sm

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sport/helpers"

	"github.com/gin-gonic/gin"
)

type DisputeRequest struct {
	Email string
	Msg   string
}

func (req *DisputeRequest) Validate() error {

	if req.Email == "" || req.Msg == "" {
		return fmt.Errorf("invalid message")
	}

	if len(req.Email) > 64 || len(req.Msg) > 256 {
		return fmt.Errorf("invalid message")
	}

	return nil
}

type SmHandlerCaller string

const (
	User  SmHandlerCaller = "user"
	Instr SmHandlerCaller = "instructor"
)

func GetSmHandlerCaller(g *gin.Context) (SmHandlerCaller, error) {
	t := g.Param("type")
	switch t {
	case "instructor":
		return Instr, nil
	default:
		return User, nil
	}
}

/*
	request help from support
*/
func (ctx *Ctx) HandlerDispute() gin.HandlerFunc {
	return func(g *gin.Context) {

		var req DisputeRequest

		var b []byte
		var err error
		if b, err = ioutil.ReadAll(g.Request.Body); err != nil {
			g.AbortWithError(400, err)
			return
		}

		if err = helpers.ParseAndValidateJson(b, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}

		c, err := GetSmHandlerCaller(g)
		if err != nil {
			g.AbortWithError(400, err)
			return
		}

		pass, err := ctx.GetPassInHandler(g, c, b)
		if err != nil {
			if !g.IsAborted() {
				g.AbortWithError(400, err)
			}
			return
		}

		switch pass.Data.State {
		case Hold:
			fallthrough
		case Link:
			fallthrough
		case LinkExpress:
			fallthrough
		case CancelOrRefund:
			fallthrough
		case Capture:
			fallthrough
		case Payout:
			fallthrough
		case Dispute:
			fallthrough
		case Refund:

			ed := map[string]string{}

			ed["ReturnEmail"] = req.Email
			ed["Msg"] = req.Msg
			ed["Type"] = pass.Data.SmType
			ed["ID"] = pass.Data.ID.String()
			ed["Who"] = string(c)

			j, err := json.Marshal(ed)
			if err != nil {
				g.AbortWithError(500, err)
				return
			}
			p := fmt.Sprintf("application/json\r\n%s", j)
			err = ctx.MoveStateToDispute(pass, StateChangeEventSource{Handler, p})
			if err != nil {
				g.AbortWithError(500, err)
				return
			}
		default:
			g.AbortWithError(
				412,
				fmt.Errorf("cant make this request for curent rsv state: %v", pass.Data.State))
			return
		}

		g.AbortWithStatus(204)
	}
}

// func (ctx *Ctx) HandlerExpireLink() gin.HandlerFunc {
// 	return func(g *gin.Context) {

// 		var b []byte
// 		var err error
// 		if b, err = ioutil.ReadAll(g.Request.Body); err != nil {
// 			g.AbortWithError(400, err)
// 			return
// 		}

// 		c, err := GetSmHandlerCaller(g)
// 		if err != nil {
// 			g.AbortWithError(400, err)
// 			return
// 		}

// 		pass, err := ctx.GetPassInHandler(g, c, b)
// 		if err != nil {
// 			if !g.IsAborted() {
// 				g.AbortWithError(400, err)
// 			}
// 			return
// 		}

// 		if pass.Data.State != Link && pass.Data.State != LinkExpress {
// 			g.AbortWithError(
// 				412,
// 				fmt.Errorf("cant expire link for curent obj state: %v", pass.Data.State))
// 			return
// 		}

// 		err = ctx.MoveStateToLinkExpire(pass, StateChangeEventSource{Handler, "expired by " + string(c)})
// 		if err != nil {
// 			g.AbortWithError(500, err)
// 			return
// 		}

// 		g.AbortWithStatus(204)
// 	}
// }

func (ctx *Ctx) HandlerCancel() gin.HandlerFunc {
	return func(g *gin.Context) {

		var b []byte
		var err error
		if b, err = ioutil.ReadAll(g.Request.Body); err != nil {
			g.AbortWithError(400, err)
			return
		}

		c, err := GetSmHandlerCaller(g)
		if err != nil {
			g.AbortWithError(400, err)
			return
		}

		pass, err := ctx.GetPassInHandler(g, c, b)
		if err != nil {
			if !g.IsAborted() {
				g.AbortWithError(400, err)
			}
			return
		}

		if pass.Data.State != Hold {
			g.AbortWithError(
				412,
				fmt.Errorf("cant make this request for curent obj state: %v", pass.Data.State))
			return
		}

		sd := StateChangeEventSource{Handler, "cancelled by " + string(c)}

		if ctx.adyen.Mockup {
			err = ctx.MoveStateToCancelOrRefund(pass, sd, nil)
		} else {
			err = ctx.MoveStateToWaitCancelOrRefund(pass, sd, nil)
		}

		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		g.AbortWithStatus(204)
	}
}

func (ctx *Ctx) HandlerRefund() gin.HandlerFunc {
	return func(g *gin.Context) {

		var b []byte
		var err error
		if b, err = ioutil.ReadAll(g.Request.Body); err != nil {
			g.AbortWithError(400, err)
			return
		}

		c, err := GetSmHandlerCaller(g)
		if err != nil {
			g.AbortWithError(400, err)
			return
		}

		pass, err := ctx.GetPassInHandler(g, c, b)
		if err != nil {
			if !g.IsAborted() {
				g.AbortWithError(400, err)
			}
			return
		}

		switch pass.Data.State {
		case WaitCapture:
			// waiting for confirmation about capture
			// TODO: maybe we could queue refund request to process it as soon as confirmation comes
			// or maybe adyen has endpoint to do it (cancel / technicalCancel ...?)
			g.AbortWithError(
				412,
				fmt.Errorf("cant make this request for curent rsv state: %v", pass.Data.State))
			return
		case Capture:

			if ctx.adyen.Mockup {
				sd := StateChangeEventSource{Handler, "refund requested by " + string(c)}
				err = ctx.MoveStateToRefund(pass, sd)
				if err != nil {
					g.AbortWithError(500, err)
				}
				break
			}

			// increase amount of user fuckups
			var tx *sql.Tx
			if c == Instr {
				if tx, err = ctx.dal.Db.Begin(); err != nil {
					g.AbortWithError(500, err)
					return
				}
				penalty := (pass.Data.InstrPrice * ctx.Config.InstrShotPenaltyPercent) / 100
				if err := ctx.instr.DalConditionalAddPenalty(pass.Data.InstructorID,
					penalty,
					ctx.Config.FreeInstrShots,
					tx,
				); err != nil {
					g.AbortWithError(500, err)
					return
				}
			}
			/*
				if this refund is made by instructor
				then we refund full amount to the user
				otherwise if client makes this cancellation
				then refund only part of the amount
			*/
			var sd StateChangeEventSource

			sd = StateChangeEventSource{Handler, "refund requested by " + string(c)}
			err = ctx.MoveStateToWaitRefund(pass, sd, true, c == Instr)
			if err != nil {
				if tx != nil {
					tx.Rollback()
				}
				if he := helpers.HttpErr(err); he != nil {
					he.WriteAndAbort(g)
				} else {
					g.AbortWithError(500, err)
				}
				return
			}
			if tx != nil {
				// if this fails right now, well our loss
				// better to loose one refund from instructor than interrupting whole process
				// ... i think
				err := tx.Commit()
				if err != nil {
					fmt.Println(err)
				}
			}
		case Payout:
			fallthrough
		case WaitPayout:
			/*
				i dont know how to process refund after we payed instructor his share
				gotta resolve it manually
			*/
			// err = ctx.MoveStateToDispute(
			// 	&rsv.Reservation, Handler, t+" refund")
			// if err != nil {
			// 	g.AbortWithError(500, err)
			// 	return
			// }
			fallthrough
		default:
			g.AbortWithError(
				412,
				fmt.Errorf("cant make this request for curent rsv state: %v", pass.Data.State))
			return
		}

		g.AbortWithStatus(204)
	}
}

type DecisionRequest struct {
	Decision InstructorDecision
}

func (req *DecisionRequest) Validate() error {
	if req.Decision != Approve && req.Decision != Reject {
		return fmt.Errorf("validate DecisionRequest: invalid decision")
	}
	return nil
}

func (ctx *Ctx) HandlerDecision() gin.HandlerFunc {
	return func(g *gin.Context) {

		var b []byte
		var err error
		if b, err = ioutil.ReadAll(g.Request.Body); err != nil {
			g.AbortWithError(400, err)
			return
		}

		var req DecisionRequest
		if err = helpers.ParseAndValidateJson(b, &req, req.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}

		pass, err := ctx.GetPassInHandler(g, Instr, b)
		if err != nil {
			if !g.IsAborted() {
				g.AbortWithError(400, err)
			}
			return
		}

		if req.Decision == Approve {

			if ctx.adyen.Mockup {
				if err := ctx.MoveStateToCapture(
					pass,
					StateChangeEventSource{Handler, "instructor confirmed"},
					0,
					map[string]interface{}{
						ctx.Model.Columns.InstructorDecision: Approve,
					},
				); err != nil {
					if err == sql.ErrNoRows {
						g.AbortWithError(404, err)
					} else {
						g.AbortWithError(500, err)
					}
					return
				}
			} else {

				// we could move it right now to wait_capture,
				// but im willing to give user some more time to decide
				if err := ctx.UpdateInstrDecision(pass, req.Decision); err != nil {

					if err == sql.ErrNoRows {
						g.AbortWithError(404, err)
					} else {
						g.AbortWithError(500, err)
					}
					return
				}
			}

		} else if req.Decision == Reject {

			switch pass.Data.State {
			case Hold:
				// no waiting here, move to cancel
				if ctx.adyen.Mockup {
					err = ctx.MoveStateToCancelOrRefund(
						pass,
						StateChangeEventSource{Handler, "instructor rejected"},
						map[string]interface{}{
							ctx.Model.Columns.InstructorDecision: Reject,
						},
					)
				} else {
					err = ctx.MoveStateToWaitCancelOrRefund(
						pass,
						StateChangeEventSource{Handler, "instructor rejected"},
						map[string]interface{}{
							ctx.Model.Columns.InstructorDecision: Reject,
						})
				}
			default:
				g.AbortWithError(
					412,
					fmt.Errorf("cant make this request for curent state: %v", pass.Data.State))
				return
			}

			if err != nil {
				g.AbortWithError(500, err)
				return
			}
		}

		g.AbortWithStatus(204)
	}
}
