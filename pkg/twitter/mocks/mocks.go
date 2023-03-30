package mocks

import _ "embed"

var (
	// SuccessResponse represents a success response from Twitter.
	//go:embed success.json
	SuccessResponse string
	// FailResponse represents a fail response from Twitter.
	//go:embed failure.json
	FailResponse string
)
