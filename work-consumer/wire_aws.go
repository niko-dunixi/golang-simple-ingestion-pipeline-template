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

func NewAwsSqsQueueSubscription(ctx context.Context, queueURL string) (*pubsub.Subscription, error) {
	queueURL = strings.Replace(queueURL, "https://", "awssqs://", 1)
	subscription, err := pubsub.OpenSubscription(ctx, queueURL)
	if err != nil {
		return nil, fmt.Errorf("could not initialize subscription (consuming side of queue) with aws sqs: %w", err)
	}
	return subscription, nil
}

func InitializeQueueSubscription(ctx context.Context, queueURL string) (*pubsub.Subscription, error) {
	wire.Build(NewAwsSqsQueueSubscription)
	return &pubsub.Subscription{}, nil
}
