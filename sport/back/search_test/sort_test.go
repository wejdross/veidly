package search_test

import (
	"encoding/json"
	"fmt"
	"sport/api"
	"sport/helpers"
	"sport/search"
	"sport/train"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSort(t *testing.T) {
	token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
	if err != nil {
		t.Fatal(err)
	}

	// create training
	s := time.Now()
	trs := []train.CreateTrainingRequest{
		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
	}

	country := "HU"

	for i := range trs {
		trs[i].Training.Title = fmt.Sprintf("sortTest%d", i)
		trs[i].Training.LocationCountry = country
		trs[i].Training.Price = (len(trs) - i + 1) * 100
		trs[i].Training.Capacity = i + 10
	}

	//ts := make(search.TestListenerMap)
	tids := make([]uuid.UUID, len(trs))

	for i := range trs {
		_tid, err := trainCtx.ApiCreateTraining(token, &trs[i])
		if err != nil {
			t.Fatal(err)
		}
		tids[i] = _tid
	}

	if err := searchCtx.RegenerateCache(); err != nil {
		t.Fatal(err)
	}

	searchCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"pl"},
					Country: country,
					Sort: []search.SearchSortRecord{
						{
							Column: search.Price,
							IsDesc: false,
						},
					},
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var r search.SearchResultsWithMetadata
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}
				for i := range r.Data {
					p := r.Data[i].Training.Price
					if i == len(r.Data)-1 {
						continue
					}
					np := r.Data[i+1].Training.Price
					if np < p {
						return fmt.Errorf("search 0 - not sorted")
					}
				}
				return mustMatchTrainings("search 0", "sortTest0", "sortTest1", "sortTest2")(b, i)
			},
		}, {
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"pl"},
					Country: country,
					Sort: []search.SearchSortRecord{
						{
							Column: search.Price,
							IsDesc: true,
						},
					},
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var r search.SearchResultsWithMetadata
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}
				for i := range r.Data {
					p := r.Data[i].Training.Price
					if i == len(r.Data)-1 {
						continue
					}
					np := r.Data[i+1].Training.Price
					if np > p {
						return fmt.Errorf("search 1 - not sorted")
					}
				}
				return mustMatchTrainings("search 1", "sortTest0", "sortTest1", "sortTest2")(b, i)
			},
		}, {
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"pl"},
					Country: country,
					Sort: []search.SearchSortRecord{
						{
							Column: search.Capacity,
							IsDesc: false,
						},
					},
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var r search.SearchResultsWithMetadata
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}
				for i := range r.Data {
					p := r.Data[i].Training.Capacity
					if i == len(r.Data)-1 {
						continue
					}
					np := r.Data[i+1].Training.Capacity
					if np < p {
						return fmt.Errorf("search 2 - not sorted")
					}
				}
				return mustMatchTrainings("search 2", "sortTest0", "sortTest1", "sortTest2")(b, i)
			},
		}, {
			RequestMethod: "POST",
			RequestUrl:    searchPath,
			RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
				SearchRequestOpts: search.SearchRequestOpts{
					Langs:   []string{"pl"},
					Country: country,
					Sort: []search.SearchSortRecord{
						{
							Column: search.Capacity,
							IsDesc: true,
						},
					},
				},
			}),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var r search.SearchResultsWithMetadata
				if err := json.Unmarshal(b, &r); err != nil {
					return err
				}
				for i := range r.Data {
					p := r.Data[i].Training.Capacity
					if i == len(r.Data)-1 {
						continue
					}
					np := r.Data[i+1].Training.Capacity
					if np > p {
						return fmt.Errorf("search 3 - not sorted")
					}
				}
				return mustMatchTrainings("search 3", "sortTest0", "sortTest1", "sortTest2")(b, i)
			},
		},
	})
}
