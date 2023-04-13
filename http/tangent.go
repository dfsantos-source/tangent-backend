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

type Location struct {
	Latitude  float32
	Longitude float32
}

var decoder = schema.NewDecoder()

func (s *Server) registerTangentRoutes(r *chi.Mux) {
	r.Get("/tangents", s.getTangent)
	r.Get("/test", s.test)
}

func (s *Server) test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test route evoked."))
}

func getMapboxResponse(w http.ResponseWriter, r *http.Request, params *TangentRequestParams, token string) ([]byte, error) {
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

	return mapBody, nil
}

func getYelpResponse(w http.ResponseWriter, r *http.Request, params *TangentRequestParams, location *Location, token string) {
	url := fmt.Sprintf(`https://api.yelp.com/v3/businesses/search?latitude=42&longitude=-71&term=food&radius=24140&sort_by=best_match&limit=20`)

	// build request URL with token
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	// make API call to Yelp
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		render.WriteJSON(w, err)
		return
	}

	// process body
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	var yelpResponse map[string]interface{}

	// parse response to json
	err = json.Unmarshal(body, &yelpResponse)
	if err != nil {
		render.WriteJSON(w, err)
		return
	}

	fmt.Println(string(fmt.Sprint(yelpResponse["businesses"])))
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

	fmt.Println(string(mapboxResponse))

	yelpToken := s.YelpUtil.GetToken()

	location := &Location{Latitude: 42.466415, Longitude: -72.555244}

	getYelpResponse(w, r, &params, location, yelpToken)
}
