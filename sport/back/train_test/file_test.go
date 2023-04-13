package train_test

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"sport/api"
	"sport/helpers"
	"sport/train"
	"sport/user"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestTrainingFile(t *testing.T) {
	token1, err := trainCtx.User.ApiCreateAndLoginUser(nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := trainCtx.Instr.ApiCreateInstructor(token1, nil); err != nil {
		t.Fatal(err)
	}
	token2, err := trainCtx.User.ApiCreateAndLoginUser(nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := trainCtx.Instr.ApiCreateInstructor(token2, nil); err != nil {
		t.Fatal(err)
	}
	trainingID, err := trainCtx.ApiCreateTraining(token1, &train.CreateTrainingRequest{
		Training: train.TrainingRequest{Title: "training1", Price: 1000, Capacity: 1, Currency: "PLN"},
	})
	if err != nil {
		t.Fatal(err)
	}

	method := "POST"
	p := "/api/training/img"

	h1 := user.GetAuthHeader(token1)

	ir, ct, err := api.CreateMultipartForm("image", "2.txt", []byte("1"))
	if err != nil {
		t.Fatal(err)
	}
	h1["Content-Type"] = ct

	tr, ct, err := api.CreateMultipartFormWithValues(
		"image", "2.txt", []byte("1"),
		[]api.AdditionalFormField{
			{
				Key: "training_id",
				Val: trainingID.String(),
			},
		})
	h2 := user.GetAuthHeader(token1)
	h2["Content-Type"] = ct
	if err != nil {
		t.Fatal(err)
	}

	fc, err := ioutil.ReadFile("../rsv_test/test_files/1.jpg")
	if err != nil {
		t.Fatal(err)
	}

	cr, ct, err := api.CreateMultipartFormWithValues(
		"image", "1.jpg", fc,
		[]api.AdditionalFormField{
			{
				Key: "training_id",
				Val: trainingID.String(),
			},
			{
				Key: "main",
				Val: "1",
			},
		})
	h3 := user.GetAuthHeader(token1)
	h3["Content-Type"] = ct
	if err != nil {
		t.Fatal(err)
	}

	tc := []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         p,
			ExpectedStatusCode: 401,
		},
		{
			RequestMethod:      method,
			RequestUrl:         p,
			RequestHeaders:     h1,
			ExpectedStatusCode: 400,
		},
		{
			RequestMethod:      method,
			RequestUrl:         p,
			RequestHeaders:     h1,
			RequestReader:      &ir,
			ExpectedStatusCode: 404,
		},
		// invalid image
		{
			RequestMethod:      method,
			RequestUrl:         p,
			RequestHeaders:     h2,
			RequestReader:      &tr,
			ExpectedStatusCode: 400,
		},
		// valid img
		{
			RequestMethod:      method,
			RequestUrl:         p,
			RequestHeaders:     h3,
			RequestReader:      &cr,
			ExpectedStatusCode: 204,
		},
	}
	if err := trainCtx.Api.TestAssertCasesErr(tc); err != nil {
		t.Fatal(err)
	}

	training, err := trainCtx.DalReadSingleTraining(train.DalReadTrainingsRequest{
		TrainingID: &trainingID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(training.Training.SecondaryImgUrls) != 0 {
		t.Fatal("invalid secondary img urls")
	}
	if training.Training.MainImgUrl == "" {
		t.Fatal("invalid main img url")
	}
	if _, err := url.Parse(training.Training.MainImgUrl); err != nil {
		t.Fatal(err)
	}
	ix := strings.LastIndex(training.Training.MainImgUrl, "/")
	if ix < 0 {
		t.Fatal("malformed url. that really shouldnt happen")
	}

	ppath := path.Join(
		trainCtx.Static.Config.Basepath,
		training.Training.MainImgID,
	)

	// file must be present and accesible
	if _, err := os.Stat(
		ppath,
	); err != nil {
		t.Fatal(err)
	}

	cr, ct, err = api.CreateMultipartFormWithValues(
		"image", "1.jpg", fc,
		[]api.AdditionalFormField{
			{
				Key: "training_id",
				Val: trainingID.String(),
			},
			{
				Key: "main",
				Val: "1",
			},
		})
	h3 = user.GetAuthHeader(token1)
	h3["Content-Type"] = ct
	if err != nil {
		t.Fatal(err)
	}

	tc = []api.TestCase{
		{
			RequestMethod:      "GET",
			RequestUrl:         training.Training.MainImgUrl,
			ExpectedStatusCode: 200,
		},
		{
			RequestMethod:      method,
			RequestUrl:         p,
			RequestHeaders:     h3,
			RequestReader:      &cr,
			ExpectedStatusCode: 204,
		},
	}
	if err := trainCtx.Api.TestAssertCasesErr(tc); err != nil {
		t.Fatal(err)
	}

	training, err = trainCtx.DalReadSingleTraining(train.DalReadTrainingsRequest{
		TrainingID: &trainingID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(training.Training.SecondaryImgUrls) != 0 {
		t.Fatal("invalid secondary img urls")
	}
	if training.Training.MainImgUrl == "" {
		t.Fatal("invalid main img url")
	}
	if _, err := url.Parse(training.Training.MainImgUrl); err != nil {
		t.Fatal(err)
	}
	ix = strings.LastIndex(training.Training.MainImgUrl, "/")
	if ix < 0 {
		t.Fatal("malformed url. that really shouldnt happen")
	}

	// file must be present and accesible
	if _, err := os.Stat(
		path.Join(
			trainCtx.Static.Config.Basepath,
			training.Training.MainImgID,
		),
	); err != nil {
		t.Fatal(err)
	}

	// previous file must be gone
	if _, err := os.Stat(ppath); err == nil {
		t.Fatal(fmt.Errorf("prev. file should be gone by now"))
	} else {
		if !os.IsNotExist(err) {
			t.Fatal(fmt.Errorf("got err: %v, expected not found", err))
		}
	}

	// access control

	cr, ct, err = api.CreateMultipartFormWithValues(
		"image", "1.jpg", fc,
		[]api.AdditionalFormField{
			{
				Key: "training_id",
				Val: trainingID.String(),
			},
			{
				Key: "main",
				Val: "1",
			},
		})
	h3 = user.GetAuthHeader(token2)
	h3["Content-Type"] = ct
	if err != nil {
		t.Fatal(err)
	}

	tc = []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         p,
			RequestHeaders:     h3,
			RequestReader:      &cr,
			ExpectedStatusCode: 404,
		},
	}
	if err := trainCtx.Api.TestAssertCasesErr(tc); err != nil {
		t.Fatal(err)
	}

	//

	// add secondary img

	cr, ct, err = api.CreateMultipartFormWithValues(
		"image", "1.jpg", fc,
		[]api.AdditionalFormField{
			{
				Key: "training_id",
				Val: trainingID.String(),
			},
		})
	h3 = user.GetAuthHeader(token1)
	h3["Content-Type"] = ct
	if err != nil {
		t.Fatal(err)
	}

	tc = []api.TestCase{
		{
			RequestMethod:      method,
			RequestUrl:         p,
			RequestHeaders:     h3,
			RequestReader:      &cr,
			ExpectedStatusCode: 204,
		},
	}
	if err := trainCtx.Api.TestAssertCasesErr(tc); err != nil {
		t.Fatal(err)
	}
	training, err = trainCtx.DalReadSingleTraining(train.DalReadTrainingsRequest{
		TrainingID: &trainingID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(training.Training.SecondaryImgUrls) != 1 {
		t.Fatal("invalid secondary img urls")
	}
	if training.Training.MainImgUrl == "" {
		t.Fatal("invalid main img url")
	}

	mainPath := path.Join(
		trainCtx.Static.Config.Basepath,
		training.Training.MainImgID,
	)

	method = "DELETE"
	tc = []api.TestCase{
		{
			RequestMethod:  method,
			RequestUrl:     p,
			RequestHeaders: user.GetAuthHeader(token1),
			RequestReader: helpers.JsonMustSerializeReader(train.DeleteTrainingImageRequest{
				ID:         uuid.New().String(),
				TrainingID: trainingID,
			}),
			ExpectedStatusCode: 404,
		},
		{
			RequestMethod:  method,
			RequestUrl:     p,
			RequestHeaders: user.GetAuthHeader(token1),
			RequestReader: helpers.JsonMustSerializeReader(train.DeleteTrainingImageRequest{
				ID:         training.Training.MainImgID,
				TrainingID: uuid.New(),
			}),
			ExpectedStatusCode: 404,
		},
		{
			RequestMethod:  method,
			RequestUrl:     p,
			RequestHeaders: user.GetAuthHeader(token2),
			RequestReader: helpers.JsonMustSerializeReader(train.DeleteTrainingImageRequest{
				ID:         training.Training.MainImgID,
				TrainingID: training.Training.ID,
			}),
			ExpectedStatusCode: 404,
		},
		// delete main
		{
			RequestMethod:  method,
			RequestUrl:     p,
			RequestHeaders: user.GetAuthHeader(token1),
			RequestReader: helpers.JsonMustSerializeReader(train.DeleteTrainingImageRequest{
				ID:         training.Training.MainImgID,
				TrainingID: training.Training.ID,
			}),
			ExpectedStatusCode: 204,
		},
	}
	if err := trainCtx.Api.TestAssertCasesErr(tc); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(mainPath); err == nil {
		t.Fatal(fmt.Errorf("image should be gone by now"))
	} else {
		if !os.IsNotExist(err) {
			t.Fatal(fmt.Errorf("image should be gone by now, err: %v", err))
		}
	}

	training, err = trainCtx.DalReadSingleTraining(train.DalReadTrainingsRequest{
		TrainingID: &trainingID,
	})

	if training.Training.MainImgID != "" {
		t.Fatal(fmt.Errorf("invalid main img id"))
	}

	secpath := path.Join(
		trainCtx.Static.Config.Basepath,
		training.Training.SecondaryImgIDs[0],
	)

	tc = []api.TestCase{
		{
			RequestMethod:  method,
			RequestUrl:     p,
			RequestHeaders: user.GetAuthHeader(token1),
			RequestReader: helpers.JsonMustSerializeReader(train.DeleteTrainingImageRequest{
				ID:         training.Training.SecondaryImgIDs[0],
				TrainingID: training.Training.ID,
			}),
			ExpectedStatusCode: 204,
		},
	}
	if err := trainCtx.Api.TestAssertCasesErr(tc); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(secpath); err == nil {
		t.Fatal(fmt.Errorf("image should be gone by now"))
	} else {
		if !os.IsNotExist(err) {
			t.Fatal(fmt.Errorf("image should be gone by now, err: %v", err))
		}
	}

	training, err = trainCtx.DalReadSingleTraining(train.DalReadTrainingsRequest{
		TrainingID: &trainingID,
	})

	if len(training.Training.SecondaryImgIDs) != 0 {
		t.Fatal(fmt.Errorf("invalid sec img ids"))
	}
}
