package schedule_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"sport/adyen_sm"
// 	"sport/api"
// 	"sport/helpers"
// 	"sport/instr"
// 	"sport/rsv"
// 	"sport/schedule"
// 	"sport/train"
// 	"sport/user"
// 	"testing"
// 	"time"

// 	"github.com/google/uuid"
// )

// func TestSchedule(t *testing.T) {

// 	instructorToken, err := rsvCtx.User.ApiCreateAndLoginUser(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	userToken, err := rsvCtx.User.ApiCreateAndLoginUser(&user.UserRequest{
// 		Email:    "1@1.1",
// 		Password: "!@213SADsadasdASD213",
// 		UserData: user.UserData{
// 			Language: "en",
// 			Country:  "US",
// 			Name:     "foo bar",
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	var userID uuid.UUID

// 	if userID, err = rsvCtx.Api.AuthorizeUserFromToken(userToken); err != nil {
// 		t.Fatal(err)
// 	}

// 	var tres *train.TrainingWithJoins
// 	var emptyRes *train.TrainingWithJoins
// 	//var instructorID uuid.UUID

// 	{
// 		instr := instr.InstructorRequest{}
// 		trainingReq := train.CreateTrainingRequest{
// 			Training: train.TrainingRequest{
// 				Title:    "GKbwB8qlSVqifs64ZXfL",
// 				Capacity: 1,
// 				Currency: "PLN",
// 				Price:    100 * 100,
// 			},
// 			Occurrences: []train.CreateOccRequest{
// 				{
// 					OccRequest: train.OccRequest{
// 						DateStart:  time.Date(2028, 02, 01, 8, 0, 0, 0, time.UTC),
// 						DateEnd:    time.Date(2028, 02, 01, 10, 0, 0, 0, time.UTC),
// 						RepeatDays: 7, // weekly
// 						Remarks:    "9KCwBNJiScZmbQNUssIxT1z4IXWjJOHTL2LumZzVXslG9iIDkdnTjFSXO0czygL38UCHHetIBHXWvzzycJCnxCWuvSrTLoF4MbVhF3Pq6sJgIYKbWEsIf7Y1sOHNJJUX",
// 						Color:      "dpD5kB1Ju",
// 					},
// 				},
// 			},
// 		}
// 		trainingReq2 := train.CreateTrainingRequest{
// 			Training: train.TrainingRequest{
// 				Title:    "Z5MK47WSXevXTChMGDB8",
// 				Capacity: 10,
// 				Currency: "PLN",
// 				Price:    100 * 40,
// 			},
// 			Occurrences: []train.CreateOccRequest{
// 				{
// 					OccRequest: train.OccRequest{
// 						DateStart:  time.Date(2028, 02, 01, 7, 0, 0, 0, time.UTC),
// 						DateEnd:    time.Date(2028, 02, 01, 8, 0, 0, 0, time.UTC),
// 						RepeatDays: 1, // daily
// 						Remarks:    "0RQkyXP6WeozhEEWJXIoEKAfh9YIQSGPDelZO7Lhf9RDn1D4oLHaCpSijMkezbMAyxJaWtzEYB4nutMxkNx8wPewJ0HkBKJKoPHWhu1uuutDSm241netJuKgqudXQOA6",
// 						Color:      "tTfLWSr5u",
// 					},
// 				},
// 				{
// 					OccRequest: train.OccRequest{
// 						DateStart:  time.Date(2028, 02, 01, 6, 0, 0, 0, time.UTC),
// 						DateEnd:    time.Date(2028, 02, 01, 7, 0, 0, 0, time.UTC),
// 						RepeatDays: 1, // daily
// 						Remarks:    "IXjw6MjpJwsVwuxaZqcrWXwBuqrMHYmSTCs17BWyrXFBD8sqOnicgJbAbQxMqZs9YLArqM85Gi5Vkkq9YyGH1kOdJyOvc0xqJDBwFkOIqRPuVGSaqbgaqewDXxH5RLqv",
// 						Color:      "nEbx1EaXp",
// 					},
// 				},
// 			},
// 		}
// 		trainingReq3 := train.CreateTrainingRequest{
// 			Training: train.TrainingRequest{
// 				Title:    "saaaHLlYFq9VaXPgebeK",
// 				Capacity: 10,
// 				Currency: "PLN",
// 				Price:    100 * 40,
// 			},
// 		}
// 		cases := []api.TestCase{
// 			{
// 				RequestMethod:      "POST",
// 				RequestUrl:         "/api/instructor",
// 				RequestReader:      helpers.JsonMustSerializeReader(instr),
// 				ExpectedStatusCode: 204,
// 				RequestHeaders:     api.JwtHeader(instructorToken),
// 			},
// 			{
// 				RequestMethod:      "POST",
// 				RequestUrl:         "/api/training",
// 				RequestReader:      helpers.JsonMustSerializeReader(trainingReq),
// 				ExpectedStatusCode: 204,
// 				RequestHeaders:     api.JwtHeader(instructorToken),
// 			},
// 			{
// 				RequestMethod:      "POST",
// 				RequestUrl:         "/api/training",
// 				RequestReader:      helpers.JsonMustSerializeReader(trainingReq2),
// 				ExpectedStatusCode: 204,
// 				RequestHeaders:     api.JwtHeader(instructorToken),
// 			},
// 			{
// 				RequestMethod:      "POST",
// 				RequestUrl:         "/api/training",
// 				RequestReader:      helpers.JsonMustSerializeReader(trainingReq3),
// 				ExpectedStatusCode: 204,
// 				RequestHeaders:     api.JwtHeader(instructorToken),
// 			},
// 			{
// 				RequestMethod:      "GET",
// 				RequestUrl:         "/api/training",
// 				ExpectedStatusCode: 200,
// 				RequestHeaders:     api.JwtHeader(instructorToken),
// 				ExpectedBodyVal: func(b []byte, i interface{}) error {
// 					var tmp []train.TrainingWithJoins
// 					var err error
// 					if err = json.Unmarshal(b, &tmp); err != nil {
// 						return err
// 					}
// 					if err = helpers.AssertErr("get instructor training", len(tmp), 3); err != nil {
// 						return err
// 					}
// 					tres = &tmp[0]
// 					for i := range tmp {
// 						if tmp[i].Training.Title == "GKbwB8qlSVqifs64ZXfL" {
// 							tres = &tmp[i]
// 						}
// 						//instructorID = tmp[i].Training.InstructorID
// 						if len(tmp[i].Occurrences) == 0 {
// 							emptyRes = &tmp[i]
// 						}
// 					}
// 					if emptyRes == nil {
// 						return fmt.Errorf("didnt find empty training in response")
// 					}
// 					if tres == nil {
// 						return fmt.Errorf("didnt find tres training in response")
// 					}
// 					return nil
// 				},
// 			},
// 		}
// 		rsvCtx.Api.TestAssertCases(t, cases)
// 	}

