//go:build !generate

package loremmarkov

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/mb-14/gomarkov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer ctxCancel()
	chain := New(ctx)
	output, err := chain.Generate(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, output)
	assert.False(t, strings.HasPrefix(output, gomarkov.StartToken), "contained gomarkov.StartToken")
	assert.False(t, strings.HasSuffix(output, gomarkov.EndToken), "contained gomarkov.EndToken")
}
