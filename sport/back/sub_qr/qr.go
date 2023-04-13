package sub_qr

import (
	"fmt"

	"github.com/google/uuid"
)

type QrCodeRequest struct {
	SubID uuid.UUID
}

type QrCode struct {
	ID uuid.UUID
	QrCodeRequest
}

type CreateQrCodeRequest struct {
	QrCodeRequest
	DataUrl bool
	Size    int
}

func (q *QrCodeRequest) Validate() error {
	if q.SubID == uuid.Nil {
		return fmt.Errorf("validate QrCodeRequest: invalid SubID")
	}
	return nil
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

func (q *QrCodeRequest) NewQrCode() *QrCode {
	return &QrCode{
		QrCodeRequest: *q,
		ID:            uuid.New(),
	}
}

func (q *QrCode) ToInterfaceUrl(ctx *Ctx) string {
	return fmt.Sprintf(ctx.Config.QrEvalUrlFmt, q.ID)
}
