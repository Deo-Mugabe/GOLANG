package domain

import "errors"

var (
	// Not Found Errors
	ErrBookingNotFound  = errors.New("booking not found")
	ErrPrisonerNotFound = errors.New("prisoner not found")
	ErrChargeNotFound   = errors.New("charge not found")
	ErrArrestNotFound   = errors.New("arrest not found")

	// Validation Errors
	ErrInvalidBookingID   = errors.New("invalid booking ID")
	ErrInvalidNameID      = errors.New("invalid name ID")
	ErrInvalidBookingDate = errors.New("invalid booking date")

	// Business Logic Errors
	ErrBookingNotActive       = errors.New("booking is not active")
	ErrCannotProcessBooking   = errors.New("booking cannot be processed")
	ErrMissingRequiredData    = errors.New("missing required data")
	ErrInvalidProcessingRange = errors.New("invalid processing time range")
)
