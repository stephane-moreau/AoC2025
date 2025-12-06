package day6

import (
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadValues(fn string) ([][]int64, [][]int64, []string, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, nil, err
	}

	lines := strings.Split(string(content), "\r\n")
	values := make([][]int64, 0)
	rtlValues := make([][]int64, 0)
	operators := make([]string, 0)
	ops := lines[len(lines)-1]

	curPos := 0
	end := 0
	for end < len(ops) {
		start := curPos
		end = start + 1
		for end < len(ops) && ops[end] == ' ' {
			end++
		}
		curPos = end
		operators = append(operators, strings.TrimSpace(ops[start:end]))

		vals := make([]int64, 0)
		for i := 0; i < len(lines)-1; i++ {
			v, err := strconv.ParseInt(strings.TrimSpace(lines[i][start:end]), 10, 64)
			if err != nil {
				return nil, nil, nil, err
			}
			vals = append(vals, v)
		}
		values = append(values, vals)

		vals = make([]int64, 0)
		for c := start; c < end; c++ {
			var sb strings.Builder
			for i := 0; i < len(lines)-1; i++ {
				sb.WriteByte(lines[i][c])
			}
			str := strings.TrimSpace(sb.String())
			if str == "" {
				continue
			}
			v, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return nil, nil, nil, err
			}
			vals = append(vals, v)
		}
		rtlValues = append(rtlValues, vals)
	}
	return values, rtlValues, operators, nil
}

func value(vals []int64, op string) int {
	var col int
	if op == "*" {
		col = 1
	}
	for _, v := range vals {
		switch op {
		case "+":
			col += int(v)
		case "*":
			col *= int(v)
		default:
			panic("invlid op")
		}
	}
	return col
}

func grandTotal(values [][]int64, ops []string) int {
	c := 0
	for i := range values {
		c += value(values[i], ops[i])
	}
	return c
}

func TestDay6Phase1(t *testing.T) {
	values, rtolValues, ops, err := loadValues("light.txt")
	require.NoError(t, err)
	assert.Equal(t, 4277556, grandTotal(values, ops))
	assert.Equal(t, 3263827, grandTotal(rtolValues, ops))

	values, rtolValues, ops, err = loadValues("large.txt")
	require.NoError(t, err)
	assert.Equal(t, 4693159084994, grandTotal(values, ops))
	assert.Equal(t, 11643736116335, grandTotal(rtolValues, ops))
}
