package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	ts     = httptest.NewServer(http.HandlerFunc(nodepingTestServer))
	secret = "123"
)

func nodepingTestServer(w http.ResponseWriter, r *http.Request) {
	user, _, ok := r.BasicAuth()
	if !ok || user != secret {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	switch r.URL.Path {
	case "/results/12345":
		fmt.Fprint(w, `[
			{
				"ci": "12345",
				"t": "SSL",
				"tg": "http://www.example.com/",
				"th": 5,
				"i": 1,
				"ra": "123456",
				"q": "1234562345",
				"s": 1559247850938,
				"sc": "valid",
				"su": true,
				"rt": 165,
				"rhp": "rh16",
				"m": "Valid cert",
				"e": 1559247851103,
				"l": {
					"1559247850938": "ro"
				},
				"f": "",
				"jobid": "1234-1234",
				"_id": "1234-1234-1234"
			}
		]`)
	case "/checks":
		fmt.Fprint(w, `{
			"201205050153W2Q4C-0J2HSIRF": {
				"_id": "201205050153W2Q4C-0J2HSIRF",
				"_rev": "37-8776f919267df3973fdb33cba0a8dd09",
				"customer_id": "201205050153W2Q4C",
				"label": "Site 1",
				"interval": 1,
				"notifications": [],
				"type": "HTTP",
				"status": "assigned",
				"modified": 1336759793520,
				"enable": "active",
				"public": false,
				"parameters": {
					"target": "http://www.example.com/",
					"threshold": 5,
					"sens": 2
				},
				"created": 1336185808566,
				"queue": "bINPckIRdv",
				"uuid": "4pybhg6m-4v1y-4enn-8tz5-tvywydu6h04k",
				"state": 0,
				"firstdown": 1336185868566
			}
		}`)
	default:
		fmt.Fprint(w, `{"result":null,"error":"BadRequest","id":1}`)
	}
}

type TestCase struct {
	APIURL        string
	Token         string
	expectedError bool
}

var TestCases = []TestCase{
	{ // everything is ok
		APIURL:        ts.URL,
		Token:         "123",
		expectedError: false,
	},
	{ // empty token
		APIURL:        ts.URL,
		Token:         "",
		expectedError: true,
	},
	{ // wrong token
		APIURL:        ts.URL,
		Token:         "321",
		expectedError: true,
	},
	{ // bad url
		APIURL:        "http://bad_url",
		Token:         "123",
		expectedError: true,
	},
	{ // not correct url
		APIURL:        "http://not correct url",
		Token:         "123",
		expectedError: true,
	},
	{ // bad json
		APIURL: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"error":"bad_json","id":1`)
		})).URL,
		Token:         "",
		expectedError: true,
	},
	{ // bad json
		APIURL: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})).URL,
		Token:         "",
		expectedError: true,
	},
}

func TestNewNodePing(t *testing.T) {
	for _, tc := range TestCases {
		_, err := NewNodePing(tc.APIURL, tc.Token)
		if err != nil && !tc.expectedError {
			t.Errorf("Error was not expected, got: %v", err)
		} else if tc.expectedError && err == nil {
			t.Error("Expected error, got nil")
		}
	}
}

func TestGetAllChecks(t *testing.T) {
	for _, tc := range TestCases {
		np := NodePing{
			APIURL: tc.APIURL,
			Token:  tc.Token,
		}

		checks, err := np.GetAllChecks()
		if err != nil && !tc.expectedError {
			t.Errorf("Error was not expected, got: %v", err)
			continue
		} else if tc.expectedError && err == nil {
			t.Error("Expected error, got nil")
			continue
		} else if err != nil && tc.expectedError {
			continue
		}
		if _, ok := checks["201205050153W2Q4C-0J2HSIRF"]; !ok {
			t.Errorf("Looks like response was not parsed correctly")
			continue
		}
	}
}

func TestGetCheckStats(t *testing.T) {
	for _, tc := range TestCases {
		np := NodePing{
			APIURL: tc.APIURL,
			Token:  tc.Token,
		}

		stats, err := np.GetCheckStats("12345")
		if err != nil && !tc.expectedError {
			t.Errorf("Error was not expected, got: %v", err)
			continue
		} else if tc.expectedError && err == nil {
			t.Error("Expected error, got nil")
			continue
		} else if err != nil && tc.expectedError {
			continue
		}
		if stats.Target != "http://www.example.com/" {
			t.Errorf("Looks like response was not parsed correctly")
			continue
		}
	}
}
