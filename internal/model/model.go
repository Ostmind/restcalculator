package model

import "sync"

type Calculation struct {
	Numbers   []float64 `json:"numbers"`
	Operation string    `json:"operation"`
	Result    float64   `json:"result"`
}

type Results struct {
	UserValues map[string][]Calculation
	Mutex      sync.Mutex
}

type Request struct {
	Value []float64 `json:"Value"`
}

func NewResult() *Results {
	return &Results{
		UserValues: make(map[string][]Calculation),
	}
}
