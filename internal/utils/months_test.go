package utils

import (
    "testing"
    "time"
)

func TestMonthsOverlap_BasicOverlap(t *testing.T) {
    aStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    aEnd := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
    bStart := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
    bEnd := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
    got := MonthsOverlap(aStart, aEnd, bStart, bEnd)
    if got != 2 {
        t.Fatalf("expected 2 overlapping months, got %d", got)
    }
}

func TestMonthsOverlap_NoOverlap(t *testing.T) {
    aStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    aEnd := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    bStart := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
    bEnd := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
    got := MonthsOverlap(aStart, aEnd, bStart, bEnd)
    if got != 0 {
        t.Fatalf("expected 0 overlapping months, got %d", got)
    }
}

func TestMonthsOverlap_IdenticalRange(t *testing.T) {
    aStart := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
    aEnd := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
    bStart := aStart
    bEnd := aEnd
    got := MonthsOverlap(aStart, aEnd, bStart, bEnd)
    if got != 3 {
        t.Fatalf("expected 3 overlapping months, got %d", got)
    }
}

func TestMonthsOverlap_InvertedInputs(t *testing.T) {
    aStart := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
    aEnd := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
    bStart := time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)
    bEnd := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
    got := MonthsOverlap(aStart, aEnd, bStart, bEnd)
    if got != 0 {
        t.Fatalf("expected 0 overlapping months when ranges inverted, got %d", got)
    }
}

func TestMonthsOverlap_TouchingEdges(t *testing.T) {
    // a: Jan..Feb, b: Mar..Apr -> touching but not overlapping
    aStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    aEnd := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
    bStart := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
    bEnd := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
    got := MonthsOverlap(aStart, aEnd, bStart, bEnd)
    if got != 0 {
        t.Fatalf("expected 0 overlapping months for touching edges, got %d", got)
    }
}

