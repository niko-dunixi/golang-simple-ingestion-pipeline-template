package zerologutil

import (
	"os"

	"github.com/rs/zerolog"
)

func init() {
	logger := zerolog.New(os.Stderr)
	zerolog.DefaultContextLogger = &logger
}
