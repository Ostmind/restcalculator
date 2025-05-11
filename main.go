package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// будет общая структура для хранения запросов в памяти, так же будет использоваться под расширение на авторизацию пользователей!?
type Calculator struct {
	Numbers   []float64
	Operation string
	Result    float64
}

type Request struct {
	Value float64
}

var (
	req  []Request
	calc []Calculator
)

func handler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST":
		var sum float64
		var calculation Calculator
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewDecoder(r.Body).Decode(&req)

		for i := 0; i < len(req); i++ {
			//fmt.Printf("Got Number %.2f\n", req[i].Value) //Debug проверка как проходит обработка запроса
			sum += req[i].Value
			calculation.Numbers = append(calculation.Numbers, req[i].Value)
		}
		calculation.Operation = "+"
		calculation.Result = sum
		calc = append(calc, calculation)
		json.NewEncoder(w).Encode(calculation)
		/*//Проверка на то что у нас правда сохраняются операции в памяти
		for i := 0; i < len(calc); i++ {
			fmt.Println("Slice: ", calc[i].Numbers, "Operation: ", calc[i].Operation, "Result: ", calc[i].Result)
		}*/
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
