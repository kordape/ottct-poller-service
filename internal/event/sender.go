package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kordape/ottct-poller-service/pkg/logger"
	"github.com/kordape/ottct-poller-service/pkg/sqs"

	msg "github.com/kordape/ottct-main-service/pkg/sqs"
)

type FakeNews struct {
	EntityId  string
	Timestamp time.Time
	Content   string
}

type SendFakeNewsEventFn func(ctx context.Context, events []FakeNews) error

// TODO: refactor this into sending batch of messages to SQS
func SendFakeNewsEventFnBuilder(client sqs.Client, log logger.Interface) SendFakeNewsEventFn {
	return func(ctx context.Context, events []FakeNews) error {
		for _, e := range events {
			raw, err := encodeEvent(toSQSEvent(e))
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

func encodeEvent(e msg.FakeNewsEvent) (string, error) {
	b, err := json.Marshal(&e)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func toSQSEvent(e FakeNews) msg.FakeNewsEvent {
	return msg.FakeNewsEvent{
		TweetContent:   e.Content,
		EntityID:       e.EntityId,
		TweetTimestamp: e.Timestamp,
	}
}
