package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	expectedResult := "hello"
	receiveResult := RemoveSubstring("hello girl", "girl")
	assert.Equal(t, expectedResult, receiveResult)
}
