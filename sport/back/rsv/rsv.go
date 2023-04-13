package rsv

import (
	"fmt"
	"sport/adyen"
	"sport/adyen_sm"
	"sport/dc"
	"sport/helpers"
	"sport/instr"
	"sport/train"
	"sport/user"
	"time"

	"github.com/google/uuid"
)

// used to create helper structures in search (there is no need for full Rsv)
type ShortenedRsv struct {
	DateStart  time.Time
	DateEnd    time.Time
	TrainingID *uuid.UUID
	ID         uuid.UUID
	Groups     train.GroupArr
}

// reservation data readable for owners (instructors and their clients)
type Rsv struct {
	ID           uuid.UUID
	InstructorID *uuid.UUID

	TrainingID *uuid.UUID
	Training   *train.Training
	Occ        train.OccWithSecondary
	Groups     train.GroupArr

	UserID   *uuid.UUID
	UserInfo user.PubUserInfo

	DateStart time.Time
	DateEnd   time.Time

	IsConfirmed        bool
	InstructorDecision adyen_sm.InstructorDecision
	// current state of reservation
	State adyen_sm.State
	// this is timeout for current state
	SmTimeout time.Time
	CreatedOn time.Time

	LinkID  string
	LinkUrl string

	// current state change retry count
	SmRetries int
	IsActive  bool

	UserContactData user.ContactData

	Dc         *dc.Dc
	DcRollback bool

	UseUserAcc bool
}

// this is DDL reservation representation
type DDLRsv struct {
	Rsv

	// payment reference
	OrderID string

	SmCache adyen_sm.SmCache

	ProcessingFee  int
	SplitIncomeFee int
	SplitPayout    int
	// this will be refunded to user
	RefundAmount int

	AccessToken uuid.UUID

	QrConfirmed bool
}

func (ctx *Ctx) RsvTokenUrl(r *DDLRsv) string {
	return fmt.Sprintf(ctx.Config.RsvDetailsUrlFmt, r.AccessToken, "token")
}

type DDLRsvWithInstr struct {
	DDLRsv
	Instructor *instr.InstructorWithUser
}

type RsvWithInstr struct {
	Rsv
	Instructor *instr.PubInstructorWithUser
}

type ApiReservationRequest struct {

	// target training ID
	TrainingID uuid.UUID
	// reservation datetime
	Occurrence time.Time

	/*
		this data is provided so instructor can see personal information about customer
	*/
	UserData    user.UserData
	ContactData user.ContactData

	UseSavedData bool

	//
	NoRedirect bool

	DcID *uuid.UUID
}

func ValidateUserDataForRsv(ui *user.PubUserInfo) error {
	const hdr = "Validate UserDataForRsv: "
	if ui.Name == "" {
		return fmt.Errorf("%sUser name must be provided", hdr)
	}
	if ui.Language == "" {
		return fmt.Errorf("%sUser language must be provided", hdr)
	}
	return nil
}

func ValidateContactDataForRsv(c *user.ContactData) error {
	return nil
}

func (c *ApiReservationRequest) Validate(ctx *Ctx) error {
	if c.TrainingID == uuid.Nil {
		return fmt.Errorf("Validate CreateReservationRequest: invalid trainingID")
	}
	if c.Occurrence.IsZero() {
		return fmt.Errorf("Validate CreateReservationRequest: invalid occurrence")
	}
	return nil
}

type AvReason uint32

const (
	// too many people registered
	AvCapacityReached AvReason = 1

	// training or occ changed or does not exist
	AvOrphaned AvReason = 2

	// training is disabled
	AvDisabled AvReason = 4

	// too late or too soon for this rsv to happen
	AvTiming AvReason = 8

	// request date doesnt match training / occ
	AvOccDoesntMatch AvReason = 16

	// Group limit reached
	AvGroupPeopleLimit   AvReason = 32
	AvGroupTrainingLimit AvReason = 64

	// instructor is on vacation
	AvOnVacation AvReason = 128

	//
)

type RsvAvailabilityInfo struct {
	IsAvailable bool

	// start time of this rsv
	DateStart time.Time

	// end time of this rsv
	DateEnd time.Time

	// optional error which explains why is this rsv not available
	// Error error
	Reason AvReason

	// occurrence which matched
	Occ *train.OccWithSecondary
}

