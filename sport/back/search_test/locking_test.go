package search_test

// func TestLock(t *testing.T) {

// 	token, _, err := trainCtx.Instr.CreateTestInstructorWithUser()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// create training
// 	s := time.Now()
// 	trs := []train.CreateTrainingRequest{
// 		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
// 		train.NewTestCreateTrainingRequest(s, s.Add(time.Minute)),
// 	}

// 	country := "AD"

// 	for i := range trs {
// 		trs[i].Training.Title = fmt.Sprintf("locktest%d", i)
// 		trs[i].Training.LocationCountry = country
// 	}

// 	tids := make([]uuid.UUID, len(trs))

// 	for i := range trs {
// 		_tid, err := trainCtx.ApiCreateTraining(token, &trs[i])
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		tids[i] = _tid
// 	}

// 	//

// 	c := make(chan struct{})
// 	oldiv := searchCtx.CacheRefreshInterval
// 	searchCtx.CacheRefreshInterval = time.Millisecond * 200
// 	/* test starts with RWLock unlocked, so gotta lock it to prevent panic */
// 	searchCtx.RWLock.Lock()
// 	go searchCtx.RunCacheGenerationAgent(c)
// 	searchCtx.DbgSlowdownCacheRefresh = time.Millisecond * 400

// 	// waiting for cache to be generated
// 	<-c

// 	defer func() {
// 		searchCtx.CacheRefreshInterval = oldiv
// 		searchCtx.DbgSlowdownCacheRefresh = 0
// 		/* terminate cache generation agent */
// 		c <- struct{}{}
// 	}()

// 	count := 40
// 	td := make([]time.Duration, count)

// 	for i := 0; i < count; i++ {
// 		s := time.Now()
// 		searchCtx.Api.TestAssertCases(t, []api.TestCase{
// 			{
// 				RequestMethod: "POST",
// 				RequestUrl:    searchPath,
// 				RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
// 					SearchRequestOpts: search.SearchRequestOpts{
// 						Langs:   []string{"pl"},
// 						Country: country,
// 					},
// 				}),
// 				ExpectedStatusCode: 200,
// 				ExpectedBodyVal:    mustMatchTrainings("search "+strconv.Itoa(i), "locktest0", "locktest1"),
// 			},
// 		})
// 		td[i] = time.Since(s)
// 		time.Sleep(time.Millisecond * 10)
// 	}

// 	f := 0

// 	for i := range td {
// 		if td[i] >= (searchCtx.DbgSlowdownCacheRefresh - 20*time.Millisecond) {
// 			f++
// 		}
// 	}

// 	if f == 0 {
// 		t.Fatal("locking isnt working")
// 	} else {
// 		fmt.Printf("threads [%d] have been throttled\n", f)
// 	}

// 	count = 5000
// 	td = make([]time.Duration, count)

// 	if err := helpers.Sem(20, count, []func(i int) error{
// 		func(i int) error {
// 			s := time.Now()
// 			if err := searchCtx.Api.TestAssertCaseErr(&api.TestCase{
// 				RequestMethod: "POST",
// 				RequestUrl:    searchPath,
// 				RequestReader: helpers.JsonMustSerializeReader(search.ApiSearchRequest{
// 					SearchRequestOpts: search.SearchRequestOpts{
// 						Langs:   []string{"pl"},
// 						Country: country,
// 					},
// 				}),
// 				ExpectedStatusCode: 200,
// 				ExpectedBodyVal:    mustMatchTrainings("search "+strconv.Itoa(i), "locktest0", "locktest1"),
// 			}); err != nil {
// 				return err
// 			}
// 			td[i] = time.Since(s)
// 			return nil
// 		},
// 	}); err != nil {
// 		t.Fatal(err)
// 	}

// 	f = 0

// 	for i := range td {
// 		if td[i] >= (searchCtx.DbgSlowdownCacheRefresh - 20*time.Millisecond) {
// 			f++
// 		}
// 	}

// 	if f == 0 {
// 		t.Fatal("locking isnt working")
// 	} else {
// 		fmt.Printf("threads [%d] have been throttled\n", f)
// 	}
// }
