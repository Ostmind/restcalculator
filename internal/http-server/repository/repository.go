package repository

import "restcalculator/internal/model"

func AppendRepository(res *model.Results, cookie string, calculation model.Calculation) error {
	res.Mutex.Lock()
	res.UserValues[cookie] = append(res.UserValues[cookie], calculation)
	res.Mutex.Unlock()
	return nil
}
