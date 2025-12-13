package day12

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

type shape []string

type area struct {
	x, y   int
	shapes []int
}

func loadData(fn string) ([]shape, []area, error) {
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
	shapes := make([]shape, 0, 5)
	areas := make([]area, 0, len(lines))

	for i := 0; i < len(lines); i++ {
		l := lines[i]
		if l == "" {
			continue
		}
		if l[len(l)-1] == ':' {
			s := []string{
				strings.ReplaceAll(lines[i+1], ".", "\x00"),
				strings.ReplaceAll(lines[i+2], ".", "\x00"),
				strings.ReplaceAll(lines[i+3], ".", "\x00"),
			}
			shapes = append(shapes, s)
			i += 3
			continue
		}
		var x, y, s1, s2, s3, s4, s5, s6 int
		_, err := fmt.Sscanf(l, "%dx%d: %d %d %d %d %d %d", &x, &y, &s1, &s2, &s3, &s4, &s5, &s6)
		if err != nil {
			return nil, nil, err
		}
		areas = append(areas, area{
			x, y,
			[]int{s1, s2, s3, s4, s5, s6},
		})
	}
	return shapes, areas, nil
}

func rotate(s shape) shape {
	newS := make(shape, len(s))
	// 123      741
	// 456  ->  852
	// 789      963
	newS[0] = string([]byte{s[2][0], s[1][0], s[0][0]})
	newS[1] = string([]byte{s[2][1], s[1][1], s[0][1]})
	newS[2] = string([]byte{s[2][2], s[1][2], s[0][2]})

	return newS
}

func applyShape(zone [][]byte, pX, pY int, s shape) bool {
	var i int
	for ; i < 9; i++ {
		y := i / 3
		x := i % 3
		if s[y][x] == 0 {
			continue
		}
		if pY+y >= len(zone) || pX+x >= len(zone[0]) || zone[pY+y][pX+x] != 0 {
			break
		}
		zone[pY+y][pX+x] = s[y][x]
	}

	if i != 9 {
		// reverse failed applicatio
		for i--; i >= 0; i-- {
			y := i / 3
			x := i % 3
			if s[y][x] == 0 {
				continue
			}
			if pY+y >= len(zone) || pX+x >= len(zone[0]) {
				continue
			}
			zone[pY+y][pX+x] = 0
		}
		return false
	}
	return true
}

func try(zone [][]byte, shapes []shape) int {
	for y := range zone {
		for x := range zone[y] {
			if zone[y][x] != 0 {
				continue
			}
			for i, s := range shapes {
				if applyShape(zone, x, y, s) {
					return i
				}
			}
		}
	}
	return -1
}

func fillArea(shapes [][]shape, a area) bool {
	zone := createZone(a.x, a.y)
	for i := range a.shapes {
		for range a.shapes[i] {

			if try(zone, shapes[i]) == -1 {
				return false
			}
		}
	}
	return true
}

func createZone(x, y int) [][]byte {
	zone := make([][]byte, y)
	for i := range zone {
		zone[i] = make([]byte, x)
	}
	return zone
}

func testShapes(shapes []shape) [][]shape {
	testShapes := make([][]shape, len(shapes))
	for i := range testShapes {
		s := shapes[i]
		testShapes[i] = append(testShapes[i], s)
		for r := 1; r <= 3; r++ {
			s = rotate(s)
			if !slices.ContainsFunc(testShapes[i], func(s1 shape) bool {
				return s1[0] == s[0] && s1[1] == s[1] && s1[2] == s[2]
			}) {
				testShapes[i] = append(testShapes[i], s)
			}
		}
	}
	return testShapes
}

func filledAreas(testShapes [][]shape, areas []area) int {
	c := 0
	for _, a := range areas {
		if fillArea(testShapes, a) {
			c++
		}
	}
	return c
}
