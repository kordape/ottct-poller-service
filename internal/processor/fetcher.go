package processor

import (
	"context"
	"fmt"
	"time"
)

const (
	fetchTweetsMinResults = 5
	fetchTweetsMaxResults = 100
)

type TweetsFetcher interface {
	FetchTweets(context.Context, FetchTweetsRequest) (FetchTweetsResponse, error)
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
	CreatedAt string
}

func (request FetchTweetsRequest) validate() error {
	if request.MaxResults < fetchTweetsMinResults || request.MaxResults > fetchTweetsMaxResults {
		return fmt.Errorf("invalid max results parameter - can range from %d to %d", fetchTweetsMinResults, fetchTweetsMaxResults)
	}

	if request.StartTime.After(request.EndTime) {
		return fmt.Errorf("start time is after end time")
	}

	return nil
}
