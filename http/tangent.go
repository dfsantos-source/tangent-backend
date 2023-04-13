package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin/render"
	"github.com/go-chi/chi/v5"
)

func (s *Server) registerTangentRoutes(r *chi.Mux) {
	r.Get("/tangents", s.getTangent)
	r.Get("/test", s.test)
}

func (s *Server) test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test route evoked."))
}

func (s *Server) getTangent(w http.ResponseWriter, r *http.Request) {
	yelpToken := s.YelpUtil.GetToken()

	url := `https://api.yelp.com/v3/businesses/search?
	latitude=42&
	longitude=-71&
	term=food&radius=24140&
	sort_by=best_match&
	limit=20
	`

	// build request URL with token
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+yelpToken)

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

	render.WriteJSON(w, yelpResponse["businesses"])
}
