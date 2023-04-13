package schedule

import (
	"fmt"
	"sport/helpers"
	"sport/rsv"
	"sport/sub"
	"sport/train"
	"time"

	"github.com/google/uuid"
)

// provided by user
type ApiScheduleFlags int

const (
	NoRsvs ApiScheduleFlags = 1
)

// provided internally in system
type ScheduleFlags int

const (
	ForInstructor ScheduleFlags = 1
	OrphanRsv     ScheduleFlags = 2
)

type ScheduleRecord struct {
	Start              time.Time
	End                time.Time
	Count              int
	CountUnconfirmed   int
	IsAvailable        bool
	AvailabilityReason rsv.AvReason
	IsOrphaned         bool
	// if rsv is not orphaned this will contain info about current occ
	Occ *train.OccWithSecondary
	// if request is made by instructor this will be filled with reservations
	Reservations []*rsv.Rsv
	// if request is made by instructor this will be filled with possible subscriptions
	Subs []*sub.SubWithJoins
}

type TrainingSchedule struct {
	Schedule []ScheduleRecord
	*train.TrainingWithJoins
}

type ScheduleBuffers struct {
	rsv.CanRegisterBuffers
	GroupedSm sub.GroupedSm
}

func (ctx *Ctx) NewSchedule(
	trainings []*train.TrainingWithJoins,
	rsvRes []rsv.RsvWithInstr,
	subs []sub.SubWithJoins,
	requestRange helpers.DateRange,
	sflags ScheduleFlags,
	aflags ApiScheduleFlags,
	buffers ScheduleBuffers,
	forceBuffers bool,
) ([]TrainingSchedule, error) {

	var t *train.TrainingWithJoins
	var o *train.OccWithSecondary
	var scheduleRes = make([]TrainingSchedule, 0, 10)
	var adScheduleRes = make([]TrainingSchedule, 0, 10)
	var adIndex = make(map[string]int)
	matchedRsvs := make(map[uuid.UUID]struct{})
	var now = time.Now()
	var scheduleEntry TrainingSchedule

	for i := 0; i < len(trainings); i++ {
		t = trainings[i]
		scheduleEntry.TrainingWithJoins = t
		scheduleEntry.Schedule = make([]ScheduleRecord, 0, len(t.Occurrences))
		for j := 0; j < len(t.Occurrences); j++ {
			o = &t.Occurrences[j]
			if o.RepeatDays <= 0 {
				if !t.Training.DateStart.IsZero() && !t.Training.DateEnd.IsZero() {
					if !o.DateStart.After(t.Training.DateStart) || !o.DateEnd.Before(t.Training.DateEnd) {
						continue
					}
				}

				// check if ovelaps with request range and add it
				if !o.DateStart.After(requestRange.End) && !o.DateEnd.Before(requestRange.Start) {
					scheduleEntry.Schedule = append(scheduleEntry.Schedule, ScheduleRecord{
						Start:       o.DateStart,
						End:         o.DateEnd,
						Count:       0,
						IsAvailable: true,
						Occ:         &t.Occurrences[j],
					})
				}
			} else {
				duration := o.DateEnd.Sub(o.DateStart)

				trainingEnd := t.Training.DateEnd
				if trainingEnd.IsZero() || requestRange.End.Before(trainingEnd) {
					trainingEnd = requestRange.End
				}
				x := o.DateStart
				for !x.After(trainingEnd) {
					y := x.Add(duration)
					if !t.Training.DateStart.IsZero() && !t.Training.DateStart.IsZero() {
						if x.Before(t.Training.DateStart) || y.After(t.Training.DateEnd) {
							goto LE
						}
					}

					// check if ovelaps with request range and add it
					if !x.After(requestRange.End) && !y.Before(requestRange.Start) {
						scheduleEntry.Schedule = append(scheduleEntry.Schedule, ScheduleRecord{
							Start:       x,
							End:         y,
							Count:       0,
							IsAvailable: true,
							Occ:         &t.Occurrences[j],
						})
					}

				LE:
					x = x.Add(time.Hour * 24 * time.Duration(o.RepeatDays))
				}
			}
		}
		scheduleRes = append(scheduleRes, scheduleEntry)
	}

	if (sflags&ForInstructor) != 0 && len(subs) > 0 {

		ix := buffers.GroupedSm
		if ix == nil {
			if forceBuffers {
				return nil, fmt.Errorf("NewSchedule: GroupedSm was nil")
			}
			ix = sub.NewGroupedSm(subs)
		}

		// iterate over schedule trainings
		for i := 0; i < len(scheduleRes); i++ {
			sms := scheduleRes[i].TrainingWithJoins.Sms
			for j := 0; j < len(sms); j++ {

				// find matching subscriptions (with training sub models)
				ss, exists := ix[sms[j].ID]
				if !exists {
					continue
				}

				// iterate over schedule items and find subs that are valid
				for k := 0; k < len(scheduleRes[i].Schedule); k++ {
					s := &scheduleRes[i].Schedule[k]
					for z := 0; z < len(ss); z++ {
						isValid := ss[z].ValidStatus(s.Start)
						if isValid != sub.SubValid {
							continue
						}

						if s.Subs == nil {
							s.Subs = make([]*sub.SubWithJoins, 0, 3)
						}
						s.Subs = append(s.Subs, &ss[z])
					}
				}
			}
		}
	}

	for i := 0; i < len(scheduleRes); i++ {
		t := scheduleRes[i].Training
		for j := 0; j < len(scheduleRes[i].Schedule); j++ {
			s := &scheduleRes[i].Schedule[j]
			for z := 0; z < len(rsvRes); z++ {
				r := &rsvRes[z]
				// check if reservation belongs to this training
				if r.DateStart.Equal(s.Start) && r.DateEnd.Equal(s.End) && r.Training.ID == t.ID {

					matchedRsvs[r.ID] = struct{}{}

					if r.IsConfirmed {
						s.Count++
					} else if (sflags & ForInstructor) != 0 {
						s.CountUnconfirmed++
					}
					// if this data is requested by instructor then add users to reservation
					if (sflags&ForInstructor) != 0 && (aflags&NoRsvs) == 0 {
						if s.Reservations == nil {
							c := t.Capacity / 2
							if c == 0 {
								c = 1
							}
							s.Reservations = make([]*rsv.Rsv, 0, c)
						}
						s.Reservations = append(s.Reservations, &r.Rsv)
						// if s.Participants == nil {
						// 	s.Participants = make([]user.UserData, 0, 4)
						// }
						// s.Participants = append(s.Participants, r.UserData)
					}

					// if we exceeded count then mark this training as unavailable
					if s.Count >= t.Capacity {
						s.AvailabilityReason |= rsv.AvCapacityReached
						s.IsAvailable = false
					}

					// if dates overlap with this training but rsv belongs to different training
					// 		-> then mark this training as no longer available
				}
				//  else if r.DateStart.Before(s.End) && r.DateEnd.After(s.Start) && r.Training.ID != t.ID {
				// 	s.IsAvailable = false
				// }
			}

			// if requestRange.End.Day() == 29 && s.Start.Day() == 23 && s.IsAvailable {
			// 	fmt.Println(helpers.JsonMustSerializeFormatStr(ts[i].Groups))
			// 	fmt.Println(helpers.JsonMustSerializeFormatStr(crpb.GroupIx))
			// }

			// i really dont want to do this since this is very slow.
			// but i dont see any other way
			// ...
			// if training is still marked as available we have to make sure that it really is
			if s.IsAvailable {
				if c, err := ctx.Rsv.CanRegister(
					scheduleRes[i].TrainingWithJoins, s.Start, &now, buffers.CanRegisterBuffers, true,
				); err != nil {
					return nil, err
				} else if !c.IsAvailable {
					s.IsAvailable = false
					s.AvailabilityReason = c.Reason
				}
			}
		}
	}

	if (sflags & OrphanRsv) != 0 {

		/*
			add orphan reservations.
			(reservations which are assigned to no longer existing trainings / sessions)
		*/
		for i := 0; i < len(rsvRes); i++ {

			r := &rsvRes[i]

			// we dont need to bother with this rsv if it has been matched to training
			if _, ok := matchedRsvs[r.ID]; ok {
				continue
			}

			// verify that reservation overlaps with request date range
			if !(rsvRes[i].DateStart.Before(requestRange.End) && r.DateEnd.After(requestRange.Start)) {
				continue
			}

			tt := TrainingSchedule{
				TrainingWithJoins: r.ToTrainingWithJoins(),
				Schedule:          []ScheduleRecord{},
			}

			tsr := ScheduleRecord{
				Start:            r.DateStart,
				End:              r.DateEnd,
				Count:            0,
				CountUnconfirmed: 0,
				// orphaned sessions are always set as not available
				IsAvailable:        false,
				AvailabilityReason: rsv.AvOrphaned,
				IsOrphaned:         true,
				Occ:                &r.Occ,
			}

			if r.IsConfirmed {
				tsr.Count++
			} else if (sflags & ForInstructor) != 0 {
				tsr.CountUnconfirmed++
			}

			if (sflags&ForInstructor) != 0 && (aflags&NoRsvs) == 0 {
				// tsr.Participants = []user.UserData{
				// 	r.UserData,
				// }
				tsr.Reservations = append(tsr.Reservations, &r.Rsv)
			}

			//var foundOcc *Occurrence
			tsr.Occ = &r.Occ

			tt.Schedule = append(tt.Schedule, tsr)

			ix := fmt.Sprintf("%d-%d-%s-%s",
				r.DateStart.Unix(), r.DateEnd.Unix(), r.Training.ID.String(), tsr.Occ.ID.String())
			// add this rsv + training to response with "orphaned" flag
			//

			if pos, ok := adIndex[ix]; ok {
				adScheduleRes[pos].Schedule[0].Count += tsr.Count
				adScheduleRes[pos].Schedule[0].CountUnconfirmed += tsr.CountUnconfirmed
				// if len(tsr.Participants) > 0 {
				// 	adScheduleRes[pos].Schedule[0].Participants =
				// 		append(adScheduleRes[pos].Schedule[0].Participants, tsr.Participants...)
				// }
				if len(tsr.Reservations) > 0 {
					adScheduleRes[pos].Schedule[0].Reservations =
						append(adScheduleRes[pos].Schedule[0].Reservations, tsr.Reservations...)
				}
			} else {
				adScheduleRes = append(adScheduleRes, tt)
				adIndex[ix] = len(adScheduleRes) - 1
			}

		}
	}

	scheduleRes = append(scheduleRes, adScheduleRes...)

	return scheduleRes, nil
}

