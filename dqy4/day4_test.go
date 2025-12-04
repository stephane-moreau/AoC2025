package day4

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type grid [][]byte

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

func (g grid) Value(x, y int) byte {
	if x < 0 || y < 0 || y >= len(g) || x >= len(g[y]) {
		return 0
	}
	return g[y][x]
}

func (g grid) numAdjValues(x, y int, v byte) int {
	c := 0
	if g.Value(x-1, y-1) == v {
		c++
	}
	if g.Value(x-1, y) == v {
		c++
	}
	if g.Value(x-1, y+1) == v {
		c++
	}
	if g.Value(x, y-1) == v {
		c++
	}
	if g.Value(x, y+1) == v {
		c++
	}
	if g.Value(x+1, y-1) == v {
		c++
	}
	if g.Value(x+1, y) == v {
		c++
	}
	if g.Value(x+1, y+1) == v {
		c++
	}
	return c
}

type point struct {
	x, y int
}

func accessibleRolls(g grid) []point {
	pts := make([]point, 0, 5*len(g))
	for y := range len(g) {
		for x := range len(g[y]) {
			if g.Value(x, y) == '@' {
				if g.numAdjValues(x, y, '@') < 4 {
					pts = append(pts, point{x, y})
				}
			}
		}
	}
	return pts
}

func removeAccessibleRolls(g grid) int {
	acc := 0
	removed := 1
	for removed != 0 {
		points := accessibleRolls(g)
		for _, p := range points {
			g[p.y][p.x] = 'x'
		}
		removed = len(points)
		acc += removed
	}
	return acc
}

func TestDay4Phase1(t *testing.T) {
	g, err := loadGrid("light.txt")
	require.NoError(t, err)
	assert.Equal(t, 13, len(accessibleRolls(g)))
	assert.Equal(t, 43, removeAccessibleRolls(g))

	g, err = loadGrid("large.txt")
	require.NoError(t, err)
	assert.Equal(t, 1480, len(accessibleRolls(g)))
	assert.Equal(t, 8899, removeAccessibleRolls(g))
}
