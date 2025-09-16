package utils

import (
    "testing"
    "time"
)

func TestParseMonthYear_Valid(t *testing.T) {
    got, err := ParseMonthYear("01-2006")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    expected := time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC)
    if !got.Equal(expected) {
        t.Fatalf("expected %v, got %v", expected, got)
    }
}

func TestParseMonthYear_InvalidFormat(t *testing.T) {
    _, err := ParseMonthYear("2006-01")
    if err == nil {
        t.Fatalf("expected error for invalid format, got nil")
    }
}

func TestParseMonthYear_OutOfRangeMonth(t *testing.T) {
    // Sscanf will parse, but constructing the date with month 13 rolls over; ensure validator is used elsewhere.
    // Here we just assert no panic and result is normalized by time.Date behavior.
    got, err := ParseMonthYear("13-2006")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // 13th month of 2006 becomes 2007-01-01 in time.Date
    expected := time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC)
    if !got.Equal(expected) {
        t.Fatalf("expected %v, got %v", expected, got)
    }
}

