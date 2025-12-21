package visitors

import (
	"github.com/pm1381/sirish/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCommentVisitor(t *testing.T) {
	expectedTargets := []string{
		"MagicNamedParamsAndResults",
		"MagicUnnamedAndNamedParamsAndResults",
		"MagicUnderscoreNames",
		"MagicNoParams",
		"MagicNoResult",
	}

	pathAbs := internal.GetTestPathHelper("comment.go", "visitors")
	commentModule := NewCommentVisitor(pathAbs)
	require.NotNil(t, commentModule)
	commentModule.Traverse()

	assert.Equal(t, len(expectedTargets), len(commentModule.GetTargets()))
	assert.ElementsMatch(t, expectedTargets, commentModule.GetTargets())
}
