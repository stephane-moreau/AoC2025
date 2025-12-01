package day1

import (
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadSequence(fn string) ([]int, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	readLines := strings.Split(string(content), "\r\n")

	input := make([]int, 0)
	for _, l := range readLines {
		d := 1
		if l[0] == 'L' {
			d = -1
		}
		v, err := strconv.Atoi(l[1:])
		if err != nil {
			return nil, err
		}
		input = append(input, v*d)
	}
	return input, nil
}

func countZeros(input []int) (int, int) {
	pos := 50
	c := 0
	pass := 0
	for _, i := range input {
		prevPos := pos
		pos += i
		if pos%100 == 0 {
			c++
		}
		if prevPos*pos < 0 {
			pass++
		}
		if pos < -100 || pos > 100 {
			// More than one full circle
			circles := pos / 100
			pos = pos - circles*100
			if pos%100 == 0 {
				if circles < 0 {
					circles++
				}
				if circles > 0 {
					circles--
				}
			}
			pass += int(math.Abs(float64(circles)))
		}
		pos = pos % 100
	}
	return c, pass
}

func TestDay1(t *testing.T) {
	light, err := loadSequence("light.txt")
	require.NoError(t, err)
	count, pass := countZeros(light)
	assert.Equal(t, 3, count)
	assert.Equal(t, 6, count+pass)

	// checking over 1 round the clock
	count, pass = countZeros(append(light, 1000))
	assert.Equal(t, 3, count)
	assert.Equal(t, 16, count+pass)

	large, err := loadSequence("large.txt")
	require.NoError(t, err)
	count, pass = countZeros(large)
	assert.Equal(t, 997, count)
	assert.Equal(t, 5978, count+pass)
}
