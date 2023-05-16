package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	models "github.com/dfsantos-source/tangent-backend/models"
	utils "github.com/dfsantos-source/tangent-backend/utils"

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
	Price       string  `json:"price"`
	Open_Now    bool    `json:"open_now"`
	Limit       int     `json:"limit"`
}

type TangentResponse struct {
	Businesses  []models.Business `json:"businesses"`
	Coordinates [][]float32       `json:"coordinates"`
}

var (
	decoder             = schema.NewDecoder()
	COORDINATE_INTERVAL = 5
)

func (s *Server) registerTangentRoutes(r *chi.Mux) {
	r.Get("/tangents", s.getTangent)
	r.Get("/test", s.test)
}

func (s *Server) test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test route evoked."))
}

func getMapboxResponse(
	w http.ResponseWriter,
	r *http.Request,
	params *TangentRequestParams,
	token string,
) (*models.MapboxResponse, error) {
	delim := "%2C"
	delim2 := "%3B"
	url := fmt.Sprintf(`https://api.mapbox.com/directions/v5/mapbox/driving/%s%s%s%s%s%s%s?alternatives=true&geometries=geojson&language=en&overview=simplified&steps=true&access_token=%s`,
		fmt.Sprint(params.Start_Lon),
		delim,
		fmt.Sprint(params.Start_Lat),
		delim2,
		fmt.Sprint(params.End_Lon),
		delim,
		fmt.Sprint(params.End_Lat),
		token,
	)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf(string(body))
	}

	var mapboxResponse models.MapboxResponse
	err = json.Unmarshal(body, &mapboxResponse)
	if err != nil {
		return nil, err
	}

	return &mapboxResponse, nil
}

func getYelpResponse(
	w http.ResponseWriter,
	r *http.Request,
	params *TangentRequestParams,
	coordinates *models.Coordinates,
	token string,
) ([]models.Business, error) {
	priceQuery := utils.ParsePrice(params.Price)
	url := fmt.Sprintf(`https://api.yelp.com/v3/businesses/search?latitude=%s&longitude=%s&term=%s&radius=%s&open_now=%s&sort_by=best_match&limit=5%s`,
		fmt.Sprint(coordinates.Latitude),
		fmt.Sprint(coordinates.Longitude),
		params.Term,
		fmt.Sprint(params.Pref_Radius),
		fmt.Sprint(params.Open_Now),
		priceQuery,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(string(body))
	}

	var yelpResponse models.YelpResponse
	err = json.Unmarshal(body, &yelpResponse)
	if err != nil {
		return nil, err
	}

	return yelpResponse.Businesses, nil
}

func getYelpResponses(
	w http.ResponseWriter,
	r *http.Request,
	params *TangentRequestParams,
	coordinates [][]float32,
	token string,
) ([]models.Business, error) {
	tangentResponse := TangentResponse{}
	aggregateList := tangentResponse.Businesses

	size := len(coordinates)
	size = size - (size % COORDINATE_INTERVAL)
	fmt.Println(size)

	channel := make(chan []models.Business, size/COORDINATE_INTERVAL)

	for i := 0; i <= size; i += COORDINATE_INTERVAL {
		if i < len(coordinates) {
			coordinate := coordinates[i]
			go func(coordinate []float32) {
				businesses, err := getYelpResponse(w, r, params, &models.Coordinates{Latitude: coordinate[1], Longitude: coordinate[0]}, token)
				if err != nil {
				}
				channel <- businesses
			}(coordinate)
		}
	}

	for i := 0; i < size/COORDINATE_INTERVAL; i++ {
		businesses := <-channel
		aggregateList = append(aggregateList, businesses...)
	}

	fmt.Println(aggregateList)

	return aggregateList, nil
}

func (s *Server) getTangent(w http.ResponseWriter, r *http.Request) {
	mapboxToken := s.MapboxUtil.GetToken()
	yelpToken := s.YelpUtil.GetToken()
	tangentResponse := TangentResponse{}

	var params = *new(TangentRequestParams)
	err := decoder.Decode(&params, r.URL.Query())
	if err != nil {
		w.WriteHeader(400)
		render.WriteJSON(w, err.Error())
		return
	}

	res, err := getMapboxResponse(w, r, &params, mapboxToken)
	if err != nil {
		w.WriteHeader(400)
		render.WriteJSON(w, err.Error())
		return
	}

	coordinates := res.Routes[0].Geometry.Coordinates
	businesses, err := getYelpResponses(w, r, &params, coordinates, yelpToken)
	if err != nil {
		w.WriteHeader(400)
		render.WriteJSON(w, err.Error())
		return
	}

	tangentResponse.Businesses = businesses
	tangentResponse.Coordinates = coordinates

	render.WriteJSON(w, tangentResponse)
}
