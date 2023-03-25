package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kordape/ottct-poller-service/pkg/logger"
	"github.com/kordape/ottct-poller-service/pkg/sqs"
)

type FakeNews struct {
	EntityId  string
	Timestamp string
	Content   string
}

type SendFakeNewsEventFn func(ctx context.Context, events []FakeNews) error

// TODO: refactor this into sending batch of messages to SQS
func SendFakeNewsEventFnBuilder(client sqs.Client, log logger.Interface) SendFakeNewsEventFn {
	return func(ctx context.Context, events []FakeNews) error {
		for _, e := range events {
			raw, err := encodeEvent(e)
			if err != nil {
				log.Error(fmt.Printf("error encoding event: %v", e))
				return fmt.Errorf("error encoding event: %s", err)
			}

			_, err = client.Send(ctx, raw)
			if err != nil {
				log.Error(fmt.Printf("error sending event: %s", e))
				return fmt.Errorf("error sending event to sqs: %s", err)
			}

			log.Error(fmt.Printf("successfully sent event: %s", e))
		}

		return nil
	}
}

func encodeEvent(event FakeNews) (string, error) {
	b, err := json.Marshal(&event)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
