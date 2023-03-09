package processor

import (
	"context"
	"errors"
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

func GetProcessFn(fetcher TweetsFetcher, classifier FakeNewsClassifier) ProcessFn {
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

		classifyRequest := make(ClassifyRequest, len(tweets))
		for i, t := range tweets {
			classifyRequest[i] = t.Text
		}

		// Classify tweets as fake or not
		classifyResponse, err := classifier.Classify(ctx, classifyRequest)
		if err != nil {
			return JobResult{
				EntityID: request.EntityID,
				Error:    err,
			}
		}

		if len(classifyResponse.Classification) != len(tweets) {
			return JobResult{
				EntityID: request.EntityID,
				Error:    errors.New("different number of predictions and tweets"),
			}
		}

		fakeTweets := []FakeNewsTweet{}
		// Filter out only fake tweets
		for i, c := range classifyResponse.Classification {
			if c == Fake {
				fakeTweets = append(fakeTweets, FakeNewsTweet{
					Content:   tweets[i].Text,
					Timestamp: tweets[i].CreatedAt,
				})
			}
		}

		return JobResult{
			EntityID:       request.EntityID,
			FakeNewsTweets: fakeTweets,
		}
	}
}
