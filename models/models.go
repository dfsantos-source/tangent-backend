package models

type Coordinates struct {
	Latitude  float32
	Longitude float32
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
