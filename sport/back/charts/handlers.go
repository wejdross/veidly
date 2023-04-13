package charts

import (
	"database/sql"
	"net/http"
	"sport/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctx *Ctx) PostCharts() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		var ChData ChartData

		if err := helpers.ReadJsonBodyFromReader(g.Request.Body, &ChData, ChData.Validate); err != nil {
			g.AbortWithError(400, err)
			return
		}

		iid, err := ctx.Instr.DalReadInstructorID(userID)
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		if err := DalPOSTChart(&ChData, userID, iid, ctx); err != nil {
			g.AbortWithError(http.StatusBadRequest, err)
			return
		}
		/*
					insert
			    into charts (label, instructor_id, user_id, timestamp, value, chart_type)
			    values ('waga', 'aeb1039f-a6e4-45c7-887b-b053c911750c',  '8b46c15b-5e95-4568-9129-80095fe81a9e', 1629754312, 90, 1);

				# analogicznie dla instruktora
				select label, timestamp,value,chart_type from charts where user_id = '6a9d61f4-f15e-404e-8739-bb428a3e5984' order by label;
		*/
		g.JSON(200, gin.H{
			"charts": "added!",
		})

	}
}

func (ctx *Ctx) GetChartsUserPerspective() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)

		if _, err := DalGetChartsUserPerspective(userID, ctx); err != nil {
			g.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}
}

/*
curl -d '{  "label": "Przyleciałem z Api, tak trzeba będzie mnie ogarnąć na bazie",  "type": 1,  "values": [    {      "timestamp": 1629318082,      "value": 100    },    {      "timestamp": 1629318182,      "value": 101    },    {      "timestamp": 1629318282,      "value": 102    },    {      "timestamp": 1629318382,      "value": 103    },    {      "timestamp": 1629318482,      "value": 104    },    {      "timestamp": 1629318582,      "value": 105    }  ]}' -H 'Accept: application/json' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzE1Njc4NzgsInVpZCI6ImUxOGFjZjBhLTcyYjgtNDEyOS1iOTdkLTVmY2E1ODVkMTUyMyJ9.SFdjeWUVBz_rtGSWr7ikDIgF_RT19i2pk59fhC98YSg" 127.0.0.1:1580/api/charts | python3 -m json.tool
curl -H 'Accept: application/json' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzE1Njc4NzgsInVpZCI6ImUxOGFjZjBhLTcyYjgtNDEyOS1iOTdkLTVmY2E1ODVkMTUyMyJ9.SFdjeWUVBz_rtGSWr7ikDIgF_RT19i2pk59fhC98YSg" 127.0.0.1:1580/api/charts/user | python3 -m json.tool
*/
