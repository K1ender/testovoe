package handlers_test

import (
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testovoe/internal/handlers"
	"testovoe/internal/models"
	"testovoe/internal/storage"
	mock_storage "testovoe/internal/storage/mocks"
	"testovoe/internal/utils"
	"testovoe/internal/validators"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

// Some attempts to write tests (they're shitty)
// TODO: improve these tests

func setupTest(t *testing.T) (*handlers.SubscriptionHandler, *mock_storage.MockSubscriptionStorage) {
	t.Helper()
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, nil)))

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("mm_yyyy", validators.MonthYearValidator)

	mockedSubscriptionStorage := mock_storage.NewMockSubscriptionStorage(ctrl)

	handler := handlers.NewSubscriptionHandler(
		&storage.Storage{
			Subscription: mockedSubscriptionStorage,
		},
		validate,
	)

	return handler, mockedSubscriptionStorage
}

func TestCreate(t *testing.T) {
	handler, mockedSubscriptionStorage := setupTest(t)

	mockedSubscriptionStorage.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(1, nil).
		Times(1)

	body := `{
		"service_name": "test",
		"start_date": "01-2006",
		"end_date": "01-2006",
		"price": 100,
		"user_id": "550e8400-e29b-41d4-a716-446655440000"
	}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/subscriptions", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	handler.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateWrongUserID(t *testing.T) {
	handler, mockedSubscriptionStorage := setupTest(t)

	mockedSubscriptionStorage.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Times(0)

	body := `{
		"service_name": "test",
		"start_date": "01-2006",
		"end_date": "01-2006",
		"price": 100,
		"user_id": "42"
	}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/subscriptions", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	handler.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGet(t *testing.T) {
	handler, mockedSubscriptionStorage := setupTest(t)

	mockedSubscriptionStorage.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(&models.Subscription{ID: 1, ServiceName: "test"}, nil).
		Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/subscriptions/1", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /subscriptions/{id}", handler.Get)

	mux.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdate(t *testing.T) {
	handler, mockedSubscriptionStorage := setupTest(t)

	mockedSubscriptionStorage.EXPECT().
		Update(gomock.Any(), gomock.Any(), &models.Subscription{
			ID:          1,
			ServiceName: "test",
			StartDate:   utils.Must(time.Parse("01-2006", "01-2006")),
			EndDate:     sql.NullTime{Time: utils.Must(time.Parse("01-2006", "01-2006")), Valid: true},
			Price:       100,
			UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		}).
		Return(nil).
		Times(1)

	body := `{
		"service_name": "test",
		"start_date": "01-2006",
		"end_date": "01-2006",
		"price": 100,
		"user_id": "550e8400-e29b-41d4-a716-446655440000"
	}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/subscriptions/1", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /subscriptions/{id}", handler.Update)

	mux.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDelete(t *testing.T) {
	handler, mockedSubscriptionStorage := setupTest(t)

	mockedSubscriptionStorage.EXPECT().
		Delete(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/subscriptions/1", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /subscriptions/{id}", handler.Delete)

	mux.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestList(t *testing.T) {
	handler, mockedSubscriptionStorage := setupTest(t)

	mockedSubscriptionStorage.EXPECT().TotalForPeriod(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil).Times(1)
	mockedSubscriptionStorage.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]models.Subscription{{ID: 1, ServiceName: "test"}}, nil).
		Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/subscriptions", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /subscriptions", handler.List)

	mux.ServeHTTP(w, r)

	type ListResponse struct {
		Data struct {
			Total         int64 `json:"total"`
			Subscriptions []struct {
				ID          int    `json:"id"`
				ServiceName string `json:"service_name"`
			} `json:"subscriptions"`
		} `json:"data"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}

	var resp ListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, int64(1), resp.Data.Total)
	assert.Equal(t, "test", resp.Data.Subscriptions[0].ServiceName)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateInvalidID(t *testing.T) {
	handler, _ := setupTest(t)

	body := `{
		"service_name": "test",
		"start_date": "01-2006",
		"end_date": "01-2006",
		"price": 100,
		"user_id": "550e8400-e29b-41d4-a716-446655440000"
	}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/subscriptions/invalid", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /subscriptions/{id}", handler.Update)

	mux.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateInvalidData(t *testing.T) {
	handler, mockedSubscriptionStorage := setupTest(t)

	mockedSubscriptionStorage.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Times(0)

	body := `{
		"service_name": "",
		"start_date": "13-2006",
		"price": -5
	}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/subscriptions", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	handler.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteNotFound(t *testing.T) {
	handler, mockedSubscriptionStorage := setupTest(t)

	mockedSubscriptionStorage.EXPECT().
		Delete(gomock.Any(), gomock.Any()).
		Return(storage.ErrNotFound).
		Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/subscriptions/999", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /subscriptions/{id}", handler.Delete)

	mux.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetNotFound(t *testing.T) {
	handler, mockedSubscriptionStorage := setupTest(t)

	mockedSubscriptionStorage.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil, storage.ErrNotFound).
		Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/subscriptions/999", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /subscriptions/{id}", handler.Get)

	mux.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateInvalidDateFormat(t *testing.T) {
	handler, mockedSubscriptionStorage := setupTest(t)

	mockedSubscriptionStorage.EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	body := `{
		"service_name": "test",
		"start_date": "2006-01-02",
		"price": 100,
		"user_id": "550e8400-e29b-41d4-a716-446655440000"
	}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/subscriptions/1", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /subscriptions/{id}", handler.Update)

	mux.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}