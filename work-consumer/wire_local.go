//go:build wireinject && !aws

package main

import (
	"context"
	"fmt"

	"github.com/google/wire"
	"github.com/niko-dunixi/golang-simple-ingestion-pipeline-template/lib/envutil"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

func NewRabbitMQSubscription(ctx context.Context, queueURL string) (*pubsub.Subscription, error) {
	if err := envutil.HasOrErr(ctx, "RABBIT_SERVER_URL"); err != nil {
		return &pubsub.Subscription{}, err
	}
	subscription, err := pubsub.OpenSubscription(ctx, queueURL)
	if err != nil {
		return nil, fmt.Errorf("could not initialize subscription (consuming side of queue) with aws sqs: %w", err)
	}
	return subscription, nil
}

func InitializeQueueSubscription(ctx context.Context, queueURL string) (*pubsub.Subscription, error) {
	wire.Build(NewRabbitMQSubscription)
	return &pubsub.Subscription{}, nil
}
