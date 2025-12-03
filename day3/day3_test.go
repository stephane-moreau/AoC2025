package day3

import (
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadSequence(fn string) ([]string, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(content), "\r\n"), nil
}

func sumMaxNumbers(lines []string, numDigits int) int {
	var sum int64
	for _, l := range lines {
		digits := make([]byte, numDigits)
		for i := range l {
			for j := range len(digits) {
				if i < len(l)-len(digits)+1+j && digits[j] < l[i] {
					digits[j] = l[i]
					for k := j + 1; k < len(digits); k++ {
						digits[k] = 0
					}
					break
				}
			}
		}
		v, err := strconv.ParseInt(string(digits), 10, 64)
		if err != nil {
			return 0
		}
		sum += v
	}
	return int(sum)
}

func TestDay3(t *testing.T) {
	lines, err := loadSequence("light.txt")
	require.NoError(t, err)
	assert.Equal(t, 357, sumMaxNumbers(lines, 2))
	assert.Equal(t, 3121910778619, sumMaxNumbers(lines, 12))

	lines, err = loadSequence("large.txt")
	require.NoError(t, err)
	assert.Equal(t, 17085, sumMaxNumbers(lines, 2))
	assert.Equal(t, 169408143086082, sumMaxNumbers(lines, 12))
}
