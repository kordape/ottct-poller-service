package processor

import "context"

type JobResult struct {
	EntityId string
	Error    error
}

type JobResults []JobResult

type ProcessEntityFn func(ctx context.Context) JobResult

func GetProcessEntityFn(entityId string) ProcessEntityFn {
	// TODO: replace with proccessor that fetches tweets, classifies and filters out fake news tweets
	return func(ctx context.Context) JobResult {
		return JobResult{
			EntityId: entityId,
		}
	}
}
