package main

type profile struct {
	id        string   `json:'id'`
	lat       float64  `json:'lat'`
	long      float64  `json:'long'`
	interests []string `json:'interests'`
}
