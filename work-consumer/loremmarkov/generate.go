//go:build generate

//go:generate go run generate.go

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mb-14/gomarkov"
	"github.com/rs/zerolog"
)

// Original text from:
//- https://en.wikipedia.org/wiki/Lorem_ipsum
//

func main() {
	log := zerolog.New(os.Stderr)
	file, err := os.Open("lorem-ipsum.txt")
	if err != nil {
		log.Fatal().Err(err).Msg("this should always work and never panic")
	}
	log.Debug().Msg("opened lorem-ipsum.txt")
	chain := gomarkov.NewChain(1)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		tokens := strings.Split(currentLine, " ")
		chain.Add(tokens)
	}
	log.Debug().Msg("consumed contents")
	chainBytes, err := chain.MarshalJSON()
	if err != nil {
		log.Fatal().Err(err).Msg("could marshal markov chain to json")
		panic(fmt.Sprintf(": %v", err))
	}
	log.Debug().Msg("storing lorem-ipsum model as json")
	if err := os.WriteFile("lorem-ipsum.json", chainBytes, 0644); err != nil {
		log.Fatal().Err(err).Msg("could save json to disk for embedding")
	}
	log.Debug().Msg("completed successfully")
}