// 	// create reservation
// 	okRsvRequest := rsv.ApiReservationRequest{
// 		TrainingID: tres.Training.ID,
// 		// take correct date start, then add random amount of RepeatDays (ex. weeks), to match training occurrence
// 		Occurrence: time.Date(2028, 2, 8, 8, 0, 0, 0, time.UTC),
// 		UserData: user.UserData{
// 			Language: "en",
// 			Country:  "US",
// 			Name:     "vivec",
// 		},
// 		UseSavedData: true,
// 	}

// 	{
// 		cases := []api.TestCase{
// 			// invalid request (date doesnt match occurrence)
// 			{
// 				RequestMethod: "POST",
// 				RequestUrl:    "/api/rsv",
// 				RequestReader: helpers.JsonMustSerializeReader(&rsv.ApiReservationRequest{
// 					TrainingID: tres.Training.ID,
// 					Occurrence: time.Date(2028, 12, 1, 25, 0, 0, 0, time.UTC),
// 					UserData: user.UserData{
// 						Language: "en",
// 						Country:  "US",
// 					},
// 				}),
// 				ExpectedStatusCode: 409,
// 				RequestHeaders:     api.JwtHeader(userToken),
// 			},
// 			{
// 				RequestMethod: "POST",
// 				RequestUrl:    "/api/rsv",
// 				RequestReader: helpers.JsonMustSerializeReader(&rsv.ApiReservationRequest{
// 					TrainingID: emptyRes.Training.ID,
// 					Occurrence: time.Date(2028, 12, 1, 25, 0, 0, 0, time.UTC),
// 					UserData: user.UserData{
// 						Language: "en",
// 						Country:  "US",
// 						Name:     "foo bar",
// 					},
// 				}),
// 				ExpectedStatusCode: 409,
// 				RequestHeaders:     api.JwtHeader(userToken),
// 			},
// 			//rsv with logged in user
// 			{
// 				RequestMethod:      "POST",
// 				RequestUrl:         "/api/rsv",
// 				RequestReader:      helpers.JsonMustSerializeReader(&okRsvRequest),
// 				ExpectedStatusCode: 303,
// 				RequestHeaders:     api.JwtHeader(userToken),
// 				ExpectedHeadersVal: func(m map[string][]string) error {
// 					c := m["Location"]
// 					if len(c) != 1 {
// 						return fmt.Errorf("invalid location hdr in response")
// 					}
// 					//fmt.Println(c[0])
// 					return nil
// 				},
// 				ExpectedBodyVal: func(b []byte, i interface{}) error {
// 					var res rsv.PostRsvResponse
// 					helpers.JsonMustDeserialize(b, &res)
// 					if res.ID == uuid.Nil {
// 						return fmt.Errorf("empty id returned")
// 					}
// 					return nil
// 				},
// 			},
// 		}
// 		rsvCtx.Api.TestAssertCases(t, cases)
// 	}

