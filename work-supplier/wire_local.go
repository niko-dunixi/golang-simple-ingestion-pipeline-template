//go:build wireinject && !aws

package main

import (
	"context"
	"fmt"

	"github.com/google/wire"
	"github.com/niko-dunixi/golang-simple-ingestion-pipeline-template/lib/envutil"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"gocloud.dev/docstore"
	_ "gocloud.dev/docstore/mongodocstore"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

func NewRabbitMQSink(ctx context.Context, queueURL string) (*pubsub.Topic, error) {
	rabbitServerURL, err := envutil.GetOrErr(ctx, "RABBIT_SERVER_URL")
	if err != nil {
		return &pubsub.Topic{}, err
	}

	if err := initailizeRabbitMQ(ctx, rabbitServerURL); err != nil {
		return nil, fmt.Errorf("could not initialize the RabbitMQ configuration: %w", err)
	}

	topic, err := pubsub.OpenTopic(ctx, queueURL)
	if err != nil {
		return nil, fmt.Errorf("could not initialize topic (producing side of queue) with aws sqs: %w", err)
	}
	return topic, nil
}

func initailizeRabbitMQ(ctx context.Context, rabbitServerURL string) error {
	log := zerolog.Ctx(ctx)
	conn, err := amqp.Dial(rabbitServerURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open a channel")
	}
	defer ch.Close()
	log.Debug().Msg("creating RabbitMQ exchange")
	exchangeName := "data-ingress"
	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	log.Debug().Msg("creating RabbitMQ queue")
	queue, err := ch.QueueDeclare(
		"data-egress",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	log.Debug().Str("queue_name", queue.Name).Msg("creating RabbitMQ binding for queue and exchange")
	err = ch.QueueBind(queue.Name, "tasks", exchangeName, true, nil)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %w", err)
	}
	return nil
}

func InitializeQueueSink(ctx context.Context, queueURL string) (*pubsub.Topic, error) {
	wire.Build(NewRabbitMQSink)
	return &pubsub.Topic{}, nil
}

func NewMongoCollection(ctx context.Context, collectionURL string) (*docstore.Collection, error) {
	return docstore.OpenCollection(ctx, collectionURL)
}

func InitializeCollection(ctx context.Context, collectionURL string) (*docstore.Collection, error) {
	wire.Build(NewMongoCollection)
	return &docstore.Collection{}, nil
}
