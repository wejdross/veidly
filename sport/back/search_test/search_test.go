package search_test

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sport/api"
	"sport/helpers"
	"sport/search"
	"sport/train"
	"strings"
	"testing"
	"time"
)

var searchPath = "/api/search"

func createTestTrainings() error {
	token, err := searchCtx.User.ApiCreateAndLoginUser(nil)
	if err != nil {
		return err
	}

	userID, err := searchCtx.Api.AuthorizeUserFromToken(token)
	if err != nil {
		return err
	}

	_, err = searchCtx.Instr.CreateTestInstructor(userID, nil)
	if err != nil {
		return err
	}

	tc := make([]api.TestCase, 0, 10)

	tc = append(tc, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/training",
		RequestReader: helpers.JsonMustSerializeReader(train.CreateTrainingRequest{
			Training: train.TrainingRequest{
				Title:           "test1",
				Capacity:        1,
				Price:           1500,
				Currency:        "PLN",
				LocationCountry: "PL",
				Tags:            nil,
				LocationLat:     0,
				LocationLng:     0,
				LocationText:    "123",
				MinAge:          10,
				MaxAge:          14,
			},
			Occurrences: []train.CreateOccRequest{
				{
					OccRequest: train.OccRequest{
						DateStart:  time.Date(2022, 02, 01, 8, 0, 0, 0, time.UTC),
						DateEnd:    time.Date(2022, 02, 01, 10, 0, 0, 0, time.UTC),
						RepeatDays: 7, // weekly
						Remarks:    helpers.CRNG_stringPanic(128),
						Color:      "ffffff",
					},
				},
			},
		}),
		ExpectedStatusCode: 204,
		RequestHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	tc = append(tc, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/training",
		RequestReader: helpers.JsonMustSerializeReader(train.CreateTrainingRequest{
			Training: train.TrainingRequest{
				Title:           "test4",
				Capacity:        5,
				Price:           9000,
				Currency:        "PLN",
				LocationCountry: "PL",
				Tags:            nil,
				LocationLat:     80,
				LocationLng:     100,
				LocationText:    "ggg",
				MinAge:          12,
				MaxAge:          16,
			},
			Occurrences: []train.CreateOccRequest{
				{
					OccRequest: train.OccRequest{
						DateStart:  time.Date(2022, 02, 01, 8, 0, 0, 0, time.UTC),
						DateEnd:    time.Date(2022, 02, 01, 10, 0, 0, 0, time.UTC),
						RepeatDays: 7, // weekly
						Remarks:    helpers.CRNG_stringPanic(128),
						Color:      "ffffff",
					},
				},
			},
		}),
		ExpectedStatusCode: 204,
		RequestHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	tc = append(tc, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/training",
		RequestReader: helpers.JsonMustSerializeReader(train.CreateTrainingRequest{
			Training: train.TrainingRequest{
				Title:           "test2",
				Capacity:        10,
				Price:           10000,
				Currency:        "PLN",
				LocationCountry: "PL",
				Tags:            nil,
				LocationLat:     30,
				LocationLng:     30,
				LocationText:    "foo",
			},
			Occurrences: []train.CreateOccRequest{
				{
					OccRequest: train.OccRequest{
						DateStart:  time.Date(2022, 02, 01, 8, 0, 0, 0, time.UTC),
						DateEnd:    time.Date(2022, 02, 01, 10, 0, 0, 0, time.UTC),
						RepeatDays: 7, // weekly
						Remarks:    helpers.CRNG_stringPanic(128),
						Color:      "ffffff",
					},
				},
			},
		}),
		ExpectedStatusCode: 204,
		RequestHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	tc = append(tc, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/training",
		RequestReader: helpers.JsonMustSerializeReader(train.CreateTrainingRequest{
			Training: train.TrainingRequest{
				Title:           "test3",
				Capacity:        20,
				Price:           2000,
				Currency:        "PLN",
				LocationCountry: "GB",
				Tags: []string{
					"tag1",
					"tag2",
				},
				LocationLat:  30,
				LocationLng:  30,
				LocationText: "foo",
			},
			Occurrences: []train.CreateOccRequest{
				{
					OccRequest: train.OccRequest{
						DateStart:  time.Date(2022, 02, 01, 8, 0, 0, 0, time.UTC),
						DateEnd:    time.Date(2022, 02, 01, 10, 0, 0, 0, time.UTC),
						RepeatDays: 7, // weekly
						Remarks:    helpers.CRNG_stringPanic(128),
						Color:      "ffffff",
					},
				},
			},
		}),
		ExpectedStatusCode: 204,
		RequestHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	tc = append(tc, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/training",
		RequestReader: helpers.JsonMustSerializeReader(train.CreateTrainingRequest{
			Training: train.TrainingRequest{
				Title:           "test5",
				Capacity:        20,
				Price:           20000,
				Currency:        "PLN",
				LocationCountry: "DE",
				Tags: []string{
					"tag1",
					"tag2",
				},
				Diff: []int32{
					1, 2, 3,
				},
				LocationLat:  70,
				LocationLng:  70,
				LocationText: "foo",
			},
			Occurrences: []train.CreateOccRequest{
				{
					OccRequest: train.OccRequest{
						// thursday:4
						DateStart:  time.Date(2021, 02, 04, 8, 0, 0, 0, time.UTC),
						DateEnd:    time.Date(2021, 02, 04, 18, 0, 0, 0, time.UTC),
						RepeatDays: 7, // weekly
					},
				},
			},
		}),
		ExpectedStatusCode: 204,
		RequestHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	tc = append(tc, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/training",
		RequestReader: helpers.JsonMustSerializeReader(train.CreateTrainingRequest{
			Training: train.TrainingRequest{
				Title:           "test6",
				Capacity:        20,
				Price:           20000,
				Currency:        "PLN",
				LocationCountry: "DE",
				Tags: []string{
					"tag1",
					"tag2",
				},
				Diff: []int32{
					3,
				},
				LocationLat:  70,
				LocationLng:  70,
				LocationText: "foo",
			},
			Occurrences: []train.CreateOccRequest{
				{
					OccRequest: train.OccRequest{
						// monday:1
						DateStart:  time.Date(2021, 02, 01, 14, 0, 0, 0, time.UTC),
						DateEnd:    time.Date(2021, 02, 01, 15, 0, 0, 0, time.UTC),
						RepeatDays: 7, // weekly
					},
				},
				{
					OccRequest: train.OccRequest{
						// friday:5
						DateStart:  time.Date(2021, 02, 05, 14, 0, 0, 0, time.UTC),
						DateEnd:    time.Date(2021, 02, 05, 16, 0, 0, 0, time.UTC),
						RepeatDays: 7, // weekly
					},
				},
			},
		}),
		ExpectedStatusCode: 204,
		RequestHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	tc = append(tc, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/training",
		RequestReader: helpers.JsonMustSerializeReader(train.CreateTrainingRequest{
			Training: train.TrainingRequest{
				Title:           "test7",
				Capacity:        20,
				Price:           2000,
				Currency:        "PLN",
				LocationCountry: "DE",
				Tags: []string{
					"tag1",
					"tag2",
				},
				LocationLat:  70,
				LocationLng:  70,
				LocationText: "foo",
			},
			Occurrences: []train.CreateOccRequest{
				{
					OccRequest: train.OccRequest{
						// tuesday:2
						DateStart:  time.Date(2021, 02, 02, 15, 0, 0, 0, time.UTC),
						DateEnd:    time.Date(2021, 02, 02, 19, 0, 0, 0, time.UTC),
						RepeatDays: 7, // weekly
					},
				},
			},
		}),
		ExpectedStatusCode: 204,
		RequestHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	tc = append(tc, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/training",
		RequestReader: helpers.JsonMustSerializeReader(train.CreateTrainingRequest{
			Training: train.TrainingRequest{
				Title:           "test8",
				Capacity:        20,
				Price:           200,
				Currency:        "PLN",
				LocationCountry: "DE",
				LocationLat:     70,
				LocationLng:     70,
				LocationText:    "foo",
			},
			Occurrences: []train.CreateOccRequest{
				{
					OccRequest: train.OccRequest{
						// saturday:6
						DateStart:  time.Date(2021, 02, 20, 20, 30, 0, 0, time.UTC),
						DateEnd:    time.Date(2021, 02, 20, 21, 0, 0, 0, time.UTC),
						RepeatDays: 3,
					},
				},
			},
		}),
		ExpectedStatusCode: 204,
		RequestHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	tc = append(tc, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/training",
		RequestReader: helpers.JsonMustSerializeReader(train.CreateTrainingRequest{
			Training: train.TrainingRequest{
				Title:           "test9",
				Capacity:        20,
				Price:           20000,
				Currency:        "PLN",
				LocationCountry: "DE",
				LocationLat:     70,
				LocationLng:     70,
				LocationText:    "foo",
			},
			Occurrences: []train.CreateOccRequest{
				{
					OccRequest: train.OccRequest{
						// saturday:6
						DateStart:  time.Date(2021, 04, 20, 20, 30, 0, 0, time.UTC),
						DateEnd:    time.Date(2021, 04, 20, 21, 0, 0, 0, time.UTC),
						RepeatDays: 0,
					},
				},
			},
		}),
		ExpectedStatusCode: 204,
		RequestHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	tc = append(tc, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/training",
		RequestReader: helpers.JsonMustSerializeReader(train.CreateTrainingRequest{
			Training: train.TrainingRequest{
				Title:           "test10",
				Capacity:        20,
				Price:           2000,
				Currency:        "PLN",
				LocationCountry: "DE",
				LocationLat:     70,
				LocationLng:     70,
				LocationText:    "foo",
			},

			Occurrences: []train.CreateOccRequest{
				{
					OccRequest: train.OccRequest{
						// saturday:6
						DateStart:  time.Date(2021, 05, 20, 20, 30, 0, 0, time.UTC),
						DateEnd:    time.Date(2021, 05, 20, 21, 0, 0, 0, time.UTC),
						RepeatDays: 10,
					},
				},
			},
		}),
		ExpectedStatusCode: 204,
		RequestHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	tc = append(tc, api.TestCase{
		RequestMethod: "POST",
		RequestUrl:    "/api/training",
		RequestReader: helpers.JsonMustSerializeReader(train.CreateTrainingRequest{
			Training: train.TrainingRequest{
				Title:           "test11",
				Capacity:        20,
				Price:           2000,
				Currency:        "PLN",
				LocationCountry: "DE",
				LocationLat:     70,
				LocationLng:     70,
				LocationText:    "foo",
				DateStart:       time.Date(2021, 06, 16, 20, 30, 0, 0, time.UTC),
				DateEnd:         time.Date(2021, 06, 21, 21, 0, 0, 0, time.UTC),
			},

			Occurrences: []train.CreateOccRequest{
				{
					OccRequest: train.OccRequest{
						// saturday:6
						DateStart:  time.Date(2021, 05, 20, 20, 30, 0, 0, time.UTC),
						DateEnd:    time.Date(2021, 05, 20, 21, 0, 0, 0, time.UTC),
						RepeatDays: 10,
					},
				},
			},
		}),
		ExpectedStatusCode: 204,
		RequestHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})

	return searchCtx.Api.TestAssertCasesSemaphoreErr(tc, runtime.NumCPU())
}

func matchResults(d []search.ScheduleWithScore, msg string, ts ...string) error {
	titles := make([]string, len(d))
	for i := range d {
		titles[i] = d[i].Training.Title
	}
	if err := helpers.AssertErr(
		fmt.Sprintf("%s: invalid number of records: [%s]", msg, strings.Join(titles, ",")),
		len(ts), len(d)); err != nil {
		return err
	}

	for _, t := range ts {
		found := false
		for j := range d {
			if d[j].Training.Title == "" {
				return fmt.Errorf("empty training returned")
			}
			if d[j].Training.Title == t {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("%s: title not found: %s in data [%s]", msg, t, strings.Join(titles, ","))
		}
	}
	return nil
}

func mustMatchTrainings(msg string, ts ...string) func(b []byte, i interface{}) error {
	return func(b []byte, i interface{}) error {
		var r search.SearchResultsWithMetadata
		if err := json.Unmarshal(b, &r); err != nil {
			return err
		}
		return matchResults(r.Data, msg, ts...)
	}
}

func TestSearch(t *testing.T) {

	if err := createTestTrainings(); err != nil {
		t.Fatal(err)
	}

	if err := searchCtx.RegenerateCache(); err != nil {
		t.Fatal(err)
	}

	tc := []api.TestCase{
		// no request
		{
			RequestMethod:      "POST",
			RequestUrl:         searchPath,
			ExpectedStatusCode: 400,
		},
		// empty request
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs: []string{"en"},
				},
			}),
			ExpectedStatusCode: 400,
		},
		//
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"en"},
					Country: "PL",
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 0", "test1", "test2", "test4"),
		},
		//
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs: []string{"en"},
					// act distance betwee 30.5 and 31 is around 73km
					DistKm: 80,
					Lat:    30.5,
					Lng:    30.5,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 1", "test2", "test3"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs: []string{"en"},
					// 100 km is very roughly 1 dg lat / lng
					DistKm: 80,
					Lat:    29.5,
					Lng:    29.5,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 2", "test2", "test3"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs: []string{"en"},
					// 100 km is very roughly 1 dg lat / lng
					DistKm: 80,
					Lat:    30,
					Lng:    29.5,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 3", "test2", "test3"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs: []string{"en"},
					// 100 km is very roughly 1 dg lat / lng
					DistKm: 80,
					Lat:    29.5,
					Lng:    30,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 4", "test2", "test3"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs: []string{"en"},
					// 100 km is very roughly 1 dg lat / lng
					DistKm: 150,
					Lat:    32,
					Lng:    32,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 5"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs: []string{"en", "pl"},
					// 100 km is very roughly 1 dg lat / lng
					DistKm:       160,
					Lat:          31,
					Lng:          31,
					MinTagPoints: 0.001,
				},
				Query: "tag",
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 6", "test3"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:       []string{"en"},
					Country:     "PL",
					CapacityMin: 0,
					CapacityMax: 9,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 7", "test1", "test4"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:       []string{"en"},
					Country:     "PL",
					CapacityMin: 0,
					CapacityMax: 9,
					Age:         11,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 7a", "test1"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:       []string{"en"},
					Country:     "PL",
					CapacityMin: 0,
					CapacityMax: 9,
					Age:         12,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 7b", "test1", "test4"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:       []string{"en"},
					Country:     "PL",
					CapacityMin: 0,
					CapacityMax: 9,
					Age:         15,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 7c", "test4"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:    []string{"en"},
					Country:  "PL",
					PriceMin: 9000,
					PriceMax: 10000,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 8", "test2", "test4"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"en"},
					Country: "DE",
					Diffs: []int{
						3,
					},
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 9", "test5", "test6"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"en"},
					Country: "DE",
					Diffs: []int{
						1, 1,
					},
				},
			}),
			ExpectedStatusCode: 400,
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"en"},
					Country: "DE",
					Diffs: []int{
						1,
					},
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 10", "test5"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"en"},
					Country: "DE",
					HrStart: "14:00",
					HrEnd:   "15:00",
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 11", "test6", "test7"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"en"},
					Country: "DE",
					HrStart: "blehasda",
					HrEnd:   "15:00",
				},
			}),
			ExpectedStatusCode: 400,
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"en"},
					Country: "DE",
					HrStart: "14:00",
					HrEnd:   "15;00",
				},
			}),
			ExpectedStatusCode: 400,
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"en"},
					Country: "DE",
					Days: []search.Weekday{
						search.Friday,
					},
					OmitEmptySchedules: true,
					DateStart:          helpers.NowMin(),
					DateEnd:            helpers.NowMin().Add(time.Hour * 24 * 100),
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 12", "test6", "test10", "test8"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:              []string{"en"},
					Country:            "DE",
					OmitEmptySchedules: true,
					Days: []search.Weekday{
						search.Monday, search.Tuesday,
					},
					DateStart: helpers.NowMin(),
					DateEnd:   helpers.NowMin().Add(time.Hour * 24 * 100),
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 13", "test6", "test7", "test10", "test8"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:       []string{"en"},
					Country:     "DE",
					DurationMin: 60,
					DurationMax: 60,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 14", "test6"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:       []string{"en"},
					Country:     "DE",
					DurationMin: 120,
					DurationMax: 240,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 15", "test7", "test6"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:              []string{"en"},
					Country:            "DE",
					DateStart:          time.Date(2021, 2, 8, 0, 0, 0, 0, time.UTC),
					DateEnd:            time.Date(2021, 2, 9, 23, 59, 0, 0, time.UTC),
					OmitEmptySchedules: true,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 16", "test6", "test7"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:              []string{"en"},
					Country:            "DE",
					DateStart:          time.Date(2021, 2, 8, 0, 0, 0, 0, time.UTC),
					DateEnd:            time.Date(2021, 2, 8, 23, 59, 0, 0, time.UTC),
					OmitEmptySchedules: true,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 17", "test6"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:              []string{"en"},
					Country:            "DE",
					DateStart:          time.Date(2021, 2, 23, 0, 0, 0, 0, time.UTC),
					DateEnd:            time.Date(2021, 2, 23, 23, 59, 0, 0, time.UTC),
					OmitEmptySchedules: true,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 18", "test8", "test7"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:              []string{"en"},
					Country:            "DE",
					DateStart:          time.Date(2021, 2, 22, 0, 0, 0, 0, time.UTC),
					DateEnd:            time.Date(2021, 2, 22, 23, 59, 0, 0, time.UTC),
					OmitEmptySchedules: true,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 19", "test6"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:              []string{"en"},
					Country:            "DE",
					DateStart:          time.Date(2021, 2, 24, 0, 0, 0, 0, time.UTC),
					DateEnd:            time.Date(2021, 2, 24, 23, 59, 0, 0, time.UTC),
					OmitEmptySchedules: true,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 20"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:              []string{"en"},
					Country:            "DE",
					DateStart:          time.Date(2021, 4, 19, 0, 0, 0, 0, time.UTC),
					DateEnd:            time.Date(2021, 4, 20, 23, 59, 0, 0, time.UTC),
					OmitEmptySchedules: true,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 21", "test9", "test6", "test7"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:              []string{"en"},
					Country:            "DE",
					DateStart:          time.Date(2021, 5, 30, 0, 0, 0, 0, time.UTC),
					DateEnd:            time.Date(2021, 5, 30, 23, 59, 0, 0, time.UTC),
					OmitEmptySchedules: true,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 22", "test10", "test8"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:              []string{"en"},
					Country:            "DE",
					DateStart:          time.Date(2021, 6, 9, 0, 0, 0, 0, time.UTC),
					DateEnd:            time.Date(2021, 6, 9, 23, 59, 0, 0, time.UTC),
					OmitEmptySchedules: true,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 23", "test10"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:              []string{"en"},
					Country:            "DE",
					DateStart:          time.Date(2021, 6, 9, 0, 0, 0, 0, time.UTC),
					DateEnd:            time.Date(2021, 6, 9, 23, 59, 0, 0, time.UTC),
					OmitEmptySchedules: true,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 24", "test10"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:              []string{"en"},
					Country:            "DE",
					DateStart:          time.Date(2021, 6, 19, 0, 0, 0, 0, time.UTC),
					DateEnd:            time.Date(2021, 6, 19, 23, 59, 0, 0, time.UTC),
					OmitEmptySchedules: true,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 25", "test10", "test11"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:              []string{"en"},
					Country:            "DE",
					DateStart:          time.Date(2021, 6, 29, 0, 0, 0, 0, time.UTC),
					DateEnd:            time.Date(2021, 6, 29, 23, 59, 0, 0, time.UTC),
					OmitEmptySchedules: true,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal:    mustMatchTrainings("search 25", "test10", "test8", "test7"),
		},
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:              []string{"en"},
					Country:            "DE",
					DateStart:          time.Date(2021, 6, 29, 0, 0, 0, 0, time.UTC),
					DateEnd:            time.Date(2022, 5, 29, 23, 59, 0, 0, time.UTC),
					OmitEmptySchedules: true,
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var sr search.SearchResultsWithMetadata
				helpers.JsonMustDeserialize(b, &sr)
				fmt.Println(helpers.JsonMustSerializeFormatStr(sr.Meta))
				for i := range sr.Data {
					fmt.Println(helpers.JsonMustSerializeFormatStr(sr.Data[i].Score))
				}
				return nil
			},
		},
	}

	// for k := range search.SortExp {
	// 	tc = append(tc, api.TestCase{
	// 		RequestMethod: "POST",
	// 		RequestUrl:    path,
	// 		RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
	// 			SearchRequest: search.SearchRequest{
	// 				Langs:     []string{"en"},
	// 				Country: "DE",
	// 				Sort: []search.SearchSortRecord{
	// 					{
	// 						Column: k,
	// 						IsDesc: true,
	// 					},
	// 				},
	// 			},
	// 		}),
	// 		ExpectedStatusCode: 200,
	// 	})
	// 	tc = append(tc, api.TestCase{
	// 		RequestMethod: "POST",
	// 		RequestUrl:    path,
	// 		RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
	// 			SearchRequest: search.SearchRequest{
	// 				Langs:     []string{"en"},
	// 				Country: "DE",
	// 				Sort: []search.SearchSortRecord{
	// 					{
	// 						Column: k,
	// 						IsDesc: false,
	// 					},
	// 				},
	// 			},
	// 		}),
	// 		ExpectedStatusCode: 200,
	// 	})
	// }

	searchCtx.Api.TestAssertCases(t, tc)
}
