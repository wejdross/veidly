package instr

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sport/lang"
	"sport/user"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CardInfo struct {
	CardNumber  string
	ExpiryMonth string
	ExpiryYear  string
	HolderName  string
	Cvc         string
}

type ProcessedCardInfo struct {
	CardRefID      string
	CardBrand      string
	CardHolderName string
	CardSummary    string
}

/* fields which instructor can freely update*/
type InstructorRequest struct {
	Tags            []string
	YearExp         int
	Disabled        bool
	KnownLangs      []string
	ProfileSections ProfileSectionArr
	InvoiceLines    []string
}

func (ir *InstructorRequest) Validate(ctx *Ctx) error {
	if len(ir.Tags) > 5 {
		return fmt.Errorf("invalid tags")
	}
	for _, t := range ir.Tags {
		if len(t) == 0 || len(t) > 96 {
			return fmt.Errorf("invalid tag")
		}
	}
	if ir.YearExp != 0 {
		thisYear := time.Now().In(time.UTC).Year()
		if ir.YearExp > thisYear || (ir.YearExp+80) < thisYear {
			return fmt.Errorf("invalid year Exp")
		}
	}

	if len(ir.KnownLangs) > lang.MaxLangs {
		return fmt.Errorf("suspicius polyglot detected")
	}

	for i := range ir.KnownLangs {
		lang := ir.KnownLangs[i]
		if !ctx.Lang.ValidateUserLang(lang) {
			return fmt.Errorf("unknown KnownLangs[%d]: %s", i, lang)
		}
	}

	if len(ir.InvoiceLines) > 5 {
		return fmt.Errorf("invalid invoicing lines")
	}

	for i := range ir.InvoiceLines {
		line := ir.InvoiceLines[i]
		if len(line) > 60 {
			return fmt.Errorf("invalid invoice line at index: %d", i)
		}
	}

	return nil
}

func (ir *InstructorRequest) NewInstructor(userID uuid.UUID) *Instructor {
	return &Instructor{
		PubInstructorInfo: PubInstructorInfo{
			ID:                uuid.New(),
			UserID:            userID,
			CreatedOn:         time.Now().In(time.UTC),
			InstructorRequest: *ir,
		},
	}
}

type InstructorRefunds map[int64]int

// this is 1:1 representation of database model
type Instructor struct {
	CardInfo         ProcessedCardInfo
	Refunds          InstructorRefunds
	QueuedPayoutCuts int
	PubInstructorInfo
}

type InstructorWithUser struct {
	Instructor
	UserInfo    user.PubUserInfo
	ContactData user.ContactData
}

func InstructorWithUserCols() string {
	return `
		i.id,
		i.user_id,
		i.created_on,
		i.tags,
		i.year_exp,
		i.known_langs,
		i.disabled,
		i.refunds,
		i.queued_payout_cuts,
		u.user_data,
		u.avatar_relpath,
		u.contact_data,
		i.card_ref_id,
		i.card_brand,
		i.card_holder_name,
		i.card_summary,
		i.bg_img_path,
		i.extra_img_paths,
		i.profile_sections,
		i.invoice_lines
	`
}

func (ir *InstructorWithUser) ScanFields() []interface{} {
	return []interface{}{
		&ir.ID,
		&ir.UserID,
		&ir.CreatedOn,
		(*pq.StringArray)(&ir.Tags),
		&ir.YearExp,
		(*pq.StringArray)(&ir.KnownLangs),
		&ir.Disabled,
		&ir.Refunds,
		&ir.QueuedPayoutCuts,
		//
		&ir.UserInfo.UserData,
		&ir.UserInfo.AvatarRelpath,
		&ir.ContactData,
		//
		&ir.CardInfo.CardRefID,
		&ir.CardInfo.CardBrand,
		&ir.CardInfo.CardHolderName,
		&ir.CardInfo.CardSummary,
		//
		&ir.BgImgPath,
		(*pq.StringArray)(&ir.ExtraImgPaths),
		&ir.ProfileSections,
		(*pq.StringArray)(&ir.InvoiceLines),
	}
}

// used to convert golang struct into db composite type representation
func (ud *InstructorRefunds) Value() (driver.Value, error) {
	return json.Marshal(ud)
}

// used to convert DB composite type representation into golang struct
func (ud *InstructorRefunds) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), ud)
}

// information about instructor which anyone can read
type PubInstructorInfo struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID
	CreatedOn time.Time
	InstructorRequest
	BgImgPath     string
	BgImgUrl      string
	ExtraImgPaths []string
	ExtraImgUrls  []string
}

func (pii *PubInstructorInfo) PostprocessAfterDbScan(c *Ctx) {
	if pii.BgImgPath == "" {
		pii.BgImgUrl = ""
	} else {
		pii.BgImgUrl = c.Static.Config.AppendToBaseUrl(pii.BgImgPath)
	}
	pii.ExtraImgUrls = make([]string, len(pii.ExtraImgPaths))
	for i := range pii.ExtraImgPaths {
		pii.ExtraImgUrls[i] = c.Static.Config.AppendToBaseUrl(pii.ExtraImgPaths[i])
	}
}

type PubInstructorWithUser struct {
	PubInstructorInfo
	UserInfo user.PubUserInfo
	Config   ConfigState
}

func PubInstructorWithUserCols(uprefix string) string {
	return fmt.Sprintf(`
		i.id,
		i.user_id,
		i.created_on,
		i.tags,
		i.year_exp,
		i.known_langs,
		i.disabled,
		%s.user_data,
		%s.avatar_relpath
	`, uprefix, uprefix)
}

