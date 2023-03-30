package twitter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kordape/ottct-poller-service/pkg/logger"
)

const (
	fetchTweetsMinResults = 5
	fetchTweetsMaxResults = 100
)

//go:generate mockery --inpackage --case snake --disable-version-string --name "TweetsFetcher"
type TweetsFetcher interface {
	FetchTweets(context.Context, logger.Interface, FetchTweetsRequest) (FetchTweetsResponse, error)
}

type FetchTweetsRequest struct {
	MaxResults int
	EntityID   string
	StartTime  time.Time
	EndTime    time.Time
}

type FetchTweetsResponse []Tweet

type Tweet struct {
	ID        string
	Text      string
	CreatedAt time.Time
}

// Make sure Client implement TweetsFetcher interface
var _ TweetsFetcher = &Client{}

type Client struct {
	httpClient  *http.Client
	bearerToken string
}

func New(client *http.Client, bearerToken string) *Client {
	return &Client{
		bearerToken: bearerToken,
		httpClient:  client,
	}
}

func (request FetchTweetsRequest) Validate() error {
	if request.MaxResults < fetchTweetsMinResults || request.MaxResults > fetchTweetsMaxResults {
		return fmt.Errorf("invalid max results parameter - can range from %d to %d", fetchTweetsMinResults, fetchTweetsMaxResults)
	}

	if request.StartTime.After(request.EndTime) {
		return fmt.Errorf("start time is after end time")
	}

	return nil
}
