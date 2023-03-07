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
	StartTime  string
	EndTime    string
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

	if request.StartTime != "" && request.EndTime != "" {
		start, err := time.Parse(time.RFC3339, request.StartTime)
		if err != nil {
			return fmt.Errorf("error parsing start time: %s", err)
		}

		end, err := time.Parse(time.RFC3339, request.EndTime)
		if err != nil {
			return fmt.Errorf("error parsing end time: %s", err)
		}

		if start.After(end) {
			return fmt.Errorf("start time is after end time")
		}
	}

	return nil
}
