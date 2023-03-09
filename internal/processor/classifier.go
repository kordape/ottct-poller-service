package processor

import (
	"context"
)

type Classification int

const (
	Real Classification = 0
	Fake Classification = 1
)

//go:generate mockery --inpackage --case snake --disable-version-string --name "FakeNewsClassifier"
type FakeNewsClassifier interface {
	Classify(ctx context.Context, request ClassifyRequest) (ClassifyResponse, error)
}

type ClassifyRequest []string

type ClassifyResponse struct {
	Classification []Classification
}
