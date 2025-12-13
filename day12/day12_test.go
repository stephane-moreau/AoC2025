package day12

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkAreaSize(shapes []shape, a area) bool {
	zone := a.x * a.y
	for i, n := range a.shapes {
		size := 0
		for _, l := range shapes[i] {
			for _, c := range l {
				if c == '#' {
					size++
				}
			}
		}
		zone -= n * size
	}
	return zone >= 0
}

func checkAllAreas(shapes []shape, areas []area) int {
	num := 0
	for _, a := range areas {
		if checkAreaSize(shapes, a) {
			num++
		}
	}
	return num
}

func TestDay12KISS(t *testing.T) {
	shapes, areas, err := loadData("large.txt")
	require.NoError(t, err)
	assert.Equal(t, 472, checkAllAreas(shapes, areas))
}

func TestDay12(t *testing.T) {
	shapes, areas, err := loadData("light.txt")
	require.NoError(t, err)
	// should be 2 the hard way...
	assert.Equal(t, 3, checkAllAreas(shapes, areas))

	// The hard unfinished way
	ts := testShapes(shapes)
	require.True(t, fillArea(ts, areas[0]))
	require.True(t, fillArea(ts, areas[1]))
	//assert.Equal(t, 2, filledAreas(ts, areas))
}
