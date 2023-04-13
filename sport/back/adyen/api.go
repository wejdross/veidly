package adyen

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type NotificationRequestItem struct {
	MerchantID   int    `json:"merchantId"`
	PosID        int    `json:"posId"`
	SessionID    string `json:"sessionId"`
	Amount       int    `json:"amount"`
	OriginAmount int    `json:"originAmount"`
	Currency     string `json:"currency"`
	OrderID      int    `json:"orderId"`
	MethodID     int    `json:"methodId"`
	Statement    string `json:"statement"`
	Sign         string `json:"sign"`
}

type RegisterTransactionRequest struct {
	// unique identifier from our system
	SessionID string `json:"sessionId"`
	// ex. 23zł 1gr is 2301
	Amount int `json:"amount"`
	// ISO currency code, ex: PLN, ...
	Currency string `json:"currency"`
	// transaction description
	Description string `json:"description"`
	// client's email
	Email string `json:"email"`
	// ISO country code, PL, DE, ...
	Country string `json:"country"`
	// ISO 639-1 language: bg, cs, de, en, es, fr, hr, hu, it, nl, pl, pt, se, sk
	Language string `json:"language"`
	// Return address
	UrlReturn string `json:"urlReturn"`
	// Notification address
	CTX_UrlStatus string `json:"urlStatus"`

	// CTX_* fields are set during setCtx() invokation

	// ffs just use AEAD cipher, and don't bother with this pitiful attempt at sEcUriTy
	CTX_Sign string `json:"sign"`
	// Shop ID
	CTX_MerchantID int `json:"merchantId"`
	// Shop id
	CTX_PosID int `json:"posId"`
}

type RegisterTransactionResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

type CreatePaymentLinkResponse struct {
	Url string `json:"url,omitempty"`
	ID  string `json:"id,omitempty"`
}

func (ctx *Ctx) AuthHeader() (key, val string) {
	key = "Authorization"
	kv := fmt.Sprintf("%d:%s", ctx.Config.Auth.Username, ctx.Config.Auth.Password)
	val = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(kv)))
	return
}

func p24hash(o interface{}) string {
	b, err := json.Marshal(o)
	fmt.Println(string(b))
	if err != nil {
		panic("failed to marshal p24hash obj, err was " + err.Error())
	}
	sum := sha512.Sum384(b)
	return hex.EncodeToString(sum[:])
}

func (r *RegisterTransactionRequest) SetCtx(ctx *Ctx) {
	var h struct {
		SessionID  string `json:"sessionId"`
		MerchantID int    `json:"merchantId"`
		Amount     int    `json:"amount"`
		Currency   string `json:"currency"`
		Crc        string `json:"crc"`
	}

	r.CTX_MerchantID = ctx.Config.Auth.Username
	r.CTX_PosID = ctx.Config.Auth.Username
	r.CTX_UrlStatus = ctx.whUrl()

	h.SessionID = r.SessionID
	h.MerchantID = r.CTX_MerchantID
	h.Amount = r.Amount
	h.Currency = r.Currency
	h.Crc = ctx.Config.Crc
	r.CTX_Sign = p24hash(h)
}

func (ctx *Ctx) p24request(
	method, path string, req, res interface{},
) error {
	url := fmt.Sprintf("%s%s", ctx.Config.BaseUrl, path)
	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}
	fmt.Println(method + " " + url + " " + string(reqBody))
	httpReq, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	httpReq.Header.Add(ctx.AuthHeader())
	httpReq.Header.Add("Content-Type", "application/json")
	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()
	resJson, resJsonErr := io.ReadAll(httpRes.Body)
	if httpRes.StatusCode < 200 || httpRes.StatusCode > 299 {
		fmt.Println(string(resJson))
		return fmt.Errorf("p24 request: '%s %s' failed, p24 returned: %s %s",
			method, path, httpRes.Status, string(resJson))
	}
	if resJsonErr != nil {
		return fmt.Errorf("p24 request: '%s %s' failed to unmarshal body, err was: %v",
			method, path, resJsonErr)
	}
	return json.Unmarshal(resJson, res)
}

func (ctx *Ctx) RegisterTransaction(
	req *RegisterTransactionRequest,
) (*CreatePaymentLinkResponse, error) {
	req.SetCtx(ctx)
	p24res := new(RegisterTransactionResponse)
	var err error
	if err != nil {
		return nil, err
	}
	if err = ctx.p24request(
		"POST", "/api/v1/transaction/register", req, p24res); err != nil {
		return nil, err
	}
	var res CreatePaymentLinkResponse
	res.ID = p24res.Data.Token
	res.Url = ctx.Config.BaseUrl + "/trnRequest/" + res.ID
	return &res, nil
}

type VerifyTransactionRequest struct {
	// unique identifier from our system
	SessionID string `json:"sessionId"`
	// ex. 23zł 1gr is 2301
	Amount int `json:"amount"`
	// ISO currency code, ex: PLN, ...
	Currency string `json:"currency"`
	// orderId assigned by p24.
	// this is taken from notification. Not response from the registerTransaction
	// -- that would be silly!
	OrderID int `json:"orderId"`
	//  {"sessionId":"str","orderId":int,"amount":int,"currency":"str","crc":"str"}
	CTX_Sign string `json:"sign"`
	// Shop ID
	CTX_MerchantID int `json:"merchantId"`
	// Shop id
	CTX_PosID int `json:"posId"`
}

func (r *VerifyTransactionRequest) SetCtx(ctx *Ctx) {
	r.CTX_MerchantID = ctx.Config.Auth.Username
	r.CTX_PosID = ctx.Config.Auth.Username
	h := struct {
		SessionID string `json:"sessionId"`
		OrderID   int    `json:"orderId"`
		Amount    int    `json:"amount"`
		Currency  string `json:"currency"`
		Crc       string `json:"crc"`
	}{
		SessionID: r.SessionID,
		OrderID:   r.OrderID,
		Amount:    r.Amount,
		Currency:  r.Currency,
		Crc:       ctx.Config.Crc,
	}
	r.CTX_Sign = p24hash(h)
}

type VerifyTransactionResponse struct {
	ResponseCode int `json:"responseCode"`
}

func (ctx *Ctx) VerifyTransaction(
	req *VerifyTransactionRequest,
) (*VerifyTransactionResponse, error) {
	res := new(VerifyTransactionResponse)
	req.SetCtx(ctx)
	return res, ctx.p24request(
		"PUT", "/api/v1/transaction/verify",
		req, res)
}

type RefundTransactionItem struct {
	OrderID     string `json:"orderId"`
	SessionID   string `json:"sessionId"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type RefundTransactionRequest struct {
	RequestID     string                  `json:"requestId"`
	Refunds       []RefundTransactionItem `json:"refunds"`
	RefundsUUid   uuid.UUID               `json:"refundsUuid"`
	CTX_UrlStatus string                  `json:"urlStatus"`
}

type RefundTransactionResponse struct {
	ResponseCode int `json:"responseCode"`
}

func (r *RefundTransactionRequest) SetCtx(ctx *Ctx) {
	r.CTX_UrlStatus = ctx.whUrl()
}

func (ctx *Ctx) RefundTransaction(
	req *RefundTransactionRequest,
) (*RefundTransactionResponse, error) {
	res := new(RefundTransactionResponse)
	req.SetCtx(ctx)
	return res, ctx.p24request(
		"POST", "/api/v1/transaction/refund",
		req, res)
}
