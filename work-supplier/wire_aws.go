//go:build wireinject && aws

package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/wire"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/awssnssqs"
)

func NewAwsSqsQueueTopic(ctx context.Context, queueURL string) (*pubsub.Topic, error) {
	// - https://gocloud.dev/howto/pubsub/publish/#sqs
	queueURL = strings.Replace(queueURL, "https://", "awssqs://", 1)

	topic, err := pubsub.OpenTopic(ctx, queueURL)
	if err != nil {
		return nil, fmt.Errorf("could not initialize topic (producing side of queue) with aws sqs: %w", err)
	}
	return topic, nil
}

func InitializeQueueSink(ctx context.Context, queueURL string) (*pubsub.Topic, error) {
	wire.Build(NewAwsSqsQueueTopic)
	return &pubsub.Topic{}, nil
}
