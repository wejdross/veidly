package schedule

import (
	"fmt"
	"sport/helpers"
	"sport/instr"
	"sport/rsv"
	"sport/sub"
	"sport/train"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctx *Ctx) HandlerGetRsv() gin.HandlerFunc {
	return func(g *gin.Context) {
		rsvs := ctx.Rsv.GetRsvsHandlerContent(g)
		if rsvs == nil {
			return
		}
		var scheduleRes = make([]TrainingSchedule, 0, 10)
		var schedItem TrainingSchedule

		//tix := make(map[uuid.UUID]int)

		for i := 0; i < len(rsvs.Rsv); i++ {
			r := &rsvs.Rsv[i]
			schedItem = ConvertRsvToSchedule(&r.Rsv, helpers.DateRange{
				Start: r.DateStart,
				End:   r.DateEnd})
			scheduleRes = append(scheduleRes, schedItem)
		}

		g.AbortWithStatusJSON(200, scheduleRes)
	}
}

/*
GET /api/rsv/schedule

DESC
	get trainings occurrences along with busy/free info

REQUEST
	- query string
		- start, end
			date range from which to generate schedule. must not be empty
			note that max allowed date range is set by configuration in field:
				- max_schedule_period_in_days
		- instructor_id
			user must either provide instructor_id or send authorization token (and be instructor)
		- training_id
			filter option
		- flags
			bit flags which specify extra behaviour (see ApiScheduleFlags)

RETURNS
	- on success application/json of type []TrainingSchedule
	- status code:
		- 200 on success
		- error status code
			- 400 when validation of request params fails
			- 500 on unexpected error
			- 401 on authorization failure (if instructor_id is not specified)
			- 404 if user is not an instructor (if instructor_id is not specified)
*/
func (ctx *Ctx) HandlerGetSched() gin.HandlerFunc {
	return func(g *gin.Context) {

		var err error
		var status = 400
		var scheduleRes = make([]TrainingSchedule, 0, 10)
		var instructorID uuid.UUID
		var trainingIDptr *uuid.UUID
		var smID *uuid.UUID
		// var adIndex = make(map[string]int)
		// matchedRsvs := make(map[uuid.UUID]struct{})
		var isForInstructor = false
		var trainings []*train.TrainingWithJoins
		var r *rsv.RsvWithInstrPagination
		var subs []sub.SubWithJoins
		var nc = make(chan error, 2)
		var flags ApiScheduleFlags
		var vix instr.VacationsGroupedByInstr
		var rcix rsv.RsvCount
		var grix rsv.GroupHash

		var requestRange helpers.DateRange
		if requestRange, err = helpers.DateRangeFromQueryString(g); err != nil {
			goto end
		}

		if !requestRange.IsNotZero() {
			err = fmt.Errorf("User Must provide date range")
			goto end
		}

		if f := g.Query("flags"); f != "" {
			if x, err := strconv.Atoi(f); err != nil {
				goto end
			} else {
				flags = ApiScheduleFlags(x)
			}
		}

		if requestRange.End.Sub(requestRange.Start) > time.Duration(ctx.Config.MaxSchedulePeriodInDays)*time.Hour*24 {
			err = fmt.Errorf("Requested schedule period (%d days) exceeded allowed period (%d days)",
				requestRange.End.Sub(requestRange.Start)/time.Hour/24,
				ctx.Config.MaxSchedulePeriodInDays)
			goto end
		}

		if tmp := g.Query("instructor_id"); tmp == "" {
			// try to get instructorID from authorization token
			var userID uuid.UUID
			if userID, err = ctx.Api.AuthorizeUserFromCtx(g); err != nil {
				status = 401
				goto end
			} else {
				if instructorID, err = ctx.Instr.DalReadInstructorID(userID); err != nil {
					if helpers.IsENF(err) {
						status = 404
					} else {
						status = 500
					}
					goto end
				}
				isForInstructor = true
			}
		} else {
			if instructorID, err = uuid.Parse(tmp); err != nil {
				goto end
			}
			if instructorID == uuid.Nil {
				err = fmt.Errorf("Instructor ID provided by user was nil")
				goto end
			}
		}

		if tmp := g.Query("training_id"); tmp != "" {
			var trainingID uuid.UUID
			if trainingID, err = uuid.Parse(tmp); err != nil {
				status = 400
				goto end
			} else {
				trainingIDptr = &trainingID
			}
		}

		if tmp := g.Query("smID"); tmp != "" {
			var _smID uuid.UUID
			if _smID, err = uuid.Parse(tmp); err != nil {
				status = 400
				goto end
			} else {
				smID = &_smID
			}
		}

		status = 500

		// T1
		go func() {
			var err error
			trainings, err = ctx.Train.DalReadTrainings(
				train.DalReadTrainingsRequest{
					InstructorID:  &instructorID,
					DateRange:     requestRange,
					TrainingID:    trainingIDptr,
					SmID:          smID,
					WithOccs:      true,
					WithGroups:    true,
					WithSubModels: true,
					// discount codes can only be visible for instructor
					WithDcs: isForInstructor,
				},
			)
			nc <- err
		}()

		// T2
		go func() {
			var err error
			var _rsv *rsv.RsvWithInstrPagination
			_rsv, err = ctx.Rsv.ReadRsvWithInstr(
				rsv.ReadRsvsArgs{
					//Pagination:     nil,
					DateRange:      requestRange,
					WithInstructor: false,
					TrainingID:     trainingIDptr,
					InstructorID:   &instructorID,
				},
			)
			if err != nil {
				nc <- err
				return
			}
			if _rsv == nil {
				nc <- fmt.Errorf("_rsv was nil, that REALLY shouldnt have happened")
				return
			}

			if trainingIDptr == nil {
				r = _rsv
			} else {
				r = &rsv.RsvWithInstrPagination{
					Pagination: _rsv.Pagination,
					Rsv:        make([]rsv.RsvWithInstr, 0, len(_rsv.Rsv)/2),
				}
				for i := range _rsv.Rsv {
					if _rsv.Rsv[i].TrainingID != nil && *_rsv.Rsv[i].TrainingID == *trainingIDptr {
						r.Rsv = append(r.Rsv, _rsv.Rsv[i])
					}
				}
			}

			rcix = rsv.NewRsvCount(r.Rsv)
			grix = rsv.NewGroupHashFromRsvs(_rsv.Rsv)

			nc <- nil
		}()

		// T3
		go func() {

			if !isForInstructor {
				nc <- nil
				return
			}

			var err error
			subs, err = ctx.Sub.DalReadSubWithJoins(
				sub.ReadSubRequest{
					DateRange:    requestRange,
					InstructorID: &instructorID,
				},
			)
			if err != nil {
				nc <- err
				return
			}

			nc <- nil
		}()

		// T4
		go func() {
			var err error
			vix, err = ctx.Instr.DalCreateVacationIndex([]string{instructorID.String()})
			nc <- err
		}()

		{
			const nothr = 4
			var errs [nothr]error
			for i := 0; i < nothr; i++ {
				errs[i] = <-nc
			}
			for i := 0; i < nothr; i++ {
				if err = errs[i]; err != nil {
					goto end
				}
			}
		}

		{
			var f ScheduleFlags
			if isForInstructor {
				f |= ForInstructor
				f |= OrphanRsv
			}
			pb := ScheduleBuffers{
				CanRegisterBuffers: rsv.CanRegisterBuffers{
					Vacations:   vix,
					RsvCount:    rcix,
					GroupedRsvs: grix,
				},
				GroupedSm: sub.NewGroupedSm(subs),
			}
			scheduleRes, err = ctx.NewSchedule(trainings, r.Rsv, subs, requestRange, f, flags, pb, true)
			if err != nil {
				status = 500
				goto end
			}
		}

		g.AbortWithStatusJSON(200, scheduleRes)
		return
	end:
		g.AbortWithError(status, err)
	}
}
