package train

import (
	"database/sql"
	"fmt"
	"sport/helpers"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (ctx *Ctx) DalCreateTrainingGroup(r *Group) error {
	q := `
		insert into training_groups (id, user_id, name, max_people, max_trainings)
		values ($1, $2, $3, $4, $5)
	`
	_, err := ctx.Dal.Db.Exec(q, r.ID, r.UserID, r.Name, r.MaxPeople, r.MaxTrainings)
	return err
}

func (ctx *Ctx) DalUpdateTrainingGroup(
	userID uuid.UUID, r *UpdateGroupRequest) error {
	tx, err := ctx.Dal.Db.Begin()
	if err != nil {
		return err
	}
	q := `
		update training_groups set
			name = $1,
			max_people = $2,
			max_trainings = $3
		where id = $4 and user_id = $5
	`
	res, err := tx.Exec(q, r.Name, r.MaxPeople, r.MaxTrainings, r.ID, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := helpers.PgMustBeOneRow(res); err != nil {
		tx.Rollback()
		return err
	}

	if err = ctx.DalUpdateTvg(r.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ctx *Ctx) DalDeleteTrainingGroup(
	userID uuid.UUID, r *DeleteGroupRequest,
) error {
	tx, err := ctx.Dal.Db.Begin()
	if err != nil {
		return err
	}

	q := "delete from training_groups where id = $1 and user_id = $2"
	res, err := tx.Exec(q, r.ID, userID)
	if err != nil {
		return err
	}

	if err := helpers.PgMustBeOneRow(res); err != nil {
		tx.Rollback()
		return err
	}

	if err = ctx.DalUpdateTvg(r.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

type ReadTrainingGroupsRequest struct {
	UserID *uuid.UUID
	IDs    []string
}

func (ctx *Ctx) DalReadTrainingGroups(req ReadTrainingGroupsRequest) ([]Group, error) {
	qb := strings.Builder{}
	qb.WriteString(`select 
			id, 
			user_id,
			name,
			max_people,
			max_trainings
		from training_groups `)
	wroteWhere := false
	args := make([]interface{}, 0, 1)
	if req.IDs != nil {
		args = append(args, pq.StringArray(req.IDs))
		qb.WriteString(fmt.Sprintf(" where id = any($%d) ", len(args)))
		//qb.WriteString(fmt.Sprintf(" id @> $%d ", len(args)))
	}

	if req.UserID != nil {
		if wroteWhere {
			qb.WriteString(" and ")
		} else {
			qb.WriteString(" where ")
			wroteWhere = true
		}
		args = append(args, *req.UserID)
		qb.WriteString(fmt.Sprintf(" user_id = $%d ", len(args)))
	}

	qr, err := ctx.Dal.Db.Query(qb.String(), args...)
	if err != nil {
		return nil, err
	}
	res := make([]Group, 0, 2)
	var tmp Group

	for qr.Next() {
		if err := qr.Scan(&tmp.ID, &tmp.UserID, &tmp.Name, &tmp.MaxPeople, &tmp.MaxTrainings); err != nil {
			qr.Close()
			return nil, err
		}
		res = append(res, tmp)
	}
	return res, qr.Close()
}

// Tvg stands for trainings_v_groups
// this function will update tvg after one of the group changes (gets updated / deleted)
func (ctx *Ctx) DalUpdateTvg(groupID uuid.UUID, tx *sql.Tx) error {
	q := `
		with cte as (
			select 
				tvg.training_id, 
				json_agg(g.*) as groups, 
				array_agg(g.id) as group_ids
			from trainings_v_groups tvg, training_groups g
			where 
				tvg.group_ids @> array[$1] and
				tvg.group_ids @> array[g.id::text]
			group by tvg.training_id
		)
		update trainings_v_groups tvg
		set 
			groups = c.groups,
			group_ids = c.group_ids
		from cte c where tvg.training_id = c.training_id 
	`
	var err error
	if tx != nil {
		_, err = tx.Exec(q, groupID)
	} else {
		_, err = ctx.Dal.Db.Exec(q, groupID)
	}

	/*

		from (select * from training_groups) as sq
		where tvg.group_ids @> array[$1] and array[sq.id::text] @> tvg.group_ids
	*/

	return err
}

// adds specified group to the training
func (ctx *Ctx) DalAddTvgBinding(
	trainingID, groupID, userID uuid.UUID,
) error {
	q := `with cte as (
		select unnest(group_ids) as group_id
		from trainings_v_groups tvg
		where tvg.training_id = $1
		union
		select $2 as group_id
	), cte2 as (
		select
			json_agg(distinct g.*) as groups,
			array_agg(distinct g.id) as group_ids
		from cte, training_groups g
		where
			cte.group_id = g.id::text and
			g.user_id = $3
	)
	insert into trainings_v_groups
		(training_id, groups, group_ids)
	select
		$1, coalesce(groups, '[]'), coalesce(group_ids, '{}')
	from cte2
	on conflict (training_id) do update
	set
		groups = excluded.groups,
		group_ids = excluded.group_ids`
	res, err := ctx.Dal.Db.Exec(q, trainingID, groupID, userID)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalRemoveTvgBinding(trainingID, groupID, userID uuid.UUID) error {
	q := `with cte as (
		select unnest(group_ids) as group_id
		from trainings_v_groups tvg
		where tvg.training_id = $1
	), cte2 as (
		select
			json_agg(distinct g.*) as groups,
			array_agg(distinct g.id) as group_ids
		from cte, training_groups g
		where
			cte.group_id != $2 and
			cte.group_id = g.id::text and
			g.user_id = $3
	)
	update trainings_v_groups tvg set
		groups = coalesce(c.groups, '[]'),
		group_ids = coalesce(c.group_ids, '{}')
	from cte2 c where tvg.training_id = $1`
	_, err := ctx.Dal.Db.Exec(q, trainingID, groupID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Ctx) DalReadTvgBindings(trainingID uuid.UUID) (GroupArr, error) {
	q := `select groups from trainings_v_groups where training_id = $1`
	var res GroupArr
	r := ctx.Dal.Db.QueryRow(q, trainingID)
	if err := r.Scan(&res); err != nil {
		return nil, err
	}
	return res, nil
}
