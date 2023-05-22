package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	tangent "github.com/dfsantos-source/tangent-backend/http"
	models "github.com/dfsantos-source/tangent-backend/models"
	utils "github.com/dfsantos-source/tangent-backend/utils"
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

// Mock function for getYelpResponses to test without depending on 3rd party API
func mockGetYelpResponses(
	w http.ResponseWriter,
	r *http.Request,
	params *tangent.TangentRequestParams,
	coordinates [][]float32,
	token string,
) ([]models.Business, error) {

	businessSet := utils.BusinessSet{
		Businesses: make([]models.Business, 0),
		Set:        make(map[string]bool),
	}

	mockedYelpResponse := []models.Business{{
		Name:         "restuarant1",
		Id:           "newrest1",
		Rating:       4.2,
		Review_Count: 15,
		Coordinates:  models.Coordinates{Latitude: 42.3335, Longitude: 54.2332},
		Price:        "2",
		Categories:   []models.Categories{{Title: "outside"}, {Title: "modern"}},
		Image_Url:    "newimage.com",
	}, {
		Name:         "restuarant1",
		Id:           "newrest1",
		Rating:       4.2,
		Review_Count: 15,
		Coordinates:  models.Coordinates{Latitude: 42.3335, Longitude: 54.2332},
		Price:        "2",
		Categories:   []models.Categories{{Title: "outside"}, {Title: "modern"}},
		Image_Url:    "newimage.com",
	}, {
		Name:         "restuarant2",
		Id:           "newrest2",
		Rating:       4.3,
		Review_Count: 30,
		Coordinates:  models.Coordinates{Latitude: 44.3335, Longitude: 64.2332},
		Price:        "2",
		Categories:   []models.Categories{{Title: "outside"}, {Title: "modern"}},
		Image_Url:    "newimage.com/restaurant",
	}}

	businessSet.AddBusinesses(mockedYelpResponse)

	return *businessSet.GetBusinesses(), nil
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
		tb.Log("teardown test")
	}, *rr
}

func setupMockTest(tb testing.TB, params tangent.TangentRequestParams) (func(tb testing.TB), httptest.ResponseRecorder) {
	fmt.Println("setup Mock test")

	tangent.GetYelpResponses = mockGetYelpResponses

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
		tb.Log("teardown test")
	}, *rr

}

func TestTangentHandler(t *testing.T) {

	teardownTest, rr := setupTest(t, defaultParams)
	defer teardownTest(t)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusOK, status)
	}

	var tangentRes tangent.TangentResponse

	if err := json.NewDecoder(rr.Body).Decode(&tangentRes); err != nil {
		t.Errorf("Error decoding tangent response")
	}
}

func TestResponseDuplicates(t *testing.T) {

	teardownTest, rr := setupTest(t, defaultParams)
	defer teardownTest(t)

	var tangentRes tangent.TangentResponse

	if err := json.NewDecoder(rr.Body).Decode(&tangentRes); err != nil {
		t.Errorf("Error decoding tangent response")
	}

	businesses := tangentRes.Businesses
	dupMap := make(map[string]bool)

	for i := 0; i < len(businesses); i++ {
		if dupMap[businesses[i].Id] {
			t.Errorf("duplicate in response")
		} else {
			dupMap[businesses[i].Id] = true
		}
	}

}

func TestMockTangentHandler(t *testing.T) {
	TeardownTest, rr := setupMockTest(t, defaultParams)
	defer TeardownTest(t)

	var tangentRes tangent.TangentResponse

	if err := json.NewDecoder(rr.Body).Decode(&tangentRes); err != nil {
		t.Errorf("Error decoding tangent response")
	}
}

func TestMockResponseDuplicates(t *testing.T) {
	TeardownTest, rr := setupMockTest(t, defaultParams)
	defer TeardownTest(t)

	var tangentRes tangent.TangentResponse

	if err := json.NewDecoder(rr.Body).Decode(&tangentRes); err != nil {
		t.Errorf("Error decoding tangent response")
	}

	fmt.Println(tangentRes.Businesses)

	businesses := tangentRes.Businesses
	dupMap := make(map[string]bool)

	for i := 0; i < len(businesses); i++ {
		if dupMap[businesses[i].Id] {
			t.Errorf("duplicate in response")
		} else {
			dupMap[businesses[i].Id] = true
		}
	}
}
