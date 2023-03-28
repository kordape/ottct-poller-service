package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SendOption is a functional option that can augment or modify a sqs.SendMessageInput request.
type SendOption func(*sqs.SendMessageInput)

// WithDelaySeconds returns a SendOption which setup the delay seconds when SendMessage.
func WithDelaySeconds(s int32) SendOption {
	return func(input *sqs.SendMessageInput) {
		input.DelaySeconds = s
	}
}

// Client represents a client that communicates with Amazon SQS about the request.
type Client interface {
	Send(ctx context.Context, msg string, options ...SendOption) (string, error)
}

type client struct {
	SQS *sqs.Client
	URL string
}

// NewClient returns a new SQS client.
func NewClient(sqsAPI *sqs.Client, queueURL string) Client {
	return &client{
		SQS: sqsAPI,
		URL: queueURL,
	}
}

func (c client) Send(ctx context.Context, msg string, options ...SendOption) (string, error) {
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(msg),
		QueueUrl:    aws.String(c.URL),
	}
	for _, option := range options {
		if option != nil {
			option(input)
		}
	}

	output, err := c.SQS.SendMessage(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to send the message into the queue: %w", err)
	}

	return aws.ToString(output.MessageId), nil
}
