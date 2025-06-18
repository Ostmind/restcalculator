package controller

import (
	"errors"
	"restcalculator/internal/http-server/repository"
	"restcalculator/internal/model"
)

const (
	SumOperation  = "+"
	MultOperation = "*"
)

func MakeCalc(num []float64, operation string, res *model.Results, cookie string) (calculation model.Calculation, err error) {

	switch operation {
	case SumOperation:
		var sum float64

		for i := range num {
			sum += num[i]

			calculation.Numbers = append(calculation.Numbers, num[i])
		}

		calculation.Operation = operation
		calculation.Result = sum

		err = repository.AppendRepository(res, cookie, calculation)

	case MultOperation:
		var mult float64 = 1

		for i := range num {
			mult *= num[i]

			calculation.Numbers = append(calculation.Numbers, num[i])
		}

		calculation.Operation = operation
		calculation.Result = mult

		err = repository.AppendRepository(res, cookie, calculation)

	default:
		err = errors.New("No Operation")
	}

	return calculation, err
}