// 	//try to timeout

// 	//create mockup payu notify request
// 	var r *rsv.DDLRsvWithInstr
// 	if r, err = rsvCtx.ReadSingleRsv(rsv.ReadRsvsArgs{UserID: &userID, WithInstructor: true}); err != nil {
// 		t.Fatal(err)
// 	}

// 	if err := rsvCtx.AdyenSm.MoveStateToWaitCapture(
// 		rsvCtx.RsvResponseToSmPassPtr(r),
// 		adyen_sm.ManualSrc,
// 		"", false,
// 	); err != nil {
// 		t.Fatal(err)
// 	}

// 	r.State = adyen_sm.WaitCapture
// 	// confirm reservation

// 	if err := rsvCtx.AdyenSm.MoveStateToCapture(
// 		rsvCtx.RsvResponseToSmPassPtr(r), adyen_sm.ManualSrc, nil,
// 	); err != nil {
// 		t.Fatal(err)
// 	}

// 	// instructorID = r.Training.InstructorID

// 	// rtf1, err := ioutil.ReadFile("../rsv_test/test_files/1.json")
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// 	// rtf2, err := ioutil.ReadFile("../rsv_test/test_files/2.json")
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// 	// rtf3, err := ioutil.ReadFile("../rsv_test/test_files/3.json")
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }

// 	// rtf1s := string(rtf1)
// 	// rtf2s := string(rtf2)
// 	// rtf3s := string(rtf3)

// 	// trs, err := rsvCtx.Train.DalReadTrainings(train.DalReadTrainingsRequest{
// 	// 	InstructorID: &instructorID,
// 	// 	WithOccs:     true,
// 	// 	WithGroups:   true,
// 	// })
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }

// 	// rf := func(f *string) {
// 	// 	*f = strings.Replace(*f, "$TRAINING1$", trs[0].Training.ID.String(), -1)
// 	// 	*f = strings.Replace(*f, "$TRAINING2$", trs[1].Training.ID.String(), -1)
// 	// 	*f = strings.Replace(*f, "$TRAINING3$", trs[2].Training.ID.String(), -1)
// 	// 	*f = strings.Replace(*f, "$CREATED_ON3$", trs[2].Training.CreatedOn.Format("2006-01-02T15:04:05Z0700"), -1)
// 	// 	*f = strings.Replace(*f, "$USER_ID$", userID.String(), -1)
// 	// 	*f = strings.Replace(*f, "$CREATED_ON$", trs[0].Training.CreatedOn.Format("2006-01-02T15:04:05Z0700"), -1)
// 	// 	*f = strings.Replace(*f, "$INSTRUCTOR$", instructorID.String(), -1)
// 	// 	*f = strings.Replace(*f, "$TRAINING1_SESSION1$", trs[0].Occurrences[0].ID.String(), -1)
// 	// 	*f = strings.Replace(*f, "$TRAINING2_SESSION1$", trs[1].Occurrences[0].ID.String(), -1)
// 	// 	*f = strings.Replace(*f, "$TRAINING2_SESSION2$", trs[1].Occurrences[1].ID.String(), -1)
// 	// 	*f = strings.Replace(*f, "$LINKID$", r.LinkID, -1)
// 	// 	*f = strings.Replace(*f, "$RSV$", r.ID.String(), -1)
// 	// 	*f = strings.Replace(*f, "$CREATED_ON_RSV$", r.CreatedOn.Format("2006-01-02T15:04:05Z0700"), -1)
// 	// }

