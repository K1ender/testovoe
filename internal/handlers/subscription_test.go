package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testovoe/internal/handlers"
	"testovoe/internal/storage"
	mock_storage "testovoe/internal/storage/mocks"
	"testovoe/internal/validators"

	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("mm_yyyy", validators.MonthYearValidator)

	mockedSubscriptionStorage := mock_storage.NewMockSubscriptionStorage(ctrl)
	mockedSubscriptionStorage.EXPECT().Create(gomock.Any(), gomock.Any()).Return(1, nil).Times(1)
	handler := handlers.NewSubscriptionHandler(
		&storage.Storage{
			Subscription: mockedSubscriptionStorage,
		},
		validate,
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/subscriptions", strings.NewReader(`{
        "service_name": "test",
        "start_date": "01-2006",
        "end_date": "01-2006",
        "price": 100,
        "user_id": "550e8400-e29b-41d4-a716-446655440000"
    }`))
	r.Header.Set("Content-Type", "application/json")

	handler.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
}
