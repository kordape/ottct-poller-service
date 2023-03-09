package processor

import (
	"context"
)

type Classification int

const (
	Real Classification = 0
	Fake Classification = 1
)

//go:generate mockery --dir=./ --name=TweetsClassifier --filename=classifier.go --output=./mocks  --outpkg=mocks
type TweetsClassifier interface {
	Classify(ctx context.Context, tweets []ClassifyTweetsRequest) (ClassifyTweetsResponse, error)
}

type ClassifyTweetsRequest struct {
	Tweet string
}

type ClassifyTweetsResponse struct {
	Classification []Classification
}
