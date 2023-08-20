package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/niko-dunixi/golang-simple-ingestion-pipeline-template/lib"
	"github.com/niko-dunixi/golang-simple-ingestion-pipeline-template/lib/envutil"
	_ "github.com/niko-dunixi/golang-simple-ingestion-pipeline-template/lib/zerologutil"
	"github.com/rs/zerolog"
	"gocloud.dev/pubsub"
	"golang.org/x/sync/errgroup"
)

func main() {
	initCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	initCtx = zerolog.Ctx(initCtx).With().Str("scope", "initialization").Logger().WithContext(initCtx)
	defer cancel()
	initLog := zerolog.Ctx(initCtx)

	maxConcurrentCount := envutil.Int(initCtx, "MAX_CONCURRENT_COUNT", 30)
	queueURL := envutil.Must(initCtx, "QUEUE_URL")
	queue, err := InitializeQueueSubscription(initCtx, queueURL)
	if err != nil {
		initLog.Fatal().Err(err).Msg("failed to initialize queue client")
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = zerolog.Ctx(context.Background()).With().Str("scope", "working").Logger().WithContext(ctx)
	defer cancel()

	messagesChannel := make(chan *pubsub.Message, maxConcurrentCount)
	defer close(messagesChannel)

	signals := make(chan os.Signal, 1)
	defer close(signals)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		cancelCount := 0
		for signal := range signals {
			log := zerolog.Ctx(ctx).With().Str("signal", signal.String()).Logger()
			if cancelCount == 0 {
				log.Info().Msg("caught signal, initiating shutdown")
				cancel()
			} else if cancelCount > 1 {
				log.Fatal().Msg("additional signal caught, forcing shutdown")
			}
			cancelCount += 1
		}
	}()

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(2)
	eg.Go(func() error {
		ctx := zerolog.Ctx(ctx).With().Str("loop", "receiving").Logger().WithContext(ctx)
		if err := receivingLoop(ctx, queue, messagesChannel); err != nil {
			return fmt.Errorf("a problem occurred in the recieving loop: %v", err)
		}
		return nil
	})
	eg.Go(func() error {
		ctx := zerolog.Ctx(ctx).With().Str("loop", "processing").Logger().WithContext(ctx)
		if err := processingLoop(ctx, messagesChannel); err != nil {
			return fmt.Errorf("a problem occurred in the proccessing loop")
		}
		return nil
	})

	initLog.Info().Msg("Starting")
	if err := eg.Wait(); err != nil {
		log := zerolog.Ctx(ctx)
		log.Fatal().Err(err).Msg("a problem occurred and we will now exit")
	}
	initLog.Info().Msg("Exiting")
}

func receivingLoop(ctx context.Context, queue *pubsub.Subscription, messagesChannel chan<- *pubsub.Message) error {
	log := zerolog.Ctx(ctx)
	log.Info().Msg("Starting loop")
	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if errors.Is(err, context.Canceled) {
				log.Info().Msg("shutting down receiving loop")
				return nil
			} else if err != nil {
				return fmt.Errorf("something caused the context to err out: %v", err)
			}
			return nil
		default:
			log.Debug().Msg("attempting to receive from queue")
			message, err := queue.Receive(ctx)
			if err != nil {
				log.Error().Err(err).Msg("could not read from subscription")
			}
			messagesChannel <- message
		}
	}
}

func processingLoop(ctx context.Context, messagesChannel <-chan *pubsub.Message) error {
	log := zerolog.Ctx(ctx)
	log.Info().Msg("Starting loop")
	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if errors.Is(err, context.Canceled) {
				log.Info().Msg("shutting down processing loop")
				return nil
			} else if err != nil {
				return fmt.Errorf("something caused the context to err out: %v", err)
			}
			return nil
		case message := <-messagesChannel:
			msgCtx := log.With().Str("message_id", message.LoggableID).Logger().WithContext(ctx)
			if err := processMessage(msgCtx, message); err != nil {
				log.Error().Err(err).Msg("could not process message")
			}
		}
	}
}

func processMessage(ctx context.Context, message *pubsub.Message) error {
	log := zerolog.Ctx(ctx)
	body := message.Body
	if body == nil {
		return fmt.Errorf("mesage body was nil")
	} else if len(body) == 0 {
		return fmt.Errorf("message body was length zero")
	}
	payload := lib.PayloadItem{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("could not parse JSON from message body: %w", err)
	}
	// All work now done, be sure to acknoledge the message so that it
	// is removed from the queue
	defer message.Ack()
	log.Info().Any("payload", payload).Msg("successfully processed")
	return nil
}
