package event

import (
	"github.com/kordape/ottct-poller-service/pkg/sqs"
)

type FakeNews struct {
	EntityId  string
	Timestamp string
	Content   string
}

type SendFakeNewsEventFn func(events []FakeNews) error

func SendFakeNewsEventFnBuilder(client sqs.Client) SendFakeNewsEventFn {
	return func(events []FakeNews) error {
		// TODO: replace with sending events to SQS
		return nil
	}
}
