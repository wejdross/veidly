package charts

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type ChartAtomicData struct {
	Timestamp int `json:"timestamp"`
	Value     int `json:"value"`
}

type ChartData struct {
	Label  string            `json:"label"`
	Type   int               `json:"type"`
	Values []ChartAtomicData `json:"values"`
}

type ExtendedChartData struct {
	Id             int       `json:"id"`
	ChD            ChartData `json:"chart_data"`
	UserData       []byte    `json:"user_data"`
	UserAvatarPath string    `json:"user_avatar"`
}

func (ChD *ChartData) Validate() error {
	if len(ChD.Values) < 1 {
		return fmt.Errorf("No values to add")
	}

	// just two types of charts would be supported at the beginning
	switch ChD.Type {
	case 0:
	case 1:
	default:
		return fmt.Errorf("Chart type not supported")
	}

	if len(ChD.Label) > 255 || len(ChD.Label) == 0 {
		return fmt.Errorf("Label too long")
	}

	return nil
}

func DalGetChartsUserPerspective(uuid uuid.UUID, ctx *Ctx) (retVal []ExtendedChartData, err error) {
	var m = []ExtendedChartData{}
	// user perspective = user query for his data
	// select charts.*,users.user_data, users.avatar_relpath from charts left join users on charts.user_id = users.id where charts.user_id = 'e18acf0a-72b8-4129-b97d-5fca585d1523';
	var rows *sql.Rows
	const q = `select 
					charts.*,
					users.user_data, 
					users.avatar_relpath 
					from charts 
					left join users on charts.user_id = users.id 
					where charts.user_id = $1
					order by label`

	rows, err = ctx.Dal.Db.Query(q, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// add here tmp values for each field and then add to structs, to be done on wednesday!!
		var tmp ExtendedChartData
		if err = rows.Scan(&tmp); err != nil {
			return nil, err
		}
		fmt.Println(tmp)
	}
	return m, nil
}

func DalGetPossibleCharts(DataObj *ChartData, uuid, iuid uuid.UUID, ctx *Ctx) (err error) {
	return nil
}

func DalPOSTChart(DataObj *ChartData, uuid, iuid uuid.UUID, ctx *Ctx) (err error) {
	const q = `
	insert into charts (
		label, 
		instr_id, 
		user_id, 
		timestamp, 
		value, 
		chart_type)
			    
		values (
			$1, 
			$2,  
			$3, 
			$4, 
			$5, 
			$6)
	`

	for _, val := range DataObj.Values {
		args := []interface{}{
			DataObj.Label,
			iuid,
			uuid,
			val.Timestamp,
			val.Value,
			DataObj.Type,
		}

		_, err = ctx.Dal.Db.Exec(q, args...)
		if err != nil {
			return err
		}
	}

	return nil
}
