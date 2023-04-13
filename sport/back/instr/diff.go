package instr

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JsonMapBuff map[string]string

func (ud *JsonMapBuff) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), ud)
}

func (ctx *Ctx) DalGetDiffs(trainingID *uuid.UUID) (map[int32]map[string]string, error) {
	qs := "select id, name from diffs"
	if trainingID != nil {
		qs += " inner join trainings t on t.id = $1 and diffs.id = any(t.diff)"
	}
	q, err := ctx.Dal.Db.Query(qs)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	ret := make(map[int32]map[string]string)
	var tmp JsonMapBuff
	var id int32
	for q.Next() {
		if err := q.Scan(&id, &tmp); err != nil {
			return nil, err
		}
		ret[id] = tmp
	}
	return ret, nil
}

func (ctx *Ctx) HandlerGetDiffs() gin.HandlerFunc {
	return func(g *gin.Context) {
		r, err := ctx.DalGetDiffs(nil)
		if err != nil {
			g.AbortWithError(500, err)
		}
		g.AbortWithStatusJSON(200, r)
	}
}
