package model

type Calculator struct {
	Numbers   []float64 `json:"numbers"`
	Operation string    `json:"operation"`
	Result    float64   `json:"result"`
}

type Results struct {
	Calc []Calculator
}

type Request struct {
	Value float64 `json:"Value"`
}
