package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"restcalculator/internal/http-server/handlers"
	"restcalculator/internal/model"
	"sync"
	"testing"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	// Просто игнорируем запись журнала
	return nil
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	// Возвращает тот же обработчик, так как нет атрибутов для сохранения
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	// Возвращает тот же обработчик, так как нет группы для сохранения
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// Всегда возвращает false, так как запись журнала игнорируется
	return false
}

func newEchoContext(t *testing.T, method, path string, body interface{}, cookie *http.Cookie) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()

	var reqBody bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&reqBody).Encode(body)
		if err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest(method, path, &reqBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if cookie != nil {
		req.AddCookie(cookie)
	}

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	return c, rec
}

func TestController_Calculation_Sum(t *testing.T) {
	logger := slog.New(&DiscardHandler{})
	res := &model.Results{
		UserValues: make(map[string][]model.Calculation),
		Mutex:      sync.Mutex{},
	}

	ctrl := handlers.New(logger, res)

	reqBody := model.Request{Value: []float64{1, 2, 3}}
	cookie := &http.Cookie{Name: "Token", Value: "user1"}

	c, rec := newEchoContext(t, http.MethodPost, "/calc/+", reqBody, cookie)
	c.SetParamNames("operation")
	c.SetParamValues("+")

	if err := ctrl.Calculation(c); err != nil {
		t.Fatalf("Calculation() returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status %d but got %d", http.StatusOK, rec.Code)
	}

	var calc model.Calculation
	if err := json.NewDecoder(rec.Body).Decode(&calc); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expectedResult := 6.0
	if calc.Result != expectedResult {
		t.Errorf("Expected result %v, got %v", expectedResult, calc.Result)
	}

	if calc.Operation != "+" {
		t.Errorf("Expected operation '+', got %v", calc.Operation)
	}

	if len(calc.Numbers) != 3 {
		t.Errorf("Expected 3 numbers, got %d", len(calc.Numbers))
	}

	res.Mutex.Lock()
	defer res.Mutex.Unlock()
	calcs, ok := res.UserValues["user1"]
	if !ok || len(calcs) == 0 {
		t.Error("Expected stored calculation for user1")
	}
}

func TestController_Calculation_Mult(t *testing.T) {
	logger := slog.New(&DiscardHandler{})
	res := &model.Results{
		UserValues: make(map[string][]model.Calculation),
		Mutex:      sync.Mutex{},
	}

	ctrl := handlers.New(logger, res)

	reqBody := model.Request{Value: []float64{2, 3, 4}}
	cookie := &http.Cookie{Name: "Token", Value: "user2"}

	c, rec := newEchoContext(t, http.MethodPost, "/calc/*", reqBody, cookie)
	c.SetParamNames("operation")
	c.SetParamValues("*")

	if err := ctrl.Calculation(c); err != nil {
		t.Fatalf("Calculation() returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status %d but got %d", http.StatusOK, rec.Code)
	}

	var calc model.Calculation
	if err := json.NewDecoder(rec.Body).Decode(&calc); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expectedResult := 1.0 * 2 * 3 * 4

	if calc.Result != expectedResult {
		t.Errorf("Expected result %v, got %v", expectedResult, calc.Result)
	}

	if calc.Operation != "*" {
		t.Errorf("Expected operation '*', got %v", calc.Operation)
	}

	if len(calc.Numbers) != 3 {
		t.Errorf("Expected 3 numbers, got %d", len(calc.Numbers))
	}

	res.Mutex.Lock()
	defer res.Mutex.Unlock()
	calcs, ok := res.UserValues["user2"]
	if !ok || len(calcs) == 0 {
		t.Error("Expected stored calculation for user2")
	}
}

func TestController_Calculation_NoCookie(t *testing.T) {
	logger := slog.New(&DiscardHandler{})
	res := &model.Results{
		UserValues: make(map[string][]model.Calculation),
		Mutex:      sync.Mutex{},
	}

	ctrl := handlers.New(logger, res)

	reqBody := model.Request{Value: []float64{1, 2, 3}}

	c, rec := newEchoContext(t, http.MethodPost, "/calc/+", reqBody, nil)
	c.SetParamNames("operation")
	c.SetParamValues("+")

	err := ctrl.Calculation(c)
	if err == nil {
		t.Fatal("Expected error due to missing cookie, got nil")
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d but got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestController_Calculation_BadJSON(t *testing.T) {
	logger := slog.New(&DiscardHandler{})
	res := &model.Results{
		UserValues: make(map[string][]model.Calculation),
		Mutex:      sync.Mutex{},
	}

	ctrl := handlers.New(logger, res)

	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/calc/+", bytes.NewBufferString("bad json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(&http.Cookie{Name: "Token", Value: "user3"})
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("operation")
	c.SetParamValues("+")

	err := ctrl.Calculation(c)
	if err == nil {
		t.Fatal("Expected error due to bad JSON, got nil")
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d but got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestController_Calculation_UnknownOperation(t *testing.T) {
	logger := slog.New(&DiscardHandler{})
	res := &model.Results{
		UserValues: make(map[string][]model.Calculation),
		Mutex:      sync.Mutex{},
	}

	ctrl := handlers.New(logger, res)

	reqBody := model.Request{Value: []float64{1, 2, 3}}
	cookie := &http.Cookie{Name: "Token", Value: "user4"}

	c, rec := newEchoContext(t, http.MethodPost, "/calc/-", reqBody, cookie)
	c.SetParamNames("operation")
	c.SetParamValues("-")

	err := ctrl.Calculation(c)
	if err == nil {
		t.Fatal("Expected error due to unknown operation, got nil")
	}

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status %d but got %d", http.StatusInternalServerError, rec.Code)
	}
}
