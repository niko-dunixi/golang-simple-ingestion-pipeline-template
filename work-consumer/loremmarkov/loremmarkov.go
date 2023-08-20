//go:build !generate

package loremmarkov

import (
	"context"
	_ "embed"
	"strings"
	"unicode"

	"github.com/mb-14/gomarkov"
	"github.com/rs/zerolog"
)

//go:embed lorem-ipsum.json
var markovChainBytes []byte

type Chain struct {
	chain *gomarkov.Chain
}

func New(ctx context.Context) Chain {
	log := zerolog.Ctx(ctx)
	chain := gomarkov.NewChain(0)
	if err := chain.UnmarshalJSON(markovChainBytes); err != nil {
		log.Fatal().Err(err).Msg("Could not restore markov chain, there is no possible way to recover")
	}
	return Chain{
		chain: chain,
	}
}

func (c Chain) Generate(ctx context.Context) (string, error) {
	resultChannel := make(chan struct {
		s   string
		err error
	})
	defer close(resultChannel)

	sb := strings.Builder{}
	tokens := []string{gomarkov.StartToken}
	for {
		// If the context has been canceled, abort
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}
		// Process chain
		next, err := c.chain.Generate(tokens[(len(tokens) - 1):])
		if err != nil {
			return "", err
		}
		tokens = append(tokens, next)
		if next == gomarkov.EndToken {
			// We've hit the end of processing. No more logic to be done.
			return sb.String(), nil
		} else if len(tokens) > 2 && !(len(next) == 1 && unicode.IsPunct(rune(next[0]))) {
			// If we've already inserted the first word (being sure to account for the StartToken)
			// we will insert a space unless it's punctuation
			sb.WriteRune(' ')
		}
		sb.WriteString(next)
	}
}
