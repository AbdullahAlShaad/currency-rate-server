package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/shaad7/currency-rate-server/controller"
	"github.com/shaad7/currency-rate-server/util"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	var dbAddress = util.DefaultDatabaseAddress
	var dbPassword = os.Getenv("DB_PASSWORD")

	var testDBName = util.TestDatabaseName
	util.CreateMySqlDatabase(dbAddress, dbPassword, testDBName)

	controller, err := controller.NewController(dbAddress, dbPassword, testDBName)
	if err != nil {
		t.Errorf("Failed to Load Controller. Error : %s", err.Error())
	}
	data := getDummyData()
	err = controller.LoadData(data)
	if err != nil {
		t.Errorf("Error happened : %s\n", err.Error())
	}
	rateServer := NewServer("127.0.0.1:8080", controller)
	r := mux.NewRouter()
	r.HandleFunc("/rates/analyze", rateServer.analyzeRate)
	r.HandleFunc("/rates/latest", rateServer.getLatestRate)
	r.HandleFunc("/rates/{date}", rateServer.ServeURLs)
	ts := httptest.NewServer(r)
	defer ts.Close()

	testGetLatestRate(ts, t)
	testGetRates(ts, t)
	testRateAnalysis(ts, t)
	testUndefinedVerbs(ts, t)
	testUndefinedPaths(ts, t)
}

func testGetLatestRate(ts *httptest.Server, t *testing.T) {
	latestData := getDataFromAPIPath(ts, t, "/rates/latest")
	if latestData["base"] != "EUR" {
		t.Errorf("Base is not Euro")
	}
	if interfaceData, ok := latestData["rates"]; ok {
		data := interfaceData.(map[string]interface{})
		if val, ok := data["C"]; ok {
			floatVal := val.(float64)
			if !util.AlmostEqual(floatVal, 10) {
				t.Errorf("Lastest Value of C Does not match. ")
			}
		} else {
			t.Errorf("Latest value for C not found")
		}

	} else {
		t.Errorf("Rates not found")
	}
}

func testGetRates(ts *httptest.Server, t *testing.T) {
	latestData := getDataFromAPIPath(ts, t, "/rates/2023-02-04")
	if latestData["base"] != "EUR" {
		t.Errorf("Base is not Euro")
	}
	if interfaceData, ok := latestData["rates"]; ok {
		data := interfaceData.(map[string]interface{})
		if val, ok := data["B"]; ok {
			floatVal := val.(float64)
			if !util.AlmostEqual(floatVal, 2) {
				t.Errorf("Lastest Value of B Does not match. ")
			}
		} else {
			t.Errorf("Latest value for B not found")
		}
	} else {
		t.Errorf("Rates not found")
	}
}

func testRateAnalysis(ts *httptest.Server, t *testing.T) {
	latestData := getDataFromAPIPath(ts, t, "/rates/analyze")
	if latestData["base"] != "EUR" {
		t.Errorf("Base is not Euro")
	}
	if interfaceData, ok := latestData["rates_analyze"]; ok {
		data := interfaceData.(map[string]interface{})
		if val, ok := data["A"]; ok {
			info := val.(map[string]interface{})
			if val, ok := info["avg"]; ok {
				avg := val.(float64)
				if !util.AlmostEqual(avg, 9) {
					t.Errorf("Average expected 5 , found : %f", avg)
				}
			} else {
				t.Errorf("Can not find average")
			}
		} else {
			t.Errorf("Latest value for A not found")
		}
	} else {
		t.Errorf("Rates not found")
	}
}

func testUndefinedVerbs(ts *httptest.Server, t *testing.T) {
	resp, err := http.Post(ts.URL+"/rates/latest", "", nil)
	if err != nil {
		t.Errorf("Could not post")
	}
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Status Could Not Matched . Found : %d", resp.StatusCode)
	}

	resp, err = http.Post(ts.URL+"/rates/analyze", "", nil)
	if err != nil {
		t.Errorf("Could not post")
	}
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Status Could Not Matched . Found : %d", resp.StatusCode)
	}
}

func testUndefinedPaths(ts *httptest.Server, t *testing.T) {
	resp, err := http.Post(ts.URL+"/rates", "", nil)
	if err != nil {
		t.Errorf("Could not post")
	}
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Status Could Not Matched . Found : %d", resp.StatusCode)
	}

	resp, err = http.Post(ts.URL+"/names", "", nil)
	if err != nil {
		t.Errorf("Could not post")
	}
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Status Could Not Matched . Found : %d", resp.StatusCode)
	}
}

func getDataFromAPIPath(ts *httptest.Server, t *testing.T, path string) map[string]interface{} {
	resp, err := http.Get(ts.URL + path)
	if err != nil {
		t.Errorf("Could not get data from %s path", path)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status Could Not Ok for path %s . Found : %d", path, resp.StatusCode)
	}
	var data map[string]interface{}
	jsonDataFromHttp, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Cannot read response body")
	}

	err = json.Unmarshal(jsonDataFromHttp, &data)
	if err != nil {
		t.Errorf("Can not get  data. Error : %s", err.Error())
	}
	return data
}

func getDummyData() *controller.RateSliceStruct {
	var data = controller.RateSlice{
		{Currency: "A",
			Date: "2023-02-05",
			Rate: "5",
		},
		{Currency: "A",
			Date: "2023-02-04",
			Rate: "10",
		},
		{Currency: "A",
			Date: "2023-02-03",
			Rate: "12",
		},

		{Currency: "B",
			Date: "2023-02-05",
			Rate: "1",
		},
		{Currency: "B",
			Date: "2023-02-04",
			Rate: "2",
		},
		{Currency: "B",
			Date: "2023-02-03",
			Rate: "3",
		},

		{Currency: "C",
			Date: "2023-02-05",
			Rate: "10",
		},
		{Currency: "C",
			Date: "2023-02-04",
			Rate: "20",
		},
		{Currency: "C",
			Date: "2023-02-03",
			Rate: "30",
		},
	}

	return &controller.RateSliceStruct{
		Rates: data,
	}
}
