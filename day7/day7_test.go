package day7

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type grid [][]byte

type point struct {
	x, y int
}

func loadGrid(fn string) (grid, error) {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\r\n")
	g := make(grid, len(lines))
	for i := range lines {
		g[i] = []byte(lines[i])
	}
	return g, nil
}

func countSplits(g grid) (int, int) {
	c := 0
	var start point
	for i := range g[0] {
		if g[0][i] == 'S' {
			start = point{i, 0}
		}
	}
	beams := map[point]int{
		start: 1,
	}
	for y := range len(g) - 1 {
		newBeams := make(map[point]int)
		for x := range g[y] {
			if b := beams[point{x, y}]; b != 0 {
				if g[y][x] == '^' {
					c++
					newBeams[point{x - 1, y + 1}] += b
					newBeams[point{x + 1, y + 1}] += b
				} else {
					newBeams[point{x, y + 1}] += b
				}
			}
		}
		beams = newBeams
	}
	b := 0
	for _, v := range beams {
		b += v
	}
	return c, b
}

func TestDay7(t *testing.T) {
	g, err := loadGrid("light.txt")
	require.NoError(t, err)
	splits, paths := countSplits(g)
	assert.Equal(t, 21, splits)
	assert.Equal(t, 40, paths)

	g, err = loadGrid("large.txt")
	require.NoError(t, err)
	splits, paths = countSplits(g)
	assert.Equal(t, 1656, splits)
	assert.Equal(t, 76624086587804, paths)
}
