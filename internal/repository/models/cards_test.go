package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCard_Getters(t *testing.T) {
	c := Card{
		ID:        "123",
		UID:       "user456",
		Name:      "John Doe",
		Number:    "4111111111111111",
		Month:     12,
		Year:      2025,
		CVC:       "123",
		PIN:       "4567",
		Meta:      "meta",
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name     string
		method   func() string
		expected string
	}{
		{"GetID", c.GetID, "123"},
		{"GetUID", c.GetUID, "user456"},
		{"GetName", c.GetName, "John Doe"},
		{"GetNumber", c.GetNumber, "4111111111111111"},
		{"GetCVC", c.GetCVC, "123"},
		{"GetPIN", c.GetPIN, "4567"},
		{"GetMeta", c.GetMeta, "meta"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.method())
		})
	}

	t.Run("GetMonth", func(t *testing.T) {
		assert.Equal(t, 12, c.GetMonth())
	})

	t.Run("GetYear", func(t *testing.T) {
		assert.Equal(t, 2025, c.GetYear())
	})

	t.Run("GetCreatedAt", func(t *testing.T) {
		assert.Equal(t, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), c.GetCreatedAt())
	})
}

func TestCard_MaskNumber(t *testing.T) {
	tests := []struct {
		number   string
		expected string
	}{
		{"4111111111111111", "**** **** **** 1111"},
		{"1234567812345678", "**** **** **** 5678"},
		{"5555555555554444", "**** **** **** 4444"},
	}

	for _, tt := range tests {
		t.Run(tt.number, func(t *testing.T) {
			c := &Card{Number: tt.number}
			c.MaskNumber()
			assert.Equal(t, tt.expected, c.GetMask())
		})
	}
}

func TestCard_Validate(t *testing.T) {
	validCard := Card{
		Month:  12,
		Year:   2025,
		CVC:    "123",
		PIN:    "4567",
		Number: "4111111111111111",
	}

	t.Run("Valid card", func(t *testing.T) {
		assert.NoError(t, validCard.Validate())
	})

	tests := []struct {
		name   string
		modify func(*Card)
		errMsg string
	}{
		{
			"Invalid month (0)",
			func(c *Card) { c.Month = 0 },
			"invalid month value",
		},
		{
			"Invalid month (13)",
			func(c *Card) { c.Month = 13 },
			"invalid month value",
		},
		{
			"Invalid year (1969)",
			func(c *Card) { c.Year = 1969 },
			"invalid year value",
		},
		{
			"Invalid CVC length (2)",
			func(c *Card) { c.CVC = "12" },
			"invalid cvc value",
		},
		{
			"Invalid CVC length (4)",
			func(c *Card) { c.CVC = "1234" },
			"invalid cvc value",
		},
		{
			"Invalid PIN length (3)",
			func(c *Card) { c.PIN = "123" },
			"invalid pin value",
		},
		{
			"Invalid PIN length (7)",
			func(c *Card) { c.PIN = "1234567" },
			"invalid pin value",
		},
		{
			"Invalid card number (Luhn fail)",
			func(c *Card) { c.Number = "4111111111111112" },
			"invalid number value",
		},
		{
			"Short card number (12 digits)",
			func(c *Card) { c.Number = "123456789012" },
			"invalid number value",
		},
		{
			"Long card number (20 digits)",
			func(c *Card) { c.Number = "12345678901234567890" },
			"invalid number value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := validCard
			tt.modify(&c)
			err := c.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}
