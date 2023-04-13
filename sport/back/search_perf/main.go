package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sport/adyen"
	"sport/api"
	"sport/config"
	"sport/dal"
	"sport/dc"
	"sport/helpers"
	"sport/instr"
	"sport/rsv"
	"sport/schedule"
	"sport/search"
	"sport/static"
	"sport/sub"
	"sport/train"
	"sport/user"
	"strings"
	"time"
)

type SearchRequest struct {
	Query    string
	Days     int
	Lat, Lng float64
	DistKm   int
}

func DoSearch(apiCtx *api.Ctx, req SearchRequest) error {

	start := helpers.NowMin()
	end := start.Add(time.Duration(req.Days) * time.Hour * 24)

	path := "/api/search"

	err := apiCtx.TestAssertCaseErr(&api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    path,
		RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
			SearchRequestOpts: search.SearchRequestOpts{
				Langs:              []string{"pl"},
				Lat:                req.Lat,
				Lng:                req.Lng,
				DistKm:             req.DistKm,
				DateStart:          start,
				DateEnd:            end,
				OmitEmptySchedules: true,
			},
			Query: req.Query,
		}),
		ExpectedStatusCode: 200,
		ExpectedBodyVal: func(b []byte, i interface{}) error {
			var sr search.SearchResultsWithMetadata
			helpers.JsonMustDeserialize(b, &sr)
			tc := len(sr.Data)
			totalSched := 0
			top := 15
			fmt.Printf("Top %d results:\n", top)
			for i := range sr.Data {
				totalSched += len(sr.Data[i].Schedule)
				if i < top {
					fmt.Printf("RESULT: \n\tinstructor=%s \n\ttitle=%s \n\ttags=%s \n\tscore: total=%f tag=%f instr=%f title=%f\n",
						sr.Data[i].UserInfo.Name,
						sr.Data[i].Training.Title,
						strings.Join(sr.Data[i].Training.Tags, ","),
						sr.Data[i].Score.TotalScore(),
						sr.Data[i].Score.TagScore,
						sr.Data[i].Score.InstrScore,
						sr.Data[i].Score.TitleScore,
					)
				}
			}
			fmt.Printf("Total %d trainings, with %d available dates\n", tc, totalSched)
			fmt.Println("---- PERF ----")
			fmt.Println(helpers.JsonMustSerializeFormatStr(sr.Meta.Perf))
			return nil
		},
	})

	return err
}

func DoOne(apiCtx *api.Ctx) error {

	var req SearchRequest

	r := bufio.NewReader(os.Stdin)

	fmt.Print("query: ")
	req.Query, _ = r.ReadString('\n')

	fmt.Print("days: ")
	if _, err := fmt.Fscanf(os.Stdin, "%d", &req.Days); err != nil {
		return err
	}

	req.Lat = 50.2649
	req.Lng = 19.0238

	// fmt.Print("lat: ")
	// if _, err := fmt.Fscanf(os.Stdin, "%f", &req.Lat); err != nil {
	// 	return err
	// }
	// fmt.Print("lng: ")
	// if _, err := fmt.Fscanf(os.Stdin, "%f", &req.Lng); err != nil {
	// 	return err
	// }

	fmt.Print("dist [km]: ")
	if _, err := fmt.Fscanf(os.Stdin, "%d", &req.DistKm); err != nil {
		return err
	}

	return DoSearch(apiCtx, req)
}

func main() {

	testdb := "sportdb_sg"
	apiCtx := api.NewApi(config.NewLocalCtx())
	dalCtx := dal.NewDal(apiCtx.Config, testdb)
	staticCtx := static.NewCtx(apiCtx)
	userCtx := user.NewCtx(apiCtx, dalCtx, staticCtx, nil)
	adyenCtx := adyen.NewMockupCtx(apiCtx)
	instrCtx := instr.NewCtx(apiCtx, dalCtx, userCtx, adyenCtx)
	dcCtx := dc.NewCtx(apiCtx, dalCtx, instrCtx)
	trainCtx := train.NewCtx(apiCtx, dalCtx, userCtx, staticCtx, instrCtx)
	subCtx := sub.NewCtx(apiCtx, dalCtx, userCtx, instrCtx,
		adyenCtx, nil, nil)
	rsvCtx := rsv.NewCtx(
		apiCtx, dalCtx, userCtx,
		instrCtx, trainCtx, adyenCtx, nil,
		nil, dcCtx, subCtx, nil)
	schedCtx := schedule.NewCtx(apiCtx, instrCtx, trainCtx, subCtx, rsvCtx)
	searchCtx := search.NewCtx(
		apiCtx, dalCtx, userCtx, instrCtx, trainCtx, rsvCtx, schedCtx)

	err := searchCtx.RegenerateCache()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("-- CACHE BUILD PERF --\n%v\n",
		helpers.JsonMustSerializeFormatStr(searchCtx.Cache.BuildPerf))

	if err := Diff(searchCtx); err != nil {
		log.Fatal(err)
	}

	for {
		err := DoOne(apiCtx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

}
