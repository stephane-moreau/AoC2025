package day2

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	light = "11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124"
	large = "9595822750-9596086139,1957-2424,88663-137581,48152-65638,12354817-12385558,435647-489419,518494-609540,2459-3699,646671-688518,195-245,295420-352048,346-514,8686839668-8686892985,51798991-51835611,8766267-8977105,2-17,967351-995831,6184891-6331321,6161577722-6161678622,912862710-913019953,6550936-6625232,4767634976-4767662856,2122995-2257010,1194-1754,779-1160,22-38,4961-6948,39-53,102-120,169741-245433,92902394-92956787,531-721,64-101,15596-20965,774184-943987,8395-11781,30178-47948,94338815-94398813"
)

func isValidCode(c int, useSplit bool) bool {
	s := strconv.Itoa(c)
	l := len(s)
	if useSplit {
		// simple double sequence
		return l%2 == 0 && s[:l/2] == s[l/2:]
	}
	for lenPart := 1; lenPart <= l/2; lenPart++ {
		isValid := true
		if l%lenPart != 0 {
			isValid = false
		}
		for start := 0; isValid && start+lenPart <= l-lenPart; start += lenPart {
			if s[start:start+lenPart] != s[start+lenPart:start+2*lenPart] {
				isValid = false
			}
		}
		if isValid {
			return true
		}
	}
	return false
}

type codeRange struct {
	start, end int
}

func parseInput(input string) []codeRange {
	var err error
	ranges := strings.Split(input, ",")
	codes := make([]codeRange, len(ranges))
	for i, r := range ranges {
		c := strings.Split(r, "-")
		var rg codeRange
		rg.start, err = strconv.Atoi(c[0])
		if err != nil {
			return nil
		}
		rg.end, err = strconv.Atoi(c[1])
		if err != nil {
			return nil
		}
		codes[i] = rg
	}
	return codes
}

func sumInvalidCodes(ranges []codeRange, useSplit bool) int {
	c := 0
	for _, r := range ranges {
		for code := r.start; code <= r.end; code++ {
			if isValidCode(code, useSplit) {
				c += code
			}
		}
	}
	return c
}

func TestValids(t *testing.T) {
	assert.True(t, isValidCode(11, false))
	assert.True(t, isValidCode(111, false))
	assert.True(t, isValidCode(565656, false))
	assert.True(t, isValidCode(824824824, false))
	assert.True(t, isValidCode(2121212121, false))
}

func TestDay2Pahes1(t *testing.T) {
	ranges := parseInput(light)
	require.NotNil(t, ranges)
	assert.Equal(t, 1227775554, sumInvalidCodes(ranges, true))
	assert.Equal(t, 4174379265, sumInvalidCodes(ranges, false))

	ranges = parseInput(large)
	require.NotNil(t, ranges)
	assert.Equal(t, 40398804950, sumInvalidCodes(ranges, true))
	assert.Equal(t, 65794984339, sumInvalidCodes(ranges, false))
}
