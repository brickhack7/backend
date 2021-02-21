package main

type Profile struct {
	ID        string   `json:'id'`
	Lat       float64  `json:'lat'`
	Long      float64  `json:'long'`
	Interests []string `json:'interests'`
}

type Location struct {
	ID       string  `json:'id'`
	Name     string  `json:'name'`
	City     string  `json:'city'`
	Distance float64 `json:'lat'`
}
