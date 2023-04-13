package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OauthCtx struct {
	GoogleConfig oauth2.Config
}

func OauthGetProviderAuthUrl(config *oauth2.Config, returnUrl string) string {
	state := ""
	if returnUrl != "" {
		state = url.QueryEscape(returnUrl)
	}
	return fmt.Sprintf("%s?client_id=%s&scope=%s&redirect_uri=%s&state=%s&response_type=code",
		config.Endpoint.AuthURL,
		url.QueryEscape(config.ClientID),
		url.QueryEscape(strings.Join(config.Scopes, " ")),
		url.QueryEscape(config.RedirectURL),
		// note that this does not provide any additional security against CSRF
		// i do not think, hovewer that its needed in our flow
		state)
}

const OauthProviderGoogle = "google"

func (req *Config) NewOauthCtx() *OauthCtx {
	ret := new(OauthCtx)
	ret.GoogleConfig = oauth2.Config{
		ClientID:     req.Google.ClientID,
		ClientSecret: req.Google.ClientSecret,
		RedirectURL:  req.Google.RedirectURL,
		Scopes:       req.Google.Scopes,
		Endpoint:     google.Endpoint,
	}
	return ret
}

type OauthUserInfo struct {
	Email      string `json:"email"`
	OauthID    string `json:"id"`
	IsVerified bool   `json:"verified_email"`
	Provider   string
}

func (ui *OauthUserInfo) Validate() error {
	const hdr = "Validate OauthUserInfo: "
	if ui.Email == "" {
		return fmt.Errorf("%sInvalid email", hdr)
	}
	if !ui.IsVerified {
		return fmt.Errorf("%sUser is not verified", hdr)
	}
	if ui.OauthID == "" {
		return fmt.Errorf("%sInvalid oauth id", hdr)
	}
	return nil
}

/*
	will not validate response
*/
func OauthGoogleGetUserInfo(token *oauth2.Token) (*OauthUserInfo, error) {
	res, err := http.Get(
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token=" +
			url.QueryEscape(token.AccessToken))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ret := new(OauthUserInfo)

	if err = json.Unmarshal(resBody, ret); err != nil {
		return nil, err
	}

	ret.Provider = OauthProviderGoogle

	return ret, nil
}
