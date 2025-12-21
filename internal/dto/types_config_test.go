package dto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsExistsWorks(t *testing.T) {
	stringType := new(Types)
	err := stringType.Set("parham")
	assert.NoError(t, err)
	err = stringType.Set("wallex")
	assert.NoError(t, err)
	err = stringType.Set("bale")
	assert.NoError(t, err)

	assert.Equal(t, true, stringType.Exists("parham"))
	assert.Equal(t, false, stringType.Exists("test"))

	rings := "[parham,wallex,bale]"
	assert.Equal(t, rings, stringType.String())
}
