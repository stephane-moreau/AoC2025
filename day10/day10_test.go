package day10

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type buttonSet [][]int

type joltage []int

type light struct {
	target  string
	buttons buttonSet
	vals    joltage
}

func loadLights(fn string) ([]light, error) {
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
	lights := make([]light, len(lines))

	for i, l := range lines {
		values := strings.Split(l, " ")
		lt := light{
			target:  values[0][1 : len(values[0])-1],
			buttons: make(buttonSet, len(values)-2),
			vals:    make(joltage, 0, len(values[len(values)-1])),
		}
		for j := 1; j < len(values)-1; j++ {
			btn := strings.Split(values[j][1:len(values[j])-1], ",")
			for _, b := range btn {
				v, err := strconv.ParseInt(b, 10, 64)
				if err != nil {
					return nil, err
				}
				lt.buttons[j-1] = append(lt.buttons[j-1], int(v))
			}
		}
		ndx := len(values) - 1
		jlt := strings.Split(values[ndx][1:len(values[ndx])-1], ",")
		for _, b := range jlt {
			v, err := strconv.ParseInt(b, 10, 64)
			if err != nil {
				return nil, err
			}
			lt.vals = append(lt.vals, int(v))
		}
		lights[i] = lt
	}
	return lights, nil
}

func minimalButtonPress(l light, joltage bool) int {
	s := make([]byte, len(l.target))
	target := []byte(l.target)
	if !joltage {
		for i := range s {
			s[i] = '.'
		}
	} else {
		for i := range target {
			target[i] = byte(l.vals[i])
		}
	}
	state := map[string]int{
		string(s): 0,
	}
	for state[string(target)] == 0 {
		newState := make(map[string]int)
		for s, pressed := range state {
		buttonsLoop:
			for _, btns := range l.buttons {
				values := []byte(s)
				for _, b := range btns {
					if joltage {
						values[b]++
						if values[b] > target[b] {
							continue buttonsLoop
						}
						continue
					}
					if values[b] == '.' {
						values[b] = '#'
					} else {
						values[b] = '.'
					}
				}
				newState[string(values)] = pressed + 1
			}
			state = newState
		}
	}
	return state[string(target)]
}

func computeButtonPress(lights []light) int {
	c := 0
	for _, l := range lights {
		c += minimalButtonPress(l, false)
	}
	return c
}

func computeVoltage(lights []light) int {
	c := 0
	for _, l := range lights {
		c += minimalButtonPress(l, true)
	}
	return c
}

func solveButtonPress(lgt light) int {
	var mtrx = make(matrix, len(lgt.vals))
	for i := range mtrx {
		mtrx[i] = make(matrixRow, len(lgt.buttons)+1)
		mtrx[i][len(lgt.buttons)] = lgt.vals[i]
	}

	for i, btns := range lgt.buttons {
		for _, b := range btns {
			mtrx[b][i] = 1
		}
	}

	mtrx, preRes, ts := triangularize(mtrx, nil)

	if ts != 0 {
		// Solve all single var
		s := dynSolve(mtrx, preRes, ts)
		if s == 0 {
			fmt.Printf("%v\n", lgt)
			return 0
		}
		return s
	}
	res := make([]int, len(lgt.buttons))
	resNdx := len(mtrx[0]) - 1
	for l := len(mtrx) - 1; l >= 0; l-- {
		foundRes := l

		if mtrx[l][resNdx] != 0 {
			p := float64(mtrx[l][resNdx]) / float64(mtrx[l][foundRes])
			if p != float64(int(p)) || p < 0 {
				panic("invalid solution")
			}
			res[foundRes] = int(p)
		}
		for m := l; m >= 0; m-- {
			mtrx[m][resNdx] = mtrx[m][resNdx] - res[foundRes]*mtrx[m][foundRes]
			mtrx[m][foundRes] = 0
		}
	}
	return sum(res)
}

func solveVoltage(lights []light) int {
	c := 0
	for i, l := range lights {
		if i%10 == 9 {
			print("#")
		}
		if i%10 == 49 {
			println("")
		}
		s := solveButtonPress(l)
		c += s
	}
	return c
}

func sumHighest(lights []light) int {
	c := 0
	for _, l := range lights {
		m := 0
		for _, v := range l.vals {
			if m < v {
				m = v
			}
		}
		c += m
	}
	return c
}

func TestOtherLineDiag(t *testing.T) {
	l := light{
		"#.####...#",
		[][]int{{2, 5, 6}, {0, 1, 3, 4, 6, 7}, {0, 1, 2, 3, 4, 5, 7, 8}, {3, 6}, {3, 4, 6, 7}, {0, 2, 3, 5, 7, 8, 9}, {0, 1, 2, 3, 4, 5, 7}, {1, 2, 3, 4, 5, 9}},
		[]int{37, 37, 56, 66, 48, 56, 48, 48, 10, 20},
	}
	assert.Equal(t, 85, solveButtonPress(l))
}

func TestOneLineIterSolver(t *testing.T) {
	l := light{
		target:  "..#.##",
		buttons: [][]int{{0, 1, 3, 4}, {0, 3, 4}, {0, 5}, {0, 1, 2}, {3, 5}, {0, 2, 3, 4}, {2, 3}},
		vals:    []int{58, 27, 37, 57, 37, 24},
	}
	assert.Equal(t, 78, solveButtonPress(l))
}

func TestOneLineIterSolver2(t *testing.T) {
	l := light{
		target: "##.#.###..",
		buttons: [][]int{
			{1}, {0, 2, 9}, {0, 2, 4}, {0, 2, 4, 5, 8, 9}, {0, 1, 3, 4, 5, 6, 7, 8}, {1, 6, 8}, {0, 4, 5}, {2, 4, 6, 8, 9},
			{0, 1, 7, 8, 9}, {1, 4, 8}, {0, 1, 3, 4, 7, 8}, {2, 3, 5, 6, 7, 9},
		},
		vals: []int{44, 30, 37, 16, 49, 27, 23, 22, 40, 36},
	}
	assert.Equal(t, 97, solveButtonPress(l))
}

func TestDay10(t *testing.T) {
	lights, err := loadLights("light.txt")
	require.NoError(t, err)
	assert.Equal(t, 7, computeButtonPress(lights))
	assert.Equal(t, 33, computeVoltage(lights))
	assert.Equal(t, 10, solveButtonPress(lights[0]))
	assert.Equal(t, 12, solveButtonPress(lights[1]))
	assert.Equal(t, 11, solveButtonPress(lights[2]))

	lights, err = loadLights("large.txt")
	require.NoError(t, err)
	assert.Equal(t, 538, computeButtonPress(lights))
	assert.Equal(t, 17716, sumHighest(lights)) // to low (obviously)
	assert.Equal(t, 20298, solveVoltage(lights))
}
