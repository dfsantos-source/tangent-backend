package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin/render"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

type TangentRequestParams struct {
	Start_Lat float32 `json:"start_lat"`
	Start_Lon float32 `json:"start_lon"`
	End_Lat   float32 `json:"end_lat"`
	End_Lon   float32 `json:"end_lon"`

	Pref_Radius float32 `json:"pref_radius"`
	Term        string  `json:"term"`
	// Price       []int   `json:"price"`
	Open_Now bool `json:"open_now"`
	Limit    int  `json:"limit"`
}

type Coordinates struct {
	Latitude  string
	Longitude string
}

type Route struct {
	Geometry struct {
		Coordinates [][]float32 `json:"coordinates"`
	} `json:"geometry"`
}

type Routes struct {
	Routes []Route `json:"routes"`
}

type Business struct {
	Id           string
	Name         string
	Rating       float32
	Review_count int
	Coordinates  Coordinates
	Price        string
	Term         string
}

type Businesses struct {
	Businesses []Business
}

type TangentResponse struct {
	Businesses  []Business  `json:"businesses"`
	Coordinates [][]float32 `json:"coordinates"`
}

var decoder = schema.NewDecoder()

func (s *Server) registerTangentRoutes(r *chi.Mux) {
	r.Get("/tangents", s.getTangent)
	r.Get("/test", s.test)
}

func (s *Server) test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test route evoked."))
}

func getMapboxResponse(w http.ResponseWriter, r *http.Request, params *TangentRequestParams, token string) (*Routes, error) {
	// delims (given by mapbox)
	delim := "%2C"
	delim2 := "%3B"

	// concatenate start and end location to url query
	url := fmt.Sprintf(`https://api.mapbox.com/directions/v5/mapbox/driving/%s%s%s%s%s%s%s?alternatives=true&geometries=geojson&language=en&overview=simplified&steps=true&access_token=%s`,
		fmt.Sprint(params.Start_Lon), delim, fmt.Sprint(params.Start_Lat), delim2, fmt.Sprint(params.End_Lon), delim, fmt.Sprint(params.End_Lat), token)
	response, err := http.Get(url)
	if err != nil {
		render.WriteJSON(w, err)
		return nil, err
	}
	mapBody, _ := ioutil.ReadAll(response.Body)

	var routes Routes
	parseErr := json.Unmarshal(mapBody, &routes)
	if err != nil {
		render.WriteJSON(w, parseErr)
		return nil, parseErr
	}

	return &routes, nil
}

func getYelpResponse(w http.ResponseWriter, r *http.Request, params *TangentRequestParams, coordinates *Coordinates, token string) ([]Business, error) {
	url := fmt.Sprintf(`https://api.yelp.com/v3/businesses/search?latitude=%s&longitude=%s&term=food&radius=24140&sort_by=best_match&limit=5`, coordinates.Latitude, coordinates.Longitude)

	// build request URL with token
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	// make API call to Yelp
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		render.WriteJSON(w, err)
		return nil, err
	}

	mapBody, _ := ioutil.ReadAll(res.Body)

	var businesses Businesses
	parseErr := json.Unmarshal(mapBody, &businesses)
	if err != nil {
		render.WriteJSON(w, parseErr)
		return nil, parseErr
	}

	return businesses.Businesses, nil
}

func runYelp(w http.ResponseWriter, r *http.Request, params *TangentRequestParams, coordinates [][]float32, token string) []Business {
	tangentResponse := TangentResponse{}
	aggregateList := tangentResponse.Businesses
	size := len(coordinates)
	size = size - (size % 5)
	fmt.Println(size)
	for i := 0; i <= size; i += 5 {
		coordinate := coordinates[i]
		fmt.Println(coordinate[0])
		fmt.Println(coordinate[1])
		fmt.Println("=======")

		businesses, err := getYelpResponse(w, r, params, &Coordinates{Latitude: fmt.Sprint(coordinate[1]), Longitude: fmt.Sprint(coordinate[0])}, token)
		if err != nil {
			render.WriteJSON(w, err)
		}
		for i := 0; i < len(businesses); i++ {
			aggregateList = append(aggregateList, businesses[i])
			fmt.Println(aggregateList[i].Name)
		}
	}
	fmt.Println("==================")
	for i := 0; i < len(aggregateList); i++ {
		fmt.Println(aggregateList[i].Name)
	}
	return aggregateList
}

func (s *Server) getTangent(w http.ResponseWriter, r *http.Request) {

	mapboxToken := s.MapboxUtil.GetToken()

	var params = *new(TangentRequestParams)
	err := decoder.Decode(&params, r.URL.Query())

	if err != nil {
		render.WriteJSON(w, err)
		return
	}

	mapboxResponse, mapboxErr := getMapboxResponse(w, r, &params, mapboxToken)
	if mapboxErr != nil {
		render.WriteJSON(w, mapboxErr)
		return
	}

	coordinates := mapboxResponse.Routes[0].Geometry.Coordinates
	yelpToken := s.YelpUtil.GetToken()

	tangentResponse := TangentResponse{}
	businesses := runYelp(w, r, &params, coordinates, yelpToken)
	tangentResponse.Businesses = businesses
	tangentResponse.Coordinates = coordinates

	render.WriteJSON(w, tangentResponse)
}
