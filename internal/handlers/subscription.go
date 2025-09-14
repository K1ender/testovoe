package handlers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"testovoe/internal/models"
	"testovoe/internal/response"
	"testovoe/internal/storage"
	"testovoe/internal/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	store    *storage.Storage
	validate *validator.Validate
}


func NewSubscriptionHandler(
	store *storage.Storage,
	validate *validator.Validate,
) *SubscriptionHandler {
	return &SubscriptionHandler{store: store, validate: validate}
}

type SubscriptionResponse struct {
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type CreateSubscriptionPayload struct {
	ServiceName string    `json:"service_name" validate:"required"`
	Price       int       `json:"price" validate:"required"`
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid"`
	StartDate   string    `json:"start_date" validate:"required,mm_yyyy"`
	EndDate     *string   `json:"end_date,omitempty" validate:"omitempty,mm_yyyy"`
}

// Create handles the creation of a new subscription.
// It reads and validates the request payload, parses the start and end dates,
// and stores the subscription in the database. Returns the ID of the created subscription.
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "create subscription", "method", r.Method)

	var payload CreateSubscriptionPayload
	if err := utils.ReadJSON(r, &payload); err != nil {
		slog.ErrorContext(ctx, "read json", "error", err)
		response.BadRequest(w, "Bad request")
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		slog.ErrorContext(ctx, "validate", "error", err)
		if verrs, ok := err.(validator.ValidationErrors); ok {
			response.ValidationError(w, verrs)
		} else {
			response.BadRequest(w, "Invalid input")
		}
		return
	}

	startDate, err := utils.ParseMonthYear(payload.StartDate)
	if err != nil {
		slog.ErrorContext(ctx, "parse start date", "error", err)
		response.BadRequest(w, "Bad request")
		return
	}

	var endDate sql.NullTime = sql.NullTime{
		Valid: false,
	}
	if payload.EndDate != nil {
		end, err := utils.ParseMonthYear(*payload.EndDate)
		if err != nil {
			slog.ErrorContext(ctx, "parse end date", "error", err)
			response.BadRequest(w, "Bad request")
			return
		}
		endDate = sql.NullTime{
			Time:  end,
			Valid: true,
		}
	}

	sub := &models.Subscription{
		ServiceName: payload.ServiceName,
		Price:       payload.Price,
		UserID:      payload.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	id, err := h.store.Subscription.Create(ctx, sub)
	if err != nil {
		slog.ErrorContext(ctx, "create subscription", "error", err)
		response.ServerError(w, "Internal server error")
		return
	}

	response.Created(w, map[string]any{
		"id": id,
	})
}

// Get retrieves a subscription by its ID and returns it in the response.
// It parses the subscription ID from the URL path, fetches the subscription from the store,
// formats the dates, and sends the response.
// If the ID is invalid, it returns a bad request error.
// If the subscription is not found, it returns a not found error.
// If there is an internal server error, it returns a server error.
func (h *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		slog.ErrorContext(ctx, "parse id", "error", err)
		response.BadRequest(w, "Bad request")
		return
	}

	sub, err := h.store.Subscription.Get(ctx, intID)
	if err != nil {
		slog.ErrorContext(ctx, "get subscription", "error", err)
		if errors.Is(err, storage.ErrNotFound) {
			response.NotFound(w, "Not found")
			return
		}
		response.ServerError(w, "Internal server error")
		return
	}

	var endDateFormated *string = nil
	if sub.EndDate.Valid {
		endDateFormated = utils.String(sub.EndDate.Time.Format("01-2006"))
	}

	response.Success(w, SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate.Format("01-2006"),
		EndDate:     endDateFormated,
	})
}

type UpdateSubscriptionPayload struct {
	ServiceName string    `json:"service_name" validate:"required"`
	Price       int       `json:"price" validate:"required"`
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid"`
	StartDate   string    `json:"start_date" validate:"required,mm_yyyy"`
	EndDate     *string   `json:"end_date,omitempty" validate:"omitempty,mm_yyyy"`
}

