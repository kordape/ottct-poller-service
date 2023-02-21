package processor

import (
	"context"
	"time"
)

type JobResult struct {
	EntityId  string
	Error     error
	Timestamp int64
	Content   string
}

type JobResults []JobResult

type ProcessEntityFn func(ctx context.Context, entityId string) JobResult

func GetProcessEntityFn() ProcessEntityFn {
	// TODO: replace with proccessor that fetches tweets, classifies and filters out fake news tweets
	return func(ctx context.Context, entityId string) JobResult {
		return JobResult{
			EntityId:  entityId,
			Timestamp: time.Now().Unix(),
			Content:   "Dummy title",
		}
	}
}
