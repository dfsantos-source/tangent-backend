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

type Geometry struct {
	Coordinates [][]float32 `json:"coordinates"`
}

type MapboxResponse struct {
	Routes []Route `json:"routes"`
}

type Categories struct {
	Title string
}

type Business struct {
	Id           string
	Name         string
	Rating       float32
	Review_Count int
	Coordinates  Coordinates
	Price        string
	Categories   []Categories
	Image_Url    string
}

type YelpResponse struct {
	Businesses []Business
}
