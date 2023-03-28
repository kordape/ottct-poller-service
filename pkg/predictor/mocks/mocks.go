package mocks

import _ "embed"

var (
	// SuccessResponse represents a success response from Predictor.
	//go:embed success.json
	SuccessResponse string
	// FailResponse represents a fail response from Predictor.
	//go:embed failure.json
	FailResponse string
)