type ReservationRequest struct {

	// target training which is to be copied into rsv
	// cannot be nil
	Tr *train.TrainingWithJoins

	// availability info about this rsv
	// cannot be nil
	AvInfo RsvAvailabilityInfo

	// response from adyen payment link creation
	// cannot be nil
	AdyenRes *adyen.CreatePaymentLinkResponse

	// timestamp when was rsv created
	DateOfRsv time.Time

	// infos about client
	UserID          *uuid.UUID
	UserContactData user.ContactData
	UserInfo        user.PubUserInfo

	PricingInfo RsvPricingInfo

	UseUserAcc bool
}

func (rr *ReservationRequest) VError(m string) error {
	return fmt.Errorf("Validate ReservationRequest: %s", m)
}

func (rr *ReservationRequest) NewReservation(ctx *Ctx, id, at *uuid.UUID) *DDLRsv {
	if id == nil {
		_id := uuid.New()
		id = &_id
	}
	if at == nil {
		_at := uuid.New()
		at = &_at
	}

	state := adyen_sm.LinkExpress
	var smto time.Time
	if ctx.NoPaymentFlow {
		state = adyen_sm.Hold
		// to avoid triggering state machine
		smto = time.Date(5000, 12, 31, 24, 59, 59, 0, time.UTC)
	} else {
		smto = rr.DateOfRsv.Add(time.Duration(ctx.Config.LinkExpressExpire))
	}

	now := time.Now().In(time.UTC)
	ret := &DDLRsv{
		Rsv: Rsv{
			DateStart:          rr.AvInfo.DateStart,
			DateEnd:            rr.AvInfo.DateEnd,
			ID:                 *id,
			UserID:             rr.UserID,
			InstructorID:       &rr.Tr.Training.InstructorID,
			UserInfo:           rr.UserInfo,
			Training:           &rr.Tr.Training,
			Occ:                *rr.AvInfo.Occ,
			Groups:             rr.Tr.Groups,
			TrainingID:         &rr.Tr.Training.ID,
			CreatedOn:          now,
			InstructorDecision: adyen_sm.Unset,
			LinkID:             rr.AdyenRes.ID,
			LinkUrl:            rr.AdyenRes.Url,
			IsActive:           true,
			IsConfirmed:        false,
			UserContactData:    rr.UserContactData,
			Dc:                 rr.PricingInfo.Dc,
			UseUserAcc:         rr.UseUserAcc,
			State:              state,
			SmTimeout:          smto,
		},
		// set in notification handler
		OrderID:        "",
		QrConfirmed:    false,
		SplitPayout:    rr.PricingInfo.SplitPayout,
		SplitIncomeFee: rr.PricingInfo.SplitIncomeFee,
		RefundAmount:   rr.PricingInfo.RefundAmount,
		ProcessingFee:  rr.PricingInfo.ProcessingFee,
		AccessToken:    *at,
		SmCache:        make(adyen_sm.SmCache),
	}

	return ret
}

func (r *Rsv) ToTrainingWithJoins() *train.TrainingWithJoins {
	x := train.TrainingWithJoins{}
	x.Training = *r.Training
	x.Groups = r.Groups
	if r.Dc != nil {
		x.Dcs = []dc.Dc{*r.Dc}
	}
	x.Occurrences = []train.OccWithSecondary{r.Occ}
	return &x
}

func ValidateRegisterTiming(
	ctx *Ctx,
	training *train.Training,
	occurrence time.Time,
	now *time.Time,
) bool {
	if now == nil {
		c := NowMin()
		now = &c
	}

	diff := occurrence.Sub(*now)
	if diff <= 0 {
		return false
	}

	if training.AllowExpress {
		if diff < time.Duration(ctx.Config.LinkExpressAtLeastBefore) {
			return false
		}
	} else {
		if training.ManualConfirm {
			if diff < time.Duration(ctx.Config.LinkManualAtLeastBefore) {
				return false
			}
		} else {
			if diff < time.Duration(ctx.Config.LinkAtLeastBefore) {
				return false
			}
		}
	}

	return true
}

type CanRegisterBuffers struct {
	Vacations   instr.VacationsGroupedByInstr
	RsvCount    RsvCount `json:"-"`
	GroupedRsvs GroupHash
}

var EmptyBuffers = CanRegisterBuffers{}

