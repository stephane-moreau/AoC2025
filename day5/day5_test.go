package day5

import (
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type idRange struct {
	start, end int64
}

func loadData(fn string) ([]idRange, []int64, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}

	lines := strings.Split(string(content), "\r\n")
	ranges := make([]idRange, 0, len(lines)/2)
	ids := make([]int64, 0, len(lines)/2)
	for _, l := range lines {
		if l == "" {
			continue
		}
		vals := strings.Split(l, "-")
		id1, err := strconv.ParseInt(vals[0], 10, 64)
		if err != nil {
			return nil, nil, err
		}
		if len(vals) == 2 {
			id2, err := strconv.ParseInt(vals[1], 10, 64)
			if err != nil {
				return nil, nil, err
			}
			ranges = append(ranges, idRange{
				id1, id2,
			})
			continue
		}
		ids = append(ids, id1)
	}
	sort.SliceStable(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})
	return ranges, ids, nil
}

func countValid(ids []int64, ranges []idRange) int {
	c := 0
foundId:
	for _, id := range ids {
		for _, r := range ranges {
			if id >= r.start && id <= r.end {
				c++
				continue foundId
			}
		}
	}
	return c
}

func countValidIds(ranges []idRange) int {
	c := 0
	id := int64(0)
	for i, r := range ranges {
		if r.end < id {
			continue
		}
		if i == 0 {
			id = r.start
		} else if id < r.start {
			id = r.start
		} else {
			c--
		}
		c += int(r.end - id + 1)
		id = r.end
	}
	return c
}

func TestDay5(t *testing.T) {
	ranges, ids, err := loadData("light.txt")
	require.NoError(t, err)
	assert.Equal(t, 3, countValid(ids, ranges))
	assert.Equal(t, 14, countValidIds(ranges))

	ranges, ids, err = loadData("large.txt")
	require.NoError(t, err)
	assert.Equal(t, 635, countValid(ids, ranges))
	assert.Equal(t, 369761800782619, countValidIds(ranges))
}