func ConvertRsvToSchedule(r *rsv.Rsv, requestRange helpers.DateRange) TrainingSchedule {
	ret := TrainingSchedule{}
	ret.TrainingWithJoins = r.ToTrainingWithJoins()
	ret.Schedule = []ScheduleRecord{
		{
			Start:            r.DateStart,
			End:              r.DateEnd,
			Count:            0,
			CountUnconfirmed: 0,
			IsAvailable:      false,
			IsOrphaned:       false,
			Occ:              &r.Occ,
			Reservations: []*rsv.Rsv{
				r,
			},
		},
	}
	return ret
}

func (ctx *Ctx) NewTrainingScheduleRecords(
	t *train.TrainingWithJoins,
	requestRange helpers.DateRange,
	buffers ScheduleBuffers,
) ([]ScheduleRecord, error) {

	if buffers.RsvCount == nil {
		return nil, fmt.Errorf("NewUserSchedule: RsvCount was nil")
	}

	var o *train.OccWithSecondary
	var now = time.Now()

	res := make([]ScheduleRecord, 0, len(t.Occurrences))
	for j := 0; j < len(t.Occurrences); j++ {
		o = &t.Occurrences[j]
		if o.RepeatDays <= 0 {
			if !t.Training.DateStart.IsZero() && !t.Training.DateEnd.IsZero() {
				if !o.DateStart.After(t.Training.DateStart) || !o.DateEnd.Before(t.Training.DateEnd) {
					continue
				}
			}

			// check if ovelaps with request range and add it
			if !o.DateStart.After(requestRange.End) && !o.DateEnd.Before(requestRange.Start) {
				res = append(res, ScheduleRecord{
					Start:       o.DateStart,
					End:         o.DateEnd,
					Count:       0,
					IsAvailable: true,
					Occ:         &t.Occurrences[j],
				})
			}
		} else {
			duration := o.DateEnd.Sub(o.DateStart)

			trainingEnd := t.Training.DateEnd
			if trainingEnd.IsZero() || requestRange.End.Before(trainingEnd) {
				trainingEnd = requestRange.End
			}
			x := o.DateStart
			for !x.After(trainingEnd) {
				y := x.Add(duration)
				if !t.Training.DateStart.IsZero() && !t.Training.DateStart.IsZero() {
					if x.Before(t.Training.DateStart) || y.After(t.Training.DateEnd) {
						goto LE
					}
				}

				// check if ovelaps with request range and add it
				if !x.After(requestRange.End) && !y.Before(requestRange.Start) {
					res = append(res, ScheduleRecord{
						Start:       x,
						End:         y,
						Count:       0,
						IsAvailable: true,
						Occ:         &t.Occurrences[j],
					})
				}

			LE:
				x = x.Add(time.Hour * 24 * time.Duration(o.RepeatDays))
			}
		}
	}

	for j := 0; j < len(res); j++ {
		s := &res[j]
		rc := buffers.RsvCount[rsv.RsvCountKey{
			TrainID:   t.Training.ID,
			DateStart: s.Start.UnixNano(),
			DateEnd:   s.End.UnixNano(),
		}]
		s.Count += rc

		// if we exceeded count then mark this training as unavailable
		if s.Count >= t.Training.Capacity {
			s.AvailabilityReason |= rsv.AvCapacityReached
			s.IsAvailable = false
		}

		if s.IsAvailable {
			if c, err := ctx.Rsv.CanRegister(
				t, s.Start, &now, buffers.CanRegisterBuffers, true,
			); err != nil {
				return nil, err
			} else if !c.IsAvailable {
				s.IsAvailable = false
				s.AvailabilityReason = c.Reason
			}
		}
	}

	return res, nil
}

/*
	this function will return data with the same ordering as trainings
*/
func (ctx *Ctx) NewUserSchedule(
	trainings []*train.TrainingWithJoins,
	requestRange helpers.DateRange,
	buffers ScheduleBuffers,
	forceBuffers bool,
) ([]TrainingSchedule, error) {
	var scheduleRes = make([]TrainingSchedule, 0, 10)

	if buffers.RsvCount == nil {
		if forceBuffers {
			return nil, fmt.Errorf("NewUserSchedule: RsvCount was nil")
		}
		var err error
		buffers.RsvCount, err = ctx.Rsv.NewDalRsvCount(&requestRange)
		if err != nil {
			return nil, err
		}
	}

	for i := 0; i < len(trainings); i++ {
		scheduleRecords, err := ctx.NewTrainingScheduleRecords(trainings[i], requestRange, buffers)
		if err != nil {
			return nil, err
		}
		scheduleRes = append(scheduleRes, TrainingSchedule{
			Schedule:          scheduleRecords,
			TrainingWithJoins: trainings[i],
		})
	}

	return scheduleRes, nil
}
