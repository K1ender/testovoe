package response

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/go-playground/validator/v10"
)

func TestBadRequest(t *testing.T) {
    rr := httptest.NewRecorder()
    if err := BadRequest(rr, "bad"); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
    }
    var resp Response
    if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
        t.Fatalf("json: %v", err)
    }
    if resp.Status != http.StatusBadRequest || resp.Message != "bad" {
        t.Fatalf("unexpected payload: %+v", resp)
    }
}

func TestValidationError(t *testing.T) {
    rr := httptest.NewRecorder()
    // Build a minimal fake validation error output by calling ValidationError with an empty list
    // Since function expects validator.ValidationErrors, simulate zero errors by passing an empty list
    var errs validator.ValidationErrors
    if err := ValidationError(rr, errs); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
    }
    var resp Response
    if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
        t.Fatalf("json: %v", err)
    }
    if resp.Status != http.StatusBadRequest || resp.Message != "validation error" {
        t.Fatalf("unexpected payload: %+v", resp)
    }
}

func TestServerError(t *testing.T) {
    rr := httptest.NewRecorder()
    if err := ServerError(rr, "boom"); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if rr.Code != http.StatusInternalServerError {
        t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
    }
}

func TestSuccessAndCreatedAndNoContent(t *testing.T) {
    // Success
    rr := httptest.NewRecorder()
    if err := Success(rr, map[string]string{"ok": "1"}); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if rr.Code != http.StatusOK {
        t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
    }

    // Created
    rr = httptest.NewRecorder()
    if err := Created(rr, map[string]int{"id": 7}); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if rr.Code != http.StatusCreated {
        t.Fatalf("expected %d, got %d", http.StatusCreated, rr.Code)
    }

    // NoContent
    rr = httptest.NewRecorder()
    if err := NoContent(rr); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if rr.Code != http.StatusNoContent {
        t.Fatalf("expected %d, got %d", http.StatusNoContent, rr.Code)
    }
}

