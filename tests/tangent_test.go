package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	tangent "github.com/dfsantos-source/tangent-backend/http"
)

var defaultParams = tangent.TangentRequestParams{
	End_Lat:     42.36025619506836,
	Term:        "restaurants",
	End_Lon:     -71.05728149414062,
	Open_Now:    false,
	Limit:       20,
	Price:       "2",
	Start_Lon:   -72.52603912353516,
	Start_Lat:   42.39096069335938,
	Pref_Radius: 24140,
}

func setupTest(tb testing.TB, params tangent.TangentRequestParams) (func(tb testing.TB), httptest.ResponseRecorder) {
	fmt.Println("setup test")

	req, err := http.NewRequest("GET", "/tangent", nil)

	if err != nil {
		tb.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("start_lat", fmt.Sprint(params.Start_Lat))
	q.Add("start_lon", fmt.Sprint(params.Start_Lon))
	q.Add("end_lat", fmt.Sprint(params.End_Lat))
	q.Add("end_lon", fmt.Sprint(params.End_Lon))
	q.Add("term", fmt.Sprint(params.Term))
	q.Add("open_now", fmt.Sprint(params.Open_Now))
	q.Add("price", fmt.Sprint(params.Price))
	q.Add("pref_radius", fmt.Sprint(params.Pref_Radius))
	q.Add("limit", fmt.Sprint(params.Limit))

	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tangent.CreateServer().GetTangent)
	handler.ServeHTTP(rr, req)

	return func(tb testing.TB) {
		fmt.Println("test teardown")
	}, *rr

}

func TestTangentHandler(t *testing.T) {

	TeardownTest, rr := setupTest(t, defaultParams)
	defer TeardownTest(t)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusOK, status)
	}

	var tangentRes tangent.TangentResponse

	if err := json.NewDecoder(rr.Body).Decode(&tangentRes); err != nil {
		t.Errorf("Error decoding tangent response")
	}
}
