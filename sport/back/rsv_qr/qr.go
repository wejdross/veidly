package rsv_qr

import (
	"fmt"

	"github.com/google/uuid"
)

type QrCodeRequest struct {
	RsvID uuid.UUID
}

type QrCode struct {
	ID uuid.UUID
	QrCodeRequest
}

type CreateQrCodeRequest struct {
	QrCodeRequest
	AccessToken *uuid.UUID
	DataUrl     bool
	Size        int
}

func (q *CreateQrCodeRequest) Validate() error {
	if q.Size < 0 || q.Size > 512 {
		return fmt.Errorf("validate CreateQrCodeRequest: invalid size")
	}
	if q.Size == 0 {
		q.Size = -1
	}
	return q.QrCodeRequest.Validate()
}

func (q *QrCodeRequest) Validate() error {
	if q.RsvID == uuid.Nil {
		return fmt.Errorf("validate QrCodeRequest: invalid RsvID")
	}
	return nil
}

func (q *QrCodeRequest) NewQrCode() *QrCode {
	return &QrCode{
		QrCodeRequest: *q,
		ID:            uuid.New(),
	}
}

func (q *QrCode) ToInterfaceUrl(ctx *Ctx) string {
	return fmt.Sprintf(ctx.Config.QrEvalUrlFmt, q.ID)
}

type QrEvalConfirmCode int

const (
	AlreadyConfirmed QrEvalConfirmCode = 1
	NotCaptured      QrEvalConfirmCode = 2
)

type EvalQrResponse struct {
	RsvID       uuid.UUID
	ConfirmCode QrEvalConfirmCode
}
