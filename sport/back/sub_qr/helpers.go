package sub_qr

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"net/url"

	"github.com/google/uuid"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func (ctx *Ctx) TestVerifyQr(png []byte, subID uuid.UUID) (uuid.UUID, error) {
	img, _, err := image.Decode(bytes.NewReader(png))
	if err != nil {
		return uuid.Nil, err
	}

	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return uuid.Nil, err
	}

	goturl := result.GetText()
	if _, err := url.ParseRequestURI(goturl); err != nil {
		return uuid.Nil, err
	}

	ucodes, err := ctx.DalReadQrForRsv(subID)
	if err != nil {
		return uuid.Nil, err
	}

	if len(ucodes) != 1 {
		return uuid.Nil, fmt.Errorf("invalid number of user codes, expected 1 got %d", len(ucodes))
	}

	expectedUrl := ucodes[0].ToInterfaceUrl(ctx)

	if goturl != expectedUrl {
		id := uuid.New()
		err := ioutil.WriteFile(id.String()+".png", png, 0600)
		if err != nil {
			return uuid.Nil, err
		}
		return uuid.Nil, fmt.Errorf("%s: invalid qr code \n\t%v \ngot \n\t%v", id, expectedUrl, goturl)
	}

	return ucodes[0].ID, nil
}
