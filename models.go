package main

type profile struct {
	id        string   `json:'id'`
	lat       float64  `json:'lat'`
	long      float64  `json:'long'`
	interests []string `json:'interests'`
}

type location struct {
	id   string  `json:'id'`
	name string  `json:'name'`
	city string  `json:'city'`
	lat  float64 `json:'lat'`
	long float64 `json:'long'`
}
