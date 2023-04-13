package api

import (
	"encoding/base64"
	"fmt"
	"sport/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func JwtHeader(token string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

// JwtRequest is used to create proper Jwt structure
type JwtRequest struct {
	Hmac64      string        `yaml:"hmac64"`
	DurationHrs time.Duration `yaml:"duration_hrs"`
}

// NewJwtRequestFromConfig creates new jwt request using configuration
func NewJwtRequestFromConfig(config *config.Ctx) (*JwtRequest, error) {
	var wrapper struct {
		Auth JwtRequest `yaml:"auth"`
	}
	if err := config.Unmarshal(&wrapper); err != nil {
		return nil, err
	}
	if err := wrapper.Auth.Validate(); err != nil {
		return nil, err
	}
	return &wrapper.Auth, nil
}

// Validate validates that jwt request holds valid data
func (req *JwtRequest) Validate() error {
	if len(req.Hmac64) == 0 {
		return fmt.Errorf("Expected hmac to be > 0 length")
	}
	if req.DurationHrs <= 0 {
		return fmt.Errorf("Expected DurationHrs to be > 0")
	}
	return nil
}

// Jwt is Used to sign and verify JWT tokens
type Jwt struct {
	SigningKey []byte
	request    *JwtRequest
}

func (req *JwtRequest) NewJwt() (*Jwt, error) {
	ret := new(Jwt)
	ret.request = req
	var err error
	if ret.SigningKey, err = base64.StdEncoding.DecodeString(
		req.Hmac64,
	); err != nil {
		return nil, err
	}
	return ret, err
}

// JwtPayload is the jwt token payload
type JwtPayload struct {
	UserID uuid.UUID
	Exp    time.Time
}

func NewJwtFromConfig(config *config.Ctx) *Jwt {
	req, err := NewJwtRequestFromConfig(config)
	if err != nil {
		panic(err)
	}
	auth, err := req.NewJwt()
	if err != nil {
		panic(err)
	}
	return auth
}

/*
	Generate jwt token, encoding {userid} in payload,
	which will be valid for {hours} amount of time.

	Signing key from global Auth context will be used as HMAC-SHA-256 key.
*/
func (auth *Jwt) GenToken(payload *JwtPayload) (string, error) {

	var jwtToken *jwt.Token
	var ret string
	var err error

	if payload == nil {
		return "", fmt.Errorf("AuthGenToken: invalid argument: payload")
	}

	if !payload.Exp.After(time.Now()) {
		payload.Exp = time.Now().Add(time.Hour * auth.request.DurationHrs)
	}

	if len(payload.UserID) == 0 {
		return "", fmt.Errorf("AuthGenToken: userid was empty")
	}

	claims := jwt.MapClaims{}

	claims["uid"] = payload.UserID
	claims["exp"] = payload.Exp.Unix()

	/* can this return nil token... ? */
	if jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, claims); jwtToken == nil {
		return "", fmt.Errorf("token was nil")
	}

	if ret, err = jwtToken.SignedString(auth.SigningKey); err != nil {
		return "", err
	}

	return ret, nil
}

/*
	Validation token callback
*/
func (auth *Jwt) TokenKeyFunc(jwtToken *jwt.Token) (interface{}, error) {
	var ok bool
	if _, ok = jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Authorize: Invalid signing method")
	}
	return auth.SigningKey, nil
}

func (auth *Jwt) ValidateToken(token string) error {
	var jwtToken *jwt.Token
	var err error

	if jwtToken, err = jwt.Parse(token, auth.TokenKeyFunc); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	err = jwtToken.Claims.Valid()

	return err
}

const jwtAuthType = "Bearer"

func (ctx *Ctx) AuthorizeUserFromToken(tb string) (uuid.UUID, error) {
	return ctx.Jwt.AuthorizeUserFromToken(tb)
}

func (ctx *Jwt) AuthorizeUserFromToken(tb string) (uuid.UUID, error) {
	var claims jwt.MapClaims
	var ok bool
	if tb == "" {
		return uuid.Nil, fmt.Errorf("JwtAuthorize: No token in header")
	}
	token, err := jwt.Parse(tb, ctx.TokenKeyFunc)
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("JwtAuthorize: Token invalid")
	}
	/* is exp validated here? */
	if claims, ok = token.Claims.(jwt.MapClaims); ok {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			return uuid.Nil, fmt.Errorf("JwtAuthorize: Token expired")
		}
	} else {
		return uuid.Nil, fmt.Errorf("JwtAuthorize: Token claims invalid")
	}

	uidi, ok := claims["uid"]
	if !ok {
		return uuid.Nil, fmt.Errorf("JwtAuthorize: uid was empty")
	}

	uids := ""
	var uid uuid.UUID

	if uids, ok = uidi.(string); !ok || len(uids) == 0 {
		return uuid.Nil, fmt.Errorf("JwtAuthorize: invalid uid")
	}

	if uid, err = uuid.Parse(uids); err != nil {
		return uuid.Nil, fmt.Errorf("JwtAuthorize: invalid uid")
	}

	return uid, nil
}

func GetTokenFromHdr(c *gin.Context) (string, error) {
	tb := c.GetHeader("Authorization")
	if tb == "" {
		return "", fmt.Errorf("ParseTokenFromHdr: No token in header")
	}

	if !strings.HasPrefix(tb, jwtAuthType) || len(tb) < len(jwtAuthType)+1 {
		return "",
			fmt.Errorf("ParseTokenFromHdr: Didnt find required Authorization: " + jwtAuthType)
	}

	tb = tb[len(jwtAuthType)+1:]

	return tb, nil
}

func (api *Ctx) AuthorizeUserFromCtx(c *gin.Context) (uuid.UUID, error) {
	tb, err := GetTokenFromHdr(c)
	if err != nil {
		return uuid.Nil, err
	}
	return api.AuthorizeUserFromToken(tb)
}

func (api *Ctx) jwtMiddleware(c *gin.Context) {

	uid, err := api.AuthorizeUserFromCtx(c)
	if err != nil {
		c.AbortWithError(401, err)
		return
	}

	c.Set("UserID", uid)

	c.Next()
}
