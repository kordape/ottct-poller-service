package processor

import (
	"context"
	"time"
)

const (
	defaultFetchCount = 5
)

type JobRequest struct {
	EntityID  string
	StartTime time.Time
	EndTime   time.Time
}

type JobResult struct {
	EntityID       string
	Error          error
	FakeNewsTweets []FakeNewsTweet
}

type FakeNewsTweet struct {
	Content   string
	Timestamp string
}

type JobResults []JobResult

type ProcessFn func(ctx context.Context, request JobRequest) JobResult

func GetProcessFn(fetcher TweetsFetcher, classifier TweetsClassifier) ProcessFn {
	return func(ctx context.Context, request JobRequest) JobResult {
		// Fetch tweets in given time window
		fetchRequest := FetchTweetsRequest{
			EntityID:   request.EntityID,
			StartTime:  request.StartTime,
			EndTime:    request.EndTime,
			MaxResults: defaultFetchCount,
		}

		if err := fetchRequest.validate(); err != nil {
			return JobResult{
				EntityID: request.EntityID,
				Error:    err,
			}
		}

		tweets, err := fetcher.FetchTweets(ctx, fetchRequest)
		if err != nil {
			return JobResult{
				EntityID: request.EntityID,
				Error:    err,
			}
		}

		fakeTweets := []FakeNewsTweet{}
		for _, tweet := range tweets {
			fakeTweets = append(fakeTweets, FakeNewsTweet{
				Content:   tweet.Text,
				Timestamp: tweet.CreatedAt,
			})
		}

		return JobResult{
			EntityID:       request.EntityID,
			FakeNewsTweets: fakeTweets,
		}
	}
}
