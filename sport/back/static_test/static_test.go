package static

import (
	"io/ioutil"
	"path"
	"sport/api"
	"testing"
)

func Initialization_Of_Environment(t *testing.T) {
	p := path.Join(staticCtx.Config.Basepath, "testdir", "cfe914e1-27d6-4af9-971a-0f64fb5ab024.jpg")
	if err := ioutil.WriteFile(p, []byte("123"), 0755); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestPathTraversal(t *testing.T) {
	staticCtx.RegisterDir("testdir")
	Initialization_Of_Environment(t)
	br, gr := GeenerateTestCases()
	staticCtx.Api.TestAssertCases(t, br)
	staticCtx.Api.TestAssertCases(t, gr)
}

func GeenerateTestCases() (badCases, goodCases []api.TestCase) {
	badCases = []api.TestCase{
		{
			// non existing file
			RequestMethod:      "GET",
			RequestUrl:         "/api/static/testdir/validFile123",
			ExpectedStatusCode: 400,
		},
		{
			// non existing file
			RequestMethod:      "GET",
			RequestUrl:         "/api/static/testdir/cfe914e1-27d6-4af9-971a-0f64fb5ab024.",
			ExpectedStatusCode: 400,
		},
		{
			// non existing file
			RequestMethod:      "GET",
			RequestUrl:         "/api/static/testdir/cfe914e1-27d6-4af9-971a-0f64fb5ab014.jpg",
			ExpectedStatusCode: 404,
		},
		{
			// non existing file
			RequestMethod:      "GET",
			RequestUrl:         "/api/static/testdir../../../../../../etc/shadow",
			ExpectedStatusCode: 404,
		},
		{
			// non existing file
			RequestMethod:      "GET",
			RequestUrl:         "/api/static/testdir%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2fetc%2fshadow",
			ExpectedStatusCode: 404,
		},
		{
			// non existing file
			RequestMethod:      "GET",
			RequestUrl:         "/api/static/testdir%252e%252e%252f%252e%252e%252f%252e%252e%252f%252e%252e%252f%252e%252e%252fetc%252fshadow",
			ExpectedStatusCode: 404,
		},
	}

	goodCases = []api.TestCase{
		{
			RequestMethod:      "GET",
			RequestUrl:         "/api/static/testdir/cfe914e1-27d6-4af9-971a-0f64fb5ab024.jpg",
			ExpectedStatusCode: 200,
		},
	}

	return badCases, goodCases
}
