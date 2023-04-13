package train

import (
	"context"
	"database/sql"
	"fmt"
	"sport/dc"
	"sport/helpers"
	"sport/sub"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type DalReadTrainingsRequest struct {

	// only one of those 2 args may be specified at a time
	UserID       *uuid.UUID
	InstructorID *uuid.UUID

	// optional args
	TrainingID  *uuid.UUID
	TrainingIDs []string
	Query       string
	DateRange   helpers.DateRange
	//
	WithOccs      bool
	WithGroups    bool
	WithDcs       bool
	WithSubModels bool
	//
	GroupIds []string
	//
	DcID *uuid.UUID
	SmID *uuid.UUID
	//
	SmIDs []string
}

// be careful when using this function.
// if both userID and training_id are nil
// then this function will result in full scan on table trainings.
/*
	this function is growing and in my opinion getting out of control
	someone, sometime may want to refactor this
*/
func (ctx *Ctx) DalReadTrainings(
	r DalReadTrainingsRequest,
) ([]*TrainingWithJoins, error) {

	withSearch := false
	if r.Query != "" {
		withSearch = true
	}

	params := make([]interface{}, 0, 3)
	nextArgIndex := 1
	q := strings.Builder{}
	wroteWithCte := false

	if r.WithDcs {
		q.WriteString(`with dccte as (
			select json_agg(dc.*) as dcs, dct.training_id
			from dc
			inner join dc_v_train dct on dct.dc_id = dc.id `)
		if r.DcID != nil {
			q.WriteString(fmt.Sprintf(` where dc_id = $%d `, nextArgIndex))
			nextArgIndex++
			params = append(params, *r.DcID)
		}
		q.WriteString(` group by dct.training_id  ) `)
		wroteWithCte = true
	}

	if r.WithSubModels {
		if wroteWithCte {
			q.WriteString(" ,")
		} else {
			q.WriteString(" with")
		}
		q.WriteString(` smcte as (
			select json_agg(distinct sm.*) as sms, t.id as training_id
			from trainings t
			left outer join sub_model_bindings smb on smb.training_id = t.id
			left outer join sub_models sm on sm.id = smb.sub_model_id or 
				(t.instructor_id = sm.instructor_id and sm.all_trainings_by_def = true)
			where sm.id is not null 
			 `)
		if r.SmID != nil {
			q.WriteString(fmt.Sprintf(` and sm.id = $%d `, nextArgIndex))
			nextArgIndex++
			params = append(params, *r.SmID)
		}
		if r.SmIDs != nil {
			q.WriteString(fmt.Sprintf(" and sm.id @> $%d", nextArgIndex))
			nextArgIndex++
			params = append(params, pq.StringArray(r.SmIDs))
		}
		q.WriteString(` group by t.id  ) `)
	}

	q.WriteString("select")
	q.WriteString(TrainingSelectFields())
	if r.WithOccs {
		q.WriteString(",")
		q.WriteString(NullableOccurrenceSelectFields())
		q.WriteString(",")
		q.WriteString(NullableSecondaryOccSelectFields())
	}
	if r.WithGroups {
		q.WriteString(",coalesce(g.groups, '[]')")
	}
	if r.WithDcs {
		q.WriteString(",coalesce(dccte.dcs, '[]')")
	}
	if r.WithSubModels {
		q.WriteString(",coalesce(smcte.sms, '[]')")
	}
	//q.WriteString("")

	occStr := ""
	if r.WithOccs {
		occStr = ` left outer join occurrences o on o.training_id = t.id 
					left outer join secondary_occs so on so.occ_id = o.id`
	}

	grp := ""
	if r.WithGroups {
		grp = ` left outer join trainings_v_groups g on g.training_id = t.id `
	}

	dcstr := ""
	if r.WithDcs {
		if r.DcID == nil {
			dcstr = ` left outer join dccte on dccte.training_id = t.id`
		} else {
			dcstr = ` inner join dccte on dccte.training_id = t.id`
		}
	}

	smstr := ""
	if r.WithSubModels {
		if r.SmID == nil {
			smstr = ` left outer join smcte on smcte.training_id = t.id`
		} else {
			smstr = ` inner join smcte on smcte.training_id = t.id`
		}
	}

	if r.InstructorID != nil {
		q.WriteString(" from trainings t ")
		q.WriteString(occStr)
		q.WriteString(grp)
		q.WriteString(dcstr)
		q.WriteString(smstr)
		q.WriteString(fmt.Sprintf(" where t.instructor_id = $%d ", nextArgIndex))
		nextArgIndex++
		params = append(params, r.InstructorID)
	} else if r.UserID != nil {
		q.WriteString(` from users u 
			inner join instructors i on i.user_id = u.id
			inner join trainings t on i.id = t.instructor_id `)
		q.WriteString(occStr)
		q.WriteString(grp)
		q.WriteString(dcstr)
		q.WriteString(smstr)
		q.WriteString(fmt.Sprintf(" where u.id = $%d ", nextArgIndex))
		nextArgIndex++
		params = append(params, r.UserID)
	} else {
		q.WriteString(" from trainings t ")
		q.WriteString(occStr)
		q.WriteString(grp)
		q.WriteString(dcstr)
		q.WriteString(smstr)
		q.WriteString(" where 1=1 ")
	}

	if r.TrainingID != nil {
		q.WriteString(fmt.Sprintf(" and t.id = $%d", nextArgIndex))
		nextArgIndex++
		params = append(params, r.TrainingID)
	}

	if r.TrainingIDs != nil {
		q.WriteString(" and ")
		helpers.PgUuidMatch(&q, r.TrainingIDs, nextArgIndex, &params, "t.id")
		// q.WriteString("and ( ")
		// for i := range r.TrainingIDs {
		// 	if i > 0 {
		// 		q.WriteString("or")
		// 	}
		// 	q.WriteString(fmt.Sprintf(" t.id = $%d::uuid ", nextArgIndex))
		// 	nextArgIndex++
		// 	params = append(params, r.TrainingIDs[i])
		// }
		// q.WriteString(") ")
	}

	if withSearch {
		q.WriteString(fmt.Sprintf(" and t.title like '%%' || $%d || '%%' limit 5 ", nextArgIndex))
		nextArgIndex++
		params = append(params, r.Query)
	}

	if r.GroupIds != nil {
		if !r.WithGroups {
			return nil, fmt.Errorf("cant have groupIds if 'WithGroups' is false")
		}
		q.WriteString(fmt.Sprintf(" and g.group_ids @> $%d ", nextArgIndex))
		nextArgIndex++
		params = append(params, pq.StringArray(r.GroupIds))
	}

	/*
		you may implement this range search to <maybe/probably> improve perf
	*/
	// if r.DateRange.IsNotZero {
	// 	q += fmt.Sprintf(` and  (it.date_start = '0001-01-01' or it.date_start <= $%d)`, nextArgIndex)
	// 	nextArgIndex++
	// 	q += fmt.Sprintf(` and (it.date_end = '0001-01-01' or it.date_end >= $%d)`, nextArgIndex)
	// 	params = append(params, r.DateRange.End.In(time.UTC), r.DateRange.Start.In(time.UTC))
	// 	nextArgIndex++
	// }

	res, err := ctx.Dal.Db.Query(q.String(), params...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	c := 1
	if r.TrainingID == nil {
		c = 10
	}
	ret := make([]*TrainingWithJoins, 0, c)
	var tmp *TrainingWithJoins
	var occ Occ
	var socc SecondaryOcc
	args := []interface{}{}
	// since occ is allocaed once those addresses wont change
	// we can preallocate array
	var occArgs []interface{}
	if r.WithOccs {
		occArgs = occ.ScanFields()
		occArgs = append(occArgs, socc.ScanFields()...)
	}
	var grpArr GroupArr
	var dcArr dc.DcArr
	var smArr sub.SmArr
	type scanIndexItem struct {
		Training int
		OccIndex map[uuid.UUID]int
	}
	// helper map used to match trainings and occs
	scanIx := make(map[uuid.UUID]*scanIndexItem)
	for res.Next() {
		tmp = new(TrainingWithJoins)
		args = tmp.Training.ScanFields()
		if r.WithOccs {
			args = append(args, occArgs...)
		}
		if r.WithGroups {
			args = append(args, &grpArr)
		}
		if r.WithDcs {
			args = append(args, &dcArr)
		}
		if r.WithSubModels {
			args = append(args, &smArr)
		}
		if err = res.Scan(args...); err != nil {
			return nil, err
		}
		if r.WithOccs {
			if ix, ok := scanIx[tmp.Training.ID]; ok {
				if oix, ok := ix.OccIndex[occ.ID]; ok {
					ret[ix.Training].Occurrences[oix].SecondaryOccs =
						append(ret[ix.Training].Occurrences[oix].SecondaryOccs, socc)
				} else {
					if socc.ID == uuid.Nil {
						ret[ix.Training].Occurrences =
							append(ret[ix.Training].Occurrences,
								OccWithSecondary{Occ: occ})
					} else {
						ret[ix.Training].Occurrences =
							append(ret[ix.Training].Occurrences,
								OccWithSecondary{Occ: occ, SecondaryOccs: []SecondaryOcc{socc}})
					}
					ix.OccIndex[occ.ID] = len(ret[ix.Training].Occurrences) - 1
				}
			} else {
				var oix map[uuid.UUID]int
				if occ.ID != uuid.Nil {
					var soccs []SecondaryOcc
					if socc.ID != uuid.Nil {
						soccs = make([]SecondaryOcc, 0, 2)
						soccs = append(soccs, socc)
					}
					tmp.Occurrences = make([]OccWithSecondary, 0, 2)
					tmp.Occurrences = append(tmp.Occurrences, OccWithSecondary{
						Occ:           occ,
						SecondaryOccs: soccs,
					})
					oix = map[uuid.UUID]int{
						occ.ID: 0,
					}
				} else {
					tmp.Occurrences = []OccWithSecondary{}
				}
				tmp.Training.PostprocessAfterDbScan(ctx)
				tmp.Groups = grpArr
				tmp.Dcs = dcArr
				tmp.Sms = smArr
				ret = append(ret, tmp)
				scanIx[tmp.Training.ID] = &scanIndexItem{
					Training: len(ret) - 1,
					OccIndex: oix,
				}
			}
		} else {
			tmp.Training.PostprocessAfterDbScan(ctx)
			tmp.Groups = grpArr
			tmp.Dcs = dcArr
			tmp.Sms = smArr
			ret = append(ret, tmp)
		}
	}

	return ret, err
}

// same as DalReadInstructorTrainings
// but performs validation on quantity of returned elements
func (ctx *Ctx) DalReadSingleTraining(r DalReadTrainingsRequest) (*TrainingWithJoins, error) {
	t, err := ctx.DalReadTrainings(r)
	if err != nil {
		return nil, err
	}
	if len(t) != 1 {
		return nil, sql.ErrNoRows
	}
	return t[0], nil
}

func (ctx *Ctx) DalDeleteTraining(id, userID uuid.UUID) error {

	instructorID, err := ctx.Instr.DalReadInstructorID(userID)
	if err != nil {
		return err
	}

	tx, err := ctx.Dal.Db.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(
		"delete from trainings where id = $1 and instructor_id = $2",
		id,
		instructorID)

	if err != nil {
		tx.Rollback()
		return err
	}
	if err := helpers.PgMustBeOneRow(res); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ctx *Ctx) DalUpdateTraining(
	tt *UpdateTrainingRequest,
	instructorID uuid.UUID,
	userID uuid.UUID,
) error {

	tx, err := ctx.Dal.Db.BeginTx(
		context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	res, err := ctx.Dal.Db.Exec(`update trainings set
			title = $1,
			description = $2,
			date_start = $3,
			date_end = $4,
			capacity = $5,
			location_text = $6,
			location_lat = $7,
			location_lng = $8,
			location_country = $9,
			price = $10,
			currency = $11,
			tags = $12,
			diff = $13,
			manual_confirm = $14,
			allow_express = $15,
			disabled = $16,

			required_gear = $17,
			recommended_gear = $18,
			instructor_gear = $19,
			min_age = $20,
			max_age = $21,
			training_supports_disabled = $22,
			place_supports_disabled = $23
		where id = $24 and instructor_id = $25`,
		tt.Title,
		tt.Description,
		tt.DateStart,
		tt.DateEnd,
		tt.Capacity,
		tt.LocationText,
		tt.LocationLat,
		tt.LocationLng,
		tt.LocationCountry,
		tt.Price,
		tt.Currency,
		pq.Array(tt.Tags),
		pq.Array(tt.Diff),
		tt.ManualConfirm,
		tt.AllowExpress,
		tt.Disabled,
		pq.StringArray(tt.RequiredGear),
		pq.StringArray(tt.RecommendedGear),
		pq.StringArray(tt.InstructorGear),
		tt.MinAge,
		tt.MaxAge,
		tt.TrainingSupportsDisabled,
		tt.PlaceSupportsDisabled,
		//
		tt.ID,
		instructorID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := helpers.PgMustBeOneRow(res); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (ctx *Ctx) DalCreateTraining(
	t *Training, userID uuid.UUID, tx *sql.Tx,
) error {

	_, err := ctx.Dal.Db.Exec(`insert into trainings (
			id,
			instructor_id,
			title,
			description,
			date_start,
			date_end,
			capacity,
			
			required_gear,
			recommended_gear,
			instructor_gear,
			min_age,
			max_age,

			location_text,
			location_lat,
			location_lng,
			location_country,
			price,
			currency,
			tags,
			diff,
			number_reviews,
			avg_mark,
			manual_confirm,
			allow_express,
			created_on,
			main_img_relpath,
			secondary_img_relpaths,
			disabled,
			training_supports_disabled,
			place_supports_disabled
		) values (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14,
			$15,
			$16,
			$17,
			$18,
			$19,
			$20,
			$21,
			$22,
			$23,
			$24,
			$25,
			$26,
			$27,
			$28,
			$29,
			$30
		)`,
		t.ID,
		t.InstructorID,
		t.Title,
		t.Description,
		t.DateStart,
		t.DateEnd,
		t.Capacity,
		pq.StringArray(t.RequiredGear),
		pq.StringArray(t.RecommendedGear),
		pq.StringArray(t.InstructorGear),
		t.MinAge,
		t.MaxAge,
		t.LocationText,
		t.LocationLat,
		t.LocationLng,
		t.LocationCountry,
		t.Price,
		t.Currency,
		pq.Array(t.Tags),
		pq.Array(t.Diff),
		t.NumberReviews,
		t.AvgMark,
		t.ManualConfirm,
		t.AllowExpress,
		t.CreatedOn,
		t.MainImgID,
		pq.StringArray(t.SecondaryImgIDs),
		t.Disabled,
		t.TrainingSupportsDisabled,
		t.PlaceSupportsDisabled,
	)

	if err != nil {
		return err
	}

	return nil
}

func (ctx *Ctx) DalUpdateTrainingMainImage(
	trainingID uuid.UUID, path string) error {

	q := `update trainings set main_img_relpath = $1 where id = $2`
	r, err := ctx.Dal.Db.Exec(q, path, trainingID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(r)
}

func (ctx *Ctx) DalAddTrainingSecondaryImage(
	trainingID uuid.UUID, path string) error {

	q := `update trainings set 
		secondary_img_relpaths = array_append(secondary_img_relpaths, $1) 
	where id = $2`
	r, err := ctx.Dal.Db.Exec(q, path, trainingID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(r)
}

func (ctx *Ctx) DalRemoveTrainingSecondaryImage(
	trainingID uuid.UUID, path string,
) error {

	q := `update trainings set 
		secondary_img_relpaths = array_remove(secondary_img_relpaths, $1) 
	where id = $2`

	r, err := ctx.Dal.Db.Exec(q, path, trainingID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(r)
}

func (ctx *Ctx) DalUpdateTrainingAvgMark(id uuid.UUID, mark int, tx *sql.Tx) error {

	if mark == 0 {
		return nil
	}

	q := ""

	if mark > 0 {
		q = `update trainings set 
				number_reviews = number_reviews + 1,
				avg_mark = ((avg_mark * number_reviews) + $1) / (number_reviews + 1)
			where id = $2`
	} else {
		q = `update trainings set 
				number_reviews = number_reviews - 1,
				avg_mark = case when number_reviews > 1 then
					((number_reviews * avg_mark) - $1) / (number_reviews - 1)
				else
					-$1
				end
			where id = $2 and avg_mark >= $1 and number_reviews >= 1`
	}

	var res sql.Result
	var err error

	args := []interface{}{
		mark, id,
	}

	if tx != nil {
		res, err = tx.Exec(q, args...)
	} else {
		res, err = ctx.Dal.Db.Exec(q, args...)
	}

	if err != nil {
		return err
	}

	return helpers.PgMustBeOneRow(res)
}
