package envutil

import (
	"context"
	"fmt"

	"strconv"

	"os"

	"github.com/rs/zerolog"
)

type AbsentKeyErr struct {
	Key string
}

func (ake AbsentKeyErr) Error() string {
	return fmt.Sprintf("key was absent: %s", ake.Key)
}

func Must(ctx context.Context, key string) string {
	log := zerolog.Ctx(ctx)
	value, err := GetOrErr(ctx, key)
	if err != nil {
		log.Fatal().Err(err).Msg("environment variable was not present")
	}
	return value
}

func GetOrErr(ctx context.Context, key string) (string, error) {
	log := zerolog.Ctx(ctx)
	value, isPresent := os.LookupEnv(key)
	if !isPresent {
		log.Debug().Str("key", key).Msg("environment variable was not present")
		return "", AbsentKeyErr{
			Key: key,
		}
	}
	log.Debug().Str("key", key).Msg("environment variable was present")
	return value, nil
}

func Int(ctx context.Context, key string, defaultValue int) int {
	log := zerolog.Ctx(ctx).With().Str("key", key).Logger()
	value, isPresent := os.LookupEnv(key)
	if !isPresent {
		log.Debug().Msg("falling back to default value for environment variable")
		return defaultValue
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse environment variable")
	}
	log.Info().Msg("parsed environment variable")
	return parsedValue
}

func Has(ctx context.Context, key string) bool {
	log := zerolog.Ctx(ctx).With().Str("key", key).Logger()
	_, isPresent := os.LookupEnv(key)
	if !isPresent {
		log.Info().Msg("key was not present")
	} else {
		log.Info().Msg("key was present")
	}
	return isPresent
}

func HasOrErr(ctx context.Context, key string) error {
	log := zerolog.Ctx(ctx).With().Str("key", key).Logger()
	_, isPresent := os.LookupEnv(key)
	if !isPresent {
		err := AbsentKeyErr{
			Key: key,
		}
		log.Info().Err(err).Msg("key was present")
		return err
	}
	log.Info().Msg("key was present")
	return nil
}

func MustHave(ctx context.Context, key string) {
	log := zerolog.Ctx(ctx).With().Str("key", key).Logger()
	_, isPresent := os.LookupEnv(key)
	if !isPresent {
		log.Fatal().Msg("key was not present")
	} else {
		log.Info().Msg("key was present")
	}
}
