package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindStartCode(t *testing.T) {
	input := []byte{0x00, 0x00, 0x01, 0x10, 0x12, 0xA1, 0x12}
	sc, exist, err := findStartCode(input)
	assert.NoError(t, err)
	assert.True(t, exist)
	assert.Equal(t, 0, sc.startIndex)
	assert.Equal(t, 2, sc.endIndex)

	input = []byte{0x00, 0x00, 0x00, 0x01, 0x10, 0x12, 0xA1, 0x12}
	sc, exist, err = findStartCode(input)
	assert.NoError(t, err)
	assert.True(t, exist)
	assert.Equal(t, 0, sc.startIndex)
	assert.Equal(t, 3, sc.endIndex)

	input = []byte{0x00, 0x00, 0x00, 0x00, 0x10, 0x12, 0xA1, 0x12}
	sc, _, err = findStartCode(input)
	assert.Error(t, err)

	input = []byte{0x23, 0x00, 0x00, 0x00, 0x01, 0x12, 0xA1, 0x12}
	sc, exist, err = findStartCode(input)
	assert.NoError(t, err)
	assert.True(t, exist)
	assert.Equal(t, 1, sc.startIndex)
	assert.Equal(t, 4, sc.endIndex)

	input = []byte{0x23, 0x00, 0x20, 0x00, 0x01, 0x12, 0xA1, 0x12}
	sc, exist, err = findStartCode(input)
	assert.NoError(t, err)
	assert.False(t, exist)
}

func TestSplitH264BitStream(t *testing.T) {
	q := NewQueue()

	streamInput := [][]byte{
		{0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
		{0x06, 0x07, 0x00, 0x00, 0x00, 0x01, 0x08},
		{0x00, 0x00, 0x00, 0x01, 0x09, 0x0A, 0x0B},
	}

	q.append(streamInput[0])
	q.append(streamInput[1])
	q.append(streamInput[2])

	output := q.pop()
	fmt.Println(output)
	output = q.pop()
	fmt.Println(output)
	output = q.pop()
	fmt.Println(output)
}

func TestBindAll(t *testing.T) {
	index := 0
	f, _ := os.Create("h264_raw_data/video_record")

	for ; index < 51; index++ {
		d, _ := ioutil.ReadFile(fmt.Sprintf("h264_raw_data/%d.bin", index))

		f.Write(d)
	}

	f.Close()
}
