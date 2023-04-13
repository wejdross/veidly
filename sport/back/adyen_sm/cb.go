package adyen_sm

import (
	"database/sql"
	"sport/adyen"

	"github.com/gin-gonic/gin"
)

type SmGenericAfterCallback func(d *Pass)
type SmGenericBeforeCallback func(
	d *Pass,
	args map[string]interface{},
	tx *sql.Tx,
	sd StateChangeEventSource) (canContinue bool, err error)

type GetPassInWebhookCallback func(ni *adyen.NotificationRequestItem) (*Pass, error)

type HandlerCallbackSource int

const (
	SrcDispute    HandlerCallbackSource = 1
	SrcExpireLink HandlerCallbackSource = 2
	SrcCancel     HandlerCallbackSource = 3
	SrcRefund     HandlerCallbackSource = 4
)

type DisputeInfo map[string]string

type SmCallbacks struct {

	// those callbacks will be called directly before or after state change
	// all are optional
	BeforeLinkExpire     SmGenericBeforeCallback
	BeforeCancelOrRefund SmGenericBeforeCallback
	AfterPayout          SmGenericAfterCallback
	BeforeWaitRefund     SmGenericBeforeCallback
	BeforeRefund         SmGenericBeforeCallback
	BeforeHold           SmGenericBeforeCallback
	BeforeCapture        SmGenericBeforeCallback
	BeforeWaitPayout     SmGenericBeforeCallback

	// mailing, all functions are optional
	EmailUserAboutHold          SmGenericAfterCallback
	EmailInstrAboutCancel       SmGenericAfterCallback
	EmailUserAboutCancel        SmGenericAfterCallback
	EmailUserAboutCapture       SmGenericAfterCallback
	EmailInstrAboutDispute      SmGenericAfterCallback
	EmailUserAboutDispute       SmGenericAfterCallback
	EmailUserAboutFailedCapture SmGenericAfterCallback
	EmailInstrAboutFailedPayout SmGenericAfterCallback

	// functions used to obtain objects in specified places
	// all are required
	GetPassInWebhook GetPassInWebhookCallback
	GetPassInAgent   func() ([]Pass, error)
	GetPassInHandler func(g *gin.Context, c SmHandlerCaller, rawBody []byte) (*Pass, error)
}
