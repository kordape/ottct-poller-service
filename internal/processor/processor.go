package processor

import (
	"context"
	"time"
)

type JobResult struct {
	EntityId       string
	Error          error
	FakeNewsTweets []FakeNewsTweet
}

type FakeNewsTweet struct {
	Content   string
	Timestamp int64
}

type JobResults []JobResult

type ProcessEntityFn func(ctx context.Context, entityId string) JobResult

func GetProcessEntityFn() ProcessEntityFn {
	// TODO: replace with proccessor that fetches tweets, classifies and filters out fake news tweets
	return func(ctx context.Context, entityId string) JobResult {
		return JobResult{
			EntityId: entityId,
			FakeNewsTweets: []FakeNewsTweet{
				{
					Content:   "Dummy content",
					Timestamp: time.Now().Unix(),
				},
			},
		}
	}
}
