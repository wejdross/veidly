package user

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/mail"
	"net/url"
	"sport/api"
	"sport/helpers"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	LangEn         string = "en"
	LangPl         string = "pl"
	LangDefault    string = LangEn
	CountryDefault string = "PL"
)

type NamedUrl struct {
	Name   string
	Url    string
	Avatar string
}

/*
	User data

	when modifying dont forget to update:
		func (ud *UserData) Value()
		func (ud *UserData) Scan(value interface{})
	(defined below)
	those functions are used in other modules
*/
type UserData struct {
	/*
		this may be full name or just nick
	*/
	Name string

	/*
		format is:
		ISO 639-1 language code
		en
		pl
		...
	*/
	Language string
	/*
		format is:
		iso 3166-1 alpha-2 country code
		US
		PL
		...
	*/
	Country string

	Urls []NamedUrl

	AboutMe string
}

// used to convert golang struct into db composite type representation
func (ud *UserData) Value() (driver.Value, error) {
	return json.Marshal(ud)
}

// used to convert DB composite type representation into golang struct
func (ud *UserData) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), ud)
}

func validateUrl(u, host string) error {
	pu, err := url.ParseRequestURI(u)
	if err != nil {
		return err
	}
	if host != "" {
		if (pu.Host != host) && (pu.Host != "www."+host) {
			return fmt.Errorf("invalid host: %s", pu.Host)
		}
	}
	return nil
}

func (ud *UserData) Validate(ctx *Ctx, edit bool) error {
	c := ctx.Config
	if len(ud.Name) > UserNameMaxLength {
		return fmt.Errorf("validate UserData: invalid firstName")
	}
	if utf8.RuneCountInString(ud.AboutMe) > 250 {
		return fmt.Errorf("validate UserData: invalid AboutMe")
	}
	if err := c.ValidateCountry(ud.Country); err != nil {
		return err
	}
	nix := make(map[string]struct{})
	uix := make(map[string]struct{})
	for i := range ud.Urls {
		if ud.Urls[i].Name == "" {
			return fmt.Errorf("url at index %d has no name", i)
		}
		if err := validateUrl(ud.Urls[i].Url, ""); err != nil {
			return err
		}
		if _, ok := nix[ud.Urls[i].Name]; ok {
			return fmt.Errorf("duplicate name")
		} else {
			nix[ud.Urls[i].Name] = struct{}{}
		}
		if _, ok := uix[ud.Urls[i].Url]; ok {
			return fmt.Errorf("duplicate url")
		} else {
			nix[ud.Urls[i].Url] = struct{}{}
		}
		// check if avatar is one of our supported avatars, or set it to default to avoid problems
		switch ud.Urls[i].Avatar {
		case "default":
		case "facebook":
		case "instagram":
		case "twitter":
		case "youtube":
		default:
			ud.Urls[i].Avatar = "default"
		}
	}

	if ud.Language != "" || !edit {
		if !ctx.LangCtx.ValidateApiLang(ud.Language) {
			ud.Language = LangDefault
			// return fmt.Errorf("unkown Lang: %s", ud.Language)
		}
	}

	return nil
}

type UserRequest struct {
	Email    string
	Password string
	UserData
}

func (ur *UserRequest) Validate(c *Ctx, withPassword bool) error {
	if len(ur.Email) == 0 || len(ur.Email) > UserNameMaxLength {
		return fmt.Errorf("Validate UserRequest: invalid email")
	}
	if withPassword {
		if len(ur.Password) == 0 || len(ur.Password) > UserPassMaxLength {
			return fmt.Errorf("Validate UserRequest: invalid password")
		}
		if err := ValidatePassword(ur.Password); err != nil {
			return err
		}
	}
	return ur.UserData.Validate(c, false)
}

func (ui *OauthUserInfo) NewUser() *User {
	return &User{
		Enabled: true,
		PrivUserInfo: PrivUserInfo{
			Email: ui.Email,
			PubUserInfo: PubUserInfo{
				UserData: UserData{
					Language: LangDefault,
					Country:  CountryDefault,
				},
			},
			ContactData: ContactData{
				Email: ui.Email,
			},
			CreatedOn:     time.Now().In(time.UTC),
			OauthProvider: ui.Provider,
		},
		ID:      uuid.New(),
		OauthID: ui.OauthID,
	}
}

func (ur *UserRequest) NewUser(mfa bool) (*User, error) {
	h, err := api.GetHash(ur.Password)
	if err != nil {
		return nil, err
	}

	ret := &User{
		PrivUserInfo: PrivUserInfo{
			Email: ur.Email,
			PubUserInfo: PubUserInfo{
				UserData: ur.UserData,
			},
			AccessFailed: 0,
			CreatedOn:    time.Now().In(time.UTC),
			ContactData: ContactData{
				Email: ur.Email,
			},
			OauthProvider: "",
		},
		Passh: h,
		// by default user is disabled onbce created
		//	- user must complete 2fa or oauth flow to become enabled.
		Enabled: false,
		ID:      uuid.New(),
		OauthID: "",
	}

	// [this is optional and may be commented out]
	// set default language
	if ret.UserData.Language == "" {
		ret.UserData.Language = LangDefault
	}

	if mfa {
		ret.MFAToken = helpers.GetUniqueToken()
	} else {
		ret.Enabled = true
	}

	return ret, nil
}

// DDL datatype
type User struct {
	ID              uuid.UUID
	Passh           []byte
	Enabled         bool
	MFAToken        string
	ForgotPassToken string
	OauthID         string
	PrivUserInfo
}

// Data which other users are allowed to read about this user
type PubUserInfo struct {
	UserData
	AvatarRelpath string
	AvatarUrl     string
}

const PubUserInfoSelectColsFmt = `
	coalesce(%s.user_data, '{}'::jsonb)
	,coalesce(%s.avatar_relpath, '')`

func GetPubUserInfoSelectColsFmt(prefix string) string {
	return fmt.Sprintf(PubUserInfoSelectColsFmt, prefix, prefix)
}

func (p *PubUserInfo) ScanFields() []interface{} {
	return []interface{}{
		&p.UserData,
		&p.AvatarRelpath,
	}
}

func (u *PubUserInfo) PostprocessAfterDbScan(c *Ctx) {
	if u.AvatarRelpath != "" {
		u.AvatarUrl = c.Static.Config.AppendToBaseUrl(u.AvatarRelpath)
	} else {
		u.AvatarUrl = ""
	}
}

type ContactData struct {
	Email string
	Phone string
	Share bool
}

func GetContactDataSelectCols(prefix string) string {
	return fmt.Sprintf(`
		coalesce(%s.contact_data, '{}'::jsonb)
	`, prefix)
}

func (p *ContactData) ScanFields() []interface{} {
	return []interface{}{
		&p,
	}
}

func (c *ContactData) Validate() error {
	if c.Email != "" {
		if len(c.Email) > 256 {
			return fmt.Errorf("validate contactdata: email too long")
		}
		if _, err := mail.ParseAddress(c.Email); err != nil {
			return fmt.Errorf("validate contactData email: %v ", err)
		}
	}

	if c.Phone != "" {
		if len(c.Email) > 32 {
			return fmt.Errorf("validate contactdata: phone too long")
		}
	}

	return nil
}

func (cd *ContactData) Value() (driver.Value, error) {
	return json.Marshal(cd)
}

func (cd *ContactData) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), cd)
}

// Data which user is allowed to read about himself
type PrivUserInfo struct {
	Email         string
	CreatedOn     time.Time
	AccessFailed  int
	OauthProvider string
	ContactData   ContactData
	PubUserInfo
}