// Update handles the HTTP request to update an existing subscription.
// It parses the subscription ID from the URL, decodes the JSON payload,
// validates the payload, parses the start and end dates, and updates
// the subscription in the database. It returns appropriate HTTP responses
// based on the outcome of these operations.
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		slog.ErrorContext(ctx, "parse id", "error", err)
		response.BadRequest(w, "Bad request")
		return
	}

	var payload UpdateSubscriptionPayload
	if err := utils.ReadJSON(r, &payload); err != nil {
		response.BadRequest(w, "Bad request")
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		response.ValidationError(w, err.(validator.ValidationErrors))
		return
	}

	startDate, err := utils.ParseMonthYear(payload.StartDate)
	if err != nil {
		slog.ErrorContext(ctx, "parse start date", "error", err)
		response.BadRequest(w, "Bad request")
		return
	}

	var endDate sql.NullTime = sql.NullTime{
		Valid: false,
	}
	if payload.EndDate != nil {
		end, err := utils.ParseMonthYear(*payload.EndDate)
		if err != nil {
			slog.ErrorContext(ctx, "parse end date", "error", err)
			response.BadRequest(w, "Bad request")
			return
		}
		endDate = sql.NullTime{
			Time:  end,
			Valid: true,
		}
	}

	sub := &models.Subscription{
		ID:          intID,
		ServiceName: payload.ServiceName,
		Price:       payload.Price,
		UserID:      payload.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := h.store.Subscription.Update(ctx, intID, sub); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			response.NotFound(w, "Not found")
			return
		}
		response.ServerError(w, "Internal server error")
		return
	}

	var endDateFormated *string = nil
	if sub.EndDate.Valid {
		endDateFormated = utils.String(sub.EndDate.Time.Format("01-2006"))
	}

	response.Success(w, SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate.Format("01-2006"),
		EndDate:     endDateFormated,
	})
}

// Delete removes a subscription by its ID.
// It extracts the ID from the request path, converts it to an integer,
// and attempts to delete the corresponding subscription from the store.
// If the ID is invalid, it responds with a bad request error.
// If the subscription is not found, it responds with a not found error.
// If an internal error occurs during deletion, it responds with a server error.
// On successful deletion, it responds with a no content status.
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		slog.ErrorContext(ctx, "parse id", "error", err)
		response.BadRequest(w, "Bad request")
		return
	}

	if err := h.store.Subscription.Delete(ctx, intID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			response.NotFound(w, "Not found")
			return
		}
		response.ServerError(w, "Internal server error")
		return
	}

	response.NoContent(w)
}

// List handles the HTTP request to list subscriptions based on query parameters.
// It retrieves subscriptions from the store and returns them in a structured JSON response.
// Query parameters:
//   - user_id: filters subscriptions by user ID
//   - service_name: filters subscriptions by service name
//   - limit: limits the number of subscriptions returned (default: 0, meaning no limit)
//   - offset: skips the specified number of subscriptions (default: 0)
//
// If limit or offset are negative, they are reset to 0.
// If userID is invalid, it defaults to uuid.Nil.
// If an error occurs during processing, an appropriate HTTP error response is sent.
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		slog.WarnContext(ctx, "parse limit", "error", err)
		limit = 0
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		slog.WarnContext(ctx, "parse offset", "error", err)
		offset = 0
	}
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}

	subscriptions, err := h.store.Subscription.List(ctx, userID, serviceName, limit, offset)
	if err != nil {
		slog.ErrorContext(ctx, "list subscriptions", "error", err)
		response.ServerError(w, "Internal server error")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		slog.WarnContext(ctx, "parse user id", "error", err)
		userUUID = uuid.Nil
	}

	total, err := h.store.Subscription.TotalForPeriod(ctx, time.Now(), time.Now(), userUUID, serviceName)
	if err != nil {
		slog.ErrorContext(ctx, "list subscriptions", "error", err)
		response.ServerError(w, "Internal server error")
		return
	}

	if len(subscriptions) == 0 {
		subscriptions = []models.Subscription{}
	}

	var resp []SubscriptionResponse
	for _, sub := range subscriptions {
		var endDateFormated *string = nil
		if sub.EndDate.Valid {
			endDateFormated = utils.String(sub.EndDate.Time.Format("01-2006"))
		}
		resp = append(resp, SubscriptionResponse{
			ID:          sub.ID,
			ServiceName: sub.ServiceName,
			Price:       sub.Price,
			UserID:      sub.UserID,
			StartDate:   sub.StartDate.Format("01-2006"),
			EndDate:     endDateFormated,
		})
	}

	response.Success(w, map[string]any{
		"subscriptions": resp,
		"total":         total,
	})
}
