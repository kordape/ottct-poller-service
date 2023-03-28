package predictor

import (
	"context"
	"net/http"
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

// Make sure Client implement FakeNewsClassifier interface
var _ FakeNewsClassifier = &Client{}

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func New(client *http.Client, baseURL string) *Client {
	return &Client{
		httpClient: client,
		baseURL:    baseURL,
	}
}
