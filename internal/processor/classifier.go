package processor

import (
	"context"
)

type Classification int

const (
	Real Classification = 0
	Fake Classification = 1
)

//go:generate mockery --dir=./ --name=FakeNewsClassifier --filename=classifier.go --output=./mocks  --outpkg=mocks
type FakeNewsClassifier interface {
	Classify(ctx context.Context, request ClassifyRequest) (ClassifyResponse, error)
}

type ClassifyRequest []string

type ClassifyResponse struct {
	Classification []Classification
}
