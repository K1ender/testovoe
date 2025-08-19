package response

import (
	"net/http"
	"testovoe/internal/utils"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func BadRequest(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusBadRequest)
	return utils.WriteJSON(w, http.StatusBadRequest, Response{
		Status:  http.StatusBadRequest,
		Message: message,
	})
}

type ValidationErr struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value any    `json:"value,omitempty"`
}

func ValidationError(w http.ResponseWriter, errs validator.ValidationErrors) error {
	w.WriteHeader(http.StatusBadRequest)
	var out []ValidationErr
	for _, err := range errs {
		out = append(out, ValidationErr{
			Field: err.Field(),
			Tag:   err.Tag(),
			Value: err.Value(),
		})
	}
	return utils.WriteJSON(w, http.StatusBadRequest, Response{
		Status:  http.StatusBadRequest,
		Message: "validation error",
		Data:    out,
	})
}

func ServerError(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusInternalServerError)
	return utils.WriteJSON(w, http.StatusInternalServerError, Response{
		Status:  http.StatusInternalServerError,
		Message: message,
	})
}

func Success(w http.ResponseWriter, data any) error {
	return utils.WriteJSON(w, http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func NotFound(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusNotFound)
	return utils.WriteJSON(w, http.StatusNotFound, Response{
		Status:  http.StatusNotFound,
		Message: message,
	})
}

func Created(w http.ResponseWriter, data any) error {
	w.WriteHeader(http.StatusCreated)
	return utils.WriteJSON(w, http.StatusCreated, Response{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    data,
	})
}

func NoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
