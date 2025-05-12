package model

type Calculator struct {
	Numbers   []float64
	Operation string
	Result    float64
}

type Results struct {
	Calc []Calculator
}

type Request struct {
	Value float64 `json:"Value"`
}
