package main

import (
	"html/template"
	"os"
	"strings"
)

type TriggerDst struct {
	TableName          string
	OldExtra, NewExtra string
}

func main() {
	dst := "../search/scripts/"
	tables := []TriggerDst{
		{
			TableName: "trainings",
			// return training id
			OldExtra: "old.id",
			NewExtra: "new.id",
		}, {
			TableName: "instructors",
			// return instructor id
			OldExtra: "old.id",
			NewExtra: "new.id",
		}, {
			TableName: "occurrences",
			// return training id
			OldExtra: "old.training_id",
			NewExtra: "new.training_id",
		}, {
			// return instructor id
			TableName: "instr_vacations",
			OldExtra:  "old.instructor_id",
			NewExtra:  "new.instructor_id",
		}, {
			// return rsv id
			TableName: "reservations",
			OldExtra:  "old.id",
			NewExtra:  "new.id",
		}, {
			// return occ id
			TableName: "secondary_occs",
			OldExtra:  "old.training_id",
			NewExtra:  "new.training_id",
		}, {
			// return sub model id
			TableName: "sub_models",
			OldExtra:  "old.id",
			NewExtra:  "new.id",
		}, {
			// return training id
			TableName: "sub_model_bindings",
			OldExtra:  "old.training_id",
			NewExtra:  "new.training_id",
		}, {
			// return sub id
			TableName: "subs",
			OldExtra:  "old.id",
			NewExtra:  "new.id",
		},
		{
			// return group id
			TableName: "training_groups",
			OldExtra:  "old.id",
			NewExtra:  "new.id",
		}, {
			// return training id
			TableName: "trainings_v_groups",
			OldExtra:  "old.training_id",
			NewExtra:  "new.training_id",
		}, {
			// return user id
			TableName: "users",
			OldExtra:  "old.id",
			NewExtra:  "new.id",
		}, {
			TableName: "reviews",
			OldExtra:  "old.training_id",
			NewExtra:  "new.training_id",
		},
	}

	t, err := template.ParseFiles("../trigger_gen/template.sql")
	if err != nil {
		panic(err)
	}

	for i := range tables {
		trigger := strings.Builder{}
		fp, err := os.OpenFile(dst+tables[i].TableName+"_trigger.sql", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		err = t.ExecuteTemplate(&trigger, "template.sql", tables[i])
		if err != nil {
			panic(err)
		}
		if _, err := fp.WriteString(trigger.String()); err != nil {
			panic(err)
		}
	}

}