func (ir *PubInstructorWithUser) ScanFields() []interface{} {
	return []interface{}{
		&ir.ID,
		&ir.UserID,
		&ir.CreatedOn,
		(*pq.StringArray)(&ir.Tags),
		&ir.YearExp,
		(*pq.StringArray)(&ir.KnownLangs),
		&ir.Disabled,
		//
		&ir.UserInfo.UserData,
		&ir.UserInfo.AvatarRelpath,
	}
}

type InstructorCert struct {
	ID           uuid.UUID
	InstructorID uuid.UUID
}

type DeletedInstructor struct {
	Email            string
	Refunds          InstructorRefunds
	QueuedPayoutCuts int
}

func (ctx *Ctx) ProcessCardInfo(c *CardInfo) (*ProcessedCardInfo, error) {

	pci := new(ProcessedCardInfo)
	pci.CardHolderName = c.HolderName
	pci.CardRefID = uuid.NewString()

	if ctx.Adyen == nil || true {
		pci.CardBrand = "test"
		pci.CardSummary = "1111"
		return pci, nil
	}

	/*

		r := adyen.PaymentRequest{
			Amount: adyen.Amount{
				Currency: "PLN", // this has to be verified and implemented (?) to the instructor
				Value:    0,
			},
			PaymentMethod: adyen.AdyenCardDetails{
				EncryptedSecurityCode: c.Cvc,
				EncryptedExpiryMonth:  c.ExpiryMonth,
				EncryptedExpiryYear:   c.ExpiryYear,
				EncryptedCardNumber:   c.CardNumber,
				HolderName:            c.HolderName,
				Type:                  "scheme",
			},
			Reference:                uuid.NewString(),
			MerchantAccount:          ctx.Adyen.Config.MerchAcc,
			ShopperReference:         pci.CardRefID,
			EnablePayOut:             "True",
			ShopperInteraction:       "Ecommerce",
			StorePaymentMethod:       true,
			RecurringProcessingModel: "CardOnFile",
		}

		res, err := ctx.Adyen.AdyenPayments(&r)
		if err != nil {
			return nil, err
		}

		var expectedResCode adyen.ResultCode = "Authorised"

		if *res.ResultCode != expectedResCode {
			return nil,
				fmt.Errorf("invalid resultCode, expected %s got %s", expectedResCode, *res.ResultCode)
		}

		pe := res.AdditionalData["payoutEligible"]
		switch pe {
		case "Y":
			fallthrough
		case "D":
			break
		default:
			return nil, fmt.Errorf("card is not eligible for payout")
		}

		pm := res.AdditionalData["paymentMethod"]
		if pm == "" {
			return nil, fmt.Errorf("invalid paymentMethod")
		}

		cs := res.AdditionalData["cardSummary"]
		if cs == "" {
			return nil, fmt.Errorf("invalid cardSummary")
		}

		// this may not be neccessary
		fa := res.AdditionalData["fundsAvailability"]
		switch fa {
		case "I":
			break
		default:
			return nil, fmt.Errorf("fast funds are not supported for this card")
		}

		pci.CardBrand = pm
		pci.CardSummary = cs

		return pci, nil
	*/
	return nil, nil
}

type ConfigState uint32

const (
	CS_NoName    ConfigState = 1
	CS_NoCard    ConfigState = 2
	CS_Disabled  ConfigState = 4
	CS_NoContact ConfigState = 8
	CS_NoAboutMe ConfigState = 16
	CS_NoInvoice ConfigState = 32
)

/// LE only
func StateToStr(s ConfigState) string {
	var one uint32 = 1
	fb := strings.Builder{}
	wrote := false
	for i := 0; i < 32; i++ {
		if (uint32(s) & one) == 1 {
			if wrote {
				fb.WriteString(fmt.Sprintf(",%d", 1<<i))
			} else {
				wrote = true
				fb.WriteString(fmt.Sprintf("%d", 1<<i))
			}
		}
		s >>= 1
	}
	return fb.String()
}

func GetInstrConfig(i *InstructorWithUser) ConfigState {

	var res ConfigState = 0

	if strings.TrimSpace(i.UserInfo.Name) == "" {
		res |= CS_NoName
	}

	// if i.CardInfo.CardRefID == "" {
	// 	res |= CS_NoCard
	// }

	if i.Instructor.InstructorRequest.Disabled {
		res |= CS_Disabled
	}

	// if strings.TrimSpace(i.ContactData.Email) == "" &&
	// 	strings.TrimSpace(i.ContactData.Phone) == "" {
	// 	res |= CS_NoContact
	// }

	if strings.TrimSpace(i.UserInfo.AboutMe) == "" {
		res |= CS_NoAboutMe
	}

	// if len(i.InvoiceLines) == 0 {
	// 	res |= CS_NoInvoice
	// }

	return res
}

// verify that instructor account is configured based on their USER id
// return 0, nil if instructor is configured and there are no problems
// return !=0, nil if instructor is not configured
// return 0, err if something went bad (instructor doesnt exist, so on)
func (ctx *Ctx) GetInstrConfigByID(userID uuid.UUID) (ConfigState, error) {

	i, err := ctx.DalReadInstructor(userID, UserID)
	if err != nil {
		return 0, err
	}

	return GetInstrConfig(i), nil
}

func (i *InstructorWithUser) ToInfo() *PubInstructorWithUser {
	return &PubInstructorWithUser{
		PubInstructorInfo: i.PubInstructorInfo,
		UserInfo:          i.UserInfo,
		Config:            GetInstrConfig(i),
	}
}
