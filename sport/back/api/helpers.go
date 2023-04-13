package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"sport/helpers"
	"strings"
	"testing"
)

func CreateMultipartForm(
	formFieldName, formFileName string, fc []byte) (bytes.Buffer, string, error) {

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile(formFieldName, formFileName)
	if err != nil {
		return b, "", err
	}
	if _, err := fw.Write(fc); err != nil {
		return b, "", err
	}
	w.Close()
	return b, w.FormDataContentType(), nil
}

type AdditionalFormField struct {
	Key, Val string
}

func CreateMultipartFormWithValues(
	formFieldName, formFileName string, fc []byte, afs []AdditionalFormField,
) (bytes.Buffer, string, error) {

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile(formFieldName, formFileName)
	if err != nil {
		return b, "", err
	}
	if _, err := fw.Write(fc); err != nil {
		return b, "", err
	}
	for i := range afs {
		w.WriteField(afs[i].Key, afs[i].Val)
	}
	w.Close()
	return b, w.FormDataContentType(), nil
}

/*
	used to automatically test api endpoinds.
*/
type TestCase struct {
	RequestMethod       string
	RequestUrl          string
	RequestReader       io.Reader
	RequestHeaders      map[string]string
	RequestQueryString  map[string]interface{}
	ExpectedStatusCode  int
	ExpectedBody        *string
	ExpectedBodyVal     func([]byte, interface{}) error
	ExpectedBodyValArgs interface{}
	ExpectedHeaders     map[string][]string
	ExpectedHeadersVal  func(map[string][]string) error
}

func (api *Ctx) TestAssertCaseErr(c *TestCase) error {
	if c.RequestReader == nil {
		c.RequestReader = strings.NewReader("")
	}
	var err error
	url := c.RequestUrl
	if len(c.RequestQueryString) > 0 {
		if !strings.HasPrefix(url, "?") {
			url += "?"
		}
		i := 0
		for k := range c.RequestQueryString {
			if i > 0 {
				url += "&"
			}
			i++
			url += fmt.Sprintf("%s=%v", k, c.RequestQueryString[k])
		}
	}
	req, err := http.NewRequest(c.RequestMethod, url, c.RequestReader)
	if err != nil {
		return err
	}
	if c.RequestHeaders != nil {
		for k := range c.RequestHeaders {
			req.Header.Add(k, c.RequestHeaders[k])
		}
	}
	recorder := httptest.NewRecorder()
	api.engine.ServeHTTP(recorder, req)
	if c.ExpectedStatusCode > 0 {
		if err = helpers.AssertErr(
			c.RequestUrl+": Invalid status code",
			recorder.Code,
			c.ExpectedStatusCode,
		); err != nil {
			return err
		}
	}
	if c.ExpectedBodyVal != nil {
		b := recorder.Body.Bytes()
		if err := c.ExpectedBodyVal(b, c.ExpectedBodyValArgs); err != nil {
			return err
		}
		if c.ExpectedBody != nil {
			if err = helpers.AssertErr(
				c.RequestUrl+": Invalid body",
				string(b),
				*c.ExpectedBody,
			); err != nil {
				return err
			}
		}
	} else if c.ExpectedBody != nil {
		if err = helpers.AssertErr(
			c.RequestUrl+": Invalid body",
			recorder.Body.String(),
			*c.ExpectedBody,
		); err != nil {
			return err
		}
	}
	if c.ExpectedHeaders != nil {
		var hdrs []string
		var ok bool
		for k := range c.ExpectedHeaders {
			if hdrs, ok = recorder.Header()[k]; !ok {
				return fmt.Errorf("Couldnt find header %s in response", k)
			}
			ehdrs := c.ExpectedHeaders[k]
			if len(hdrs) != len(ehdrs) {
				return fmt.Errorf("Invalid number of header values in response. Expected %d got %d", len(ehdrs), len(hdrs))
			}
			for i, ehdr := range ehdrs {
				if hdrs[i] != ehdr {
					return fmt.Errorf("Invalid header value in response. Expected %s got %s", hdrs[i], ehdr)
				}
			}
		}
	}
	if c.ExpectedHeadersVal != nil {
		if err := c.ExpectedHeadersVal(map[string][]string(recorder.Header())); err != nil {
			return err
		}
	}
	return nil
}

/*
	Test ApiTestCase array on global context.
	you must initialize Api global context before calling this function.
*/
func (api *Ctx) TestAssertCases(t *testing.T, tcs []TestCase) {
	var c *TestCase
	for i := 0; i < len(tcs); i++ {
		c = &tcs[i]
		err := api.TestAssertCaseErr(c)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (api *Ctx) TestAssertCasesErr(tcs []TestCase) error {
	var c *TestCase
	for i := 0; i < len(tcs); i++ {
		c = &tcs[i]
		err := api.TestAssertCaseErr(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (api *Ctx) TestAssertCasesParallel(t *testing.T, tcs []TestCase) {
	cqty := len(tcs)
	if cqty == 0 {
		return
	}
	errs := make([]error, cqty)
	waiter := make(chan interface{})
	for i := 0; i < cqty; i++ {
		go func(_i int) {
			c := &tcs[_i]
			errs[_i] = api.TestAssertCaseErr(c)
			waiter <- nil
		}(i)
	}

	for i := 0; i < cqty; i++ {
		<-waiter
	}

	for _, err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
}

type void struct{}

var null = void{}

func (api *Ctx) TestAssertCasesSemaphoreErr(tcs []TestCase, c int) error {
	cqty := len(tcs)
	if cqty == 0 {
		return nil // nothing to do
	}
	waiter := make(chan error, cqty)
	semaphore := make(chan void, c)
	for i := 0; i < cqty; i++ {
		semaphore <- null
		go func(_i int) {
			c := &tcs[_i]
			waiter <- api.TestAssertCaseErr(c)
			<-semaphore
		}(i)
	}
	for i := 0; i < cqty; i++ {
		err := <-waiter
		if err != nil {
			return err
		}
	}
	return nil
}

func (api *Ctx) TestAssertCasesSemaphore(t *testing.T, tcs []TestCase, c int) {
	if err := api.TestAssertCasesSemaphoreErr(tcs, c); err != nil {
		t.Fatal(err)
	}
}
