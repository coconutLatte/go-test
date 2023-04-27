package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestA(t *testing.T) {
	input := []byte{0x00, 0x00, 0x01, 0x01, 0x02, 0x03, 0x00, 0x00, 0x00, 0x01, 0x04, 0x05, 0x00, 0x00, 0x01, 0x06}
	nalus, err := splitNalus(input)
	assert.NoError(t, err)
	fmt.Println(nalus)
}