// 	// rf(&rtf1s)
// 	// rf(&rtf2s)
// 	// rf(&rtf3s)

// 	cases := []api.TestCase{
// 		{
// 			RequestMethod:      "GET",
// 			RequestUrl:         "/api/rsv/t/user",
// 			ExpectedStatusCode: 401,
// 		},
// 		// read rsv
// 		{
// 			RequestMethod:      "GET",
// 			RequestUrl:         "/api/rsv/t/user?page=0&size=10",
// 			ExpectedStatusCode: 200,
// 			RequestHeaders:     api.JwtHeader(userToken),
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				var res rsv.RsvWithInstrPagination
// 				if err = json.Unmarshal(b, &res); err != nil {
// 					return err
// 				}
// 				return rsvCtx.AssertRsvResponse(
// 					&res,
// 					[]*rsv.ApiReservationRequest{&okRsvRequest},
// 					&userID,
// 					1,
// 					0,
// 					10)
// 			},
// 		},
// 		{
// 			RequestMethod:      "GET",
// 			RequestUrl:         "/api/rsv/t/instructor?page=0&size=10",
// 			ExpectedStatusCode: 200,
// 			RequestHeaders:     api.JwtHeader(instructorToken),
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				var res rsv.RsvWithInstrPagination
// 				if err = json.Unmarshal(b, &res); err != nil {
// 					return err
// 				}
// 				return rsvCtx.AssertRsvResponse(
// 					&res,
// 					[]*rsv.ApiReservationRequest{&okRsvRequest},
// 					&userID,
// 					1,
// 					0,
// 					10)
// 			},
// 		},
// 		{
// 			RequestMethod:      "GET",
// 			RequestUrl:         "/api/rsv/t/user?page=1&size=10",
// 			ExpectedStatusCode: 200,
// 			RequestHeaders:     api.JwtHeader(userToken),
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				var res rsv.RsvWithInstrPagination
// 				if err = json.Unmarshal(b, &res); err != nil {
// 					return err
// 				}
// 				return rsvCtx.AssertRsvResponse(
// 					&res,
// 					nil,
// 					&userID,
// 					1,
// 					1,
// 					10)
// 			},
// 		},
// 		{
// 			RequestMethod: "GET",
// 			RequestUrl: fmt.Sprintf("/api/rsv/t/user?start=%d&end=%d",
// 				okRsvRequest.Occurrence.Unix(), okRsvRequest.Occurrence.Unix()),
// 			ExpectedStatusCode: 200,
// 			RequestHeaders:     api.JwtHeader(userToken),
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				var res rsv.RsvWithInstrPagination
// 				if err = json.Unmarshal(b, &res); err != nil {
// 					return err
// 				}
// 				return rsvCtx.AssertRsvResponse(
// 					&res,
// 					[]*rsv.ApiReservationRequest{&okRsvRequest},
// 					&userID,
// 					1,
// 					0,
// 					rsvCtx.Config.MaxPageSize)
// 			},
// 		},
// 		{
// 			RequestMethod: "GET",
// 			RequestUrl: fmt.Sprintf("/api/rsv/t/user?start=%d&end=%d",
// 				okRsvRequest.Occurrence.Add(-2*time.Second).Unix(), okRsvRequest.Occurrence.Add(-1*time.Second).Unix()),
// 			ExpectedStatusCode: 200,
// 			RequestHeaders:     api.JwtHeader(userToken),
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				var res rsv.RsvWithInstrPagination
// 				if err = json.Unmarshal(b, &res); err != nil {
// 					return err
// 				}
// 				return rsvCtx.AssertRsvResponse(
// 					&res,
// 					[]*rsv.ApiReservationRequest{},
// 					&userID,
// 					0,
// 					0,
// 					rsvCtx.Config.MaxPageSize)
// 			},
// 		},
// 		{
// 			RequestMethod: "GET",
// 			RequestUrl: fmt.Sprintf("/api/rsv/t/user?start=%d&end=%d",
// 				okRsvRequest.Occurrence.Add(-2*time.Second).Unix(), okRsvRequest.Occurrence.Add(time.Hour*100).Unix()),
// 			ExpectedStatusCode: 200,
// 			RequestHeaders:     api.JwtHeader(userToken),
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				var res rsv.RsvWithInstrPagination
// 				if err = json.Unmarshal(b, &res); err != nil {
// 					return err
// 				}
// 				return rsvCtx.AssertRsvResponse(
// 					&res,
// 					[]*rsv.ApiReservationRequest{&okRsvRequest},
// 					&userID,
// 					1,
// 					0,
// 					rsvCtx.Config.MaxPageSize)
// 			},
// 		},
// 		{
// 			RequestMethod: "GET",
// 			RequestUrl: fmt.Sprintf("/api/schedule?start=%d&end=%d&instructor_id=%s",
// 				time.Date(2028, 02, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 				time.Date(2028, 03, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 				r.InstructorID),
// 			ExpectedStatusCode: 200,
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				var res []schedule.TrainingSchedule
// 				if err = json.Unmarshal(b, &res); err != nil {
// 					return err
// 				}
// 				//resJson := helpers.JsonMustSerializeFormatStr(res)
// 				// if resJson != rtf1s {
// 				// 	if err := ioutil.WriteFile("../rsv_test/test_files/result/expected.json", []byte(rtf1s), 0600); err != nil {
// 				// 		return err
// 				// 	}
// 				// 	if err := ioutil.WriteFile("../rsv_test/test_files/result/got.json", []byte(resJson), 0600); err != nil {
// 				// 		return err
// 				// 	}
// 				// 	return fmt.Errorf("invalid rsv 1 schedule")
// 				// }
// 				return nil
// 			},
// 		},
// 		// delete training
// 		{
// 			RequestMethod: "DELETE",
// 			RequestUrl:    "/api/training",
// 			RequestHeaders: map[string]string{
// 				"Authorization": "Bearer " + instructorToken,
// 			},
// 			RequestReader: helpers.JsonMustSerializeReader(train.ObjectKey{
// 				ID: okRsvRequest.TrainingID,
// 			}),
// 			ExpectedStatusCode: 204,
// 		},
// 		// now rsv should be still there but it should be orpaned.
// 		// note that this should only show when we are searching as instructor
// 		// other users shouldnt be able to see this instructor's orphaned rsvs
// 		// im setting flags = 1 (NoRsvs) so we dont generate reservations for schedule
// 		{
// 			RequestMethod: "GET",
// 			RequestUrl: fmt.Sprintf("/api/schedule?start=%d&end=%d&flags=1",
// 				time.Date(2028, 02, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 				time.Date(2028, 03, 1, 0, 0, 0, 0, time.UTC).Unix()),
// 			ExpectedStatusCode: 200,
// 			RequestHeaders: map[string]string{
// 				"Authorization": "Bearer " + instructorToken,
// 			},
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				var res []schedule.TrainingSchedule
// 				if err = json.Unmarshal(b, &res); err != nil {
// 					return err
// 				}
// 				// resJson := helpers.JsonMustSerializeFormatStr(res)
// 				// if resJson != rtf2s {
// 				// 	if err := ioutil.WriteFile("../rsv_test/test_files/result/expected.json", []byte(rtf2s), 0600); err != nil {
// 				// 		return err
// 				// 	}
// 				// 	if err := ioutil.WriteFile("../rsv_test/test_files/result/got.json", []byte(resJson), 0600); err != nil {
// 				// 		return err
// 				// 	}
// 				// 	return fmt.Errorf("invalid rsv 2 schedule")
// 				// }
// 				return nil
// 			},
// 		},
// 		{
// 			RequestMethod: "GET",
// 			RequestUrl: fmt.Sprintf("/api/schedule?start=%d&end=%d",
// 				time.Date(2028, 02, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 				time.Date(2028, 03, 1, 0, 0, 0, 0, time.UTC).Unix()),
// 			RequestHeaders:     api.JwtHeader(instructorToken),
// 			ExpectedStatusCode: 200,
// 			ExpectedBodyVal: func(b []byte, i interface{}) error {
// 				var res []schedule.TrainingSchedule
// 				if err = json.Unmarshal(b, &res); err != nil {
// 					return err
// 				}
// 				// resJson := helpers.JsonMustSerializeFormatStr(res)
// 				// if resJson != rtf3s {
// 				// 	if err := ioutil.WriteFile("../rsv_test/test_files/result/expected.json", []byte(rtf3s), 0600); err != nil {
// 				// 		return err
// 				// 	}
// 				// 	if err := ioutil.WriteFile("../rsv_test/test_files/result/got.json", []byte(resJson), 0600); err != nil {
// 				// 		return err
// 				// 	}
// 				// 	return fmt.Errorf("invalid rsv 3 schedule")
// 				// }
// 				return nil
// 			},
// 		},
// 	}
// 	rsvCtx.Api.TestAssertCases(t, cases)

// }