/*
	validate that

	1. occurrence matches allowed occ for training
	2. capacity for this session is not exceeded
	3. is not too late / too soon to register
	4. instructor is not on vacation

*/
func (ctx *Ctx) CanRegister(
	t *train.TrainingWithJoins,
	start time.Time,
	now *time.Time,
	// when u run this function in loop you should provide vacations and rsvGrp
	// to reuse resources
	buffers CanRegisterBuffers,
	// if any of the used buffers is nil, then error will be returned
	forceBuffers bool,
) (RsvAvailabilityInfo, error) {

	if now == nil {
		c := NowMin()
		now = &c
	}

	var res RsvAvailabilityInfo
	var err error
	res.DateStart = start

	if t.Training.Disabled {
		res.Reason |= AvDisabled
		return res, nil
	}

	// is not too late or too soon to make reservation
	if !ValidateRegisterTiming(ctx, &t.Training, start, now) {
		res.Reason |= AvTiming
		//res.Error = fmt.Errorf("too late to make reservation for this training")
		return res, nil
	}

	// occurrence matches
	res.Occ, res.DateEnd, err = t.ValidateOccurrence(start)
	if err != nil {
		//res.Error = err
		res.Reason = AvOccDoesntMatch
		return res, nil
	}

	instructorID := t.Training.InstructorID

	var groupIx GroupHash
	if buffers.GroupedRsvs == nil {
		if forceBuffers {
			return res, fmt.Errorf("CanRegister: GroupedRsvs buffer is nil")
		}
		// build ix
		ixRsvs, err := ctx.ReadRsvWithInstr(
			ReadRsvsArgs{
				DateRange: helpers.DateRange{
					Start: start,
					End:   res.DateEnd,
				},
				InstructorID: &instructorID,
			},
		)
		if err != nil {
			return res, err
		}
		groupIx = NewGroupHashFromRsvs(ixRsvs.Rsv)
	} else {
		groupIx = buffers.GroupedRsvs
	}

	for i := range t.Groups {
		v := groupIx[t.Groups[i].ID]
		for j := range v {
			if !overlaps(res.DateStart, res.DateEnd, v[j].DateStart, v[j].DateEnd) {
				continue
			}
			v[j].People++
			if _, ok := v[j].Trainings[t.Training.ID]; !ok {
				v[j].Trainings[t.Training.ID] = struct{}{}
			}
			if v[j].People > t.Groups[i].MaxPeople {
				//res.Error = fmt.Errorf("Group '%s' is full", t.Groups[i].Name)
				res.Reason = AvGroupPeopleLimit
				return res, nil
			}
			if len(v[j].Trainings) > t.Groups[i].MaxTrainings {
				//res.Error = fmt.Errorf("Group '%s' is full", t.Groups[i].Name)
				res.Reason = AvGroupTrainingLimit
				return res, nil
			}
		}
	}

	// instr on vacation
	if buffers.Vacations != nil {
		v := buffers.Vacations[instructorID]
		for i := range v {
			if overlaps(v[i].DateStart, v[i].DateEnd, res.DateStart, res.DateEnd) {
				//res.Error = fmt.Errorf("instructor is on vacation")
				res.Reason = AvOnVacation
				return res, nil
			}
		}
	} else {
		if forceBuffers {
			return res, fmt.Errorf("CanRegister: Vacations buffer is nil")
		}
		if vac, err := ctx.Instr.DalIsInstructorOnVacation(
			instructorID, start, res.DateEnd); err != nil {
			return res, err
		} else if vac {
			//res.Error = fmt.Errorf("instructor is on vacation")
			res.Reason = AvOnVacation
			return res, nil
		}
	}

	// capacity for this session is not exceeded
	if buffers.RsvCount != nil {
		key := RsvCountKey{
			TrainID:   t.Training.ID,
			DateStart: start.UnixNano(),
			DateEnd:   res.DateEnd.UnixNano(),
		}
		if x := buffers.RsvCount[key]; x >= t.Training.Capacity {
			//res.Error = fmt.Errorf("rsv is not available")
			res.Reason = AvCapacityReached
			return res, nil
		}
	} else {
		if forceBuffers {
			return res, fmt.Errorf("CanRegister: RsvCount buffer is nil")
		}
		tmp, err := ctx.ValidateCapacity(&t.Training, start, nil)
		if err != nil {
			return res, err
		}

		if !tmp {
			//res.Error = fmt.Errorf("rsv is not available")
			res.Reason = AvCapacityReached
			return res, nil
		}
	}

	res.IsAvailable = true
	//res.Error = nil
	return res, nil
}
