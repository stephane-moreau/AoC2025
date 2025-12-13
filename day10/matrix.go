package day10

import (
	"math"
	"slices"
	"sort"
)

func abs(f matrixCell) matrixCell {
	if f >= 0 {
		return f
	}
	return -f
}

func sign(f matrixCell) matrixCell {
	if f >= 0 {
		return 1
	}
	return -1
}

func sum(res []int) int {
	c := 0
	for _, r := range res {
		if r < 0 {
			continue
		}
		c += r
	}
	return c
}

type matrixCell = int
type matrixRow = []matrixCell
type matrix = [][]matrixCell

func clone(mtrx matrix) matrix {
	newMtrx := make(matrix, len(mtrx))
	for i := range newMtrx {
		newMtrx[i] = slices.Clone(mtrx[i])
	}
	return newMtrx
}

func findNonSolvedIndexes(mtrx matrix) []int {
	varIndex := make([]int, 0)
	for c := 0; c < len(mtrx[0])-1; c++ {
		countNonNil := 0
		for j := 0; j < len(mtrx); j++ {
			if mtrx[j][c] != 0 {
				countNonNil++
			}
		}
		if countNonNil > 1 {
			varIndex = append(varIndex, c)
		}
	}
	return varIndex
}

func removeEmptyLines(mtrx matrix) (matrix, bool) {
	var removed bool
	for m := len(mtrx) - 1; m >= 0; m-- {
		fullOfZero := true
		for i := range mtrx[m] {
			if mtrx[m][i] != 0 {
				fullOfZero = false
				break
			}
		}
		if fullOfZero {
			removed = true
			mtrx = append(mtrx[:m], mtrx[m+1:]...)
		}
	}
	return mtrx, removed
}

func solveProportionalLines(mtrx matrix, res []int) (matrix, []int) {
	resNdx := len(mtrx[0]) - 1
	removed := true
	for removed {
		removed = false
		for m := range mtrx {
			propIndex := -1
			for i := range len(mtrx[m]) - 1 {
				if mtrx[m][i] != 0 {
					if propIndex == -1 {
						propIndex = i
					} else {
						propIndex = -1
						break
					}
				}
			}
			if propIndex != -1 {
				res[propIndex] = int(mtrx[m][resNdx] / mtrx[m][propIndex])
				if res[propIndex] < 0 {
					return nil, nil
				}
				for pm := len(mtrx) - 1; pm >= 0; pm-- {
					mtrx[pm][resNdx] = mtrx[pm][resNdx] - mtrx[pm][propIndex]*matrixCell(res[propIndex])
					mtrx[pm][propIndex] = 0
				}
			}
		}
		mtrx, removed = removeEmptyLines(mtrx)
	}
	return mtrx, res
}

func apply(mtrx matrix, res []int) (matrix, []int, bool) {
	resNdx := len(mtrx[0]) - 1
	newmtrx := slices.Clone(mtrx)
	newRes := slices.Clone(res)
	for m := range newmtrx {
		newmtrx[m] = slices.Clone(mtrx[m])
		for i, r := range newRes {
			if r == -1 {
				continue
			}
			if newmtrx[m][i] != 0 {
				newmtrx[m][resNdx] = newmtrx[m][resNdx] - r*newmtrx[m][i]
				newmtrx[m][i] = 0
			}
		}
	}
	newmtrx, newRes = solveProportionalLines(newmtrx, newRes)
	return newmtrx, newRes, len(newmtrx) == 0 && newRes != nil
}

func triangularize(mtrx matrix, knownRes []int) (matrix, []int, []int) {
	// Move as close as possible to upper triangular
	for l := range mtrx {
		sort.SliceStable(mtrx[l:], func(i, j int) bool {
			return abs(mtrx[i+l][l]) > abs(mtrx[j+l][l])
		})

		for o := l + 1; o < len(mtrx); o++ {
			if mtrx[o][l] == 0 {
				break
			}
			diag := mtrx[l][l]
			mult := mtrx[o][l]
			for c := l; c < len(mtrx[o]); c++ {
				mtrx[o][c] = mtrx[l][c]*mult - mtrx[o][c]*diag
			}
		}
	}
	if knownRes == nil {
		// top down diagonalization
		for l := range mtrx {

			for firstNonZero := l + 1; firstNonZero < len(mtrx[l])-1; firstNonZero++ {
				if mtrx[l][firstNonZero] == 0 {
					continue
				}

				lineNdx := l + 1
			cancelLine:
				for ; lineNdx < len(mtrx); lineNdx++ {
					if mtrx[lineNdx][firstNonZero] == 0 {
						continue
					}

					for j := 0; j < firstNonZero; j++ {
						if mtrx[lineNdx][j] != 0 {
							continue cancelLine
						}
					}
					break
				}
				if lineNdx == len(mtrx) {
					continue
				}

				// multiply line l by line[firstNonZero] and remove mtrx[l][firstNonZero]*line[o]
				d := mtrx[l][firstNonZero]
				for o := 0; o < len(mtrx[l]); o++ {
					mtrx[l][o] = mtrx[l][o]*mtrx[lineNdx][firstNonZero] - d*mtrx[lineNdx][o]
				}
			}
		}
	}
notEmpty:
	for i := len(mtrx) - 1; i >= 0; i-- {
		for _, c := range mtrx[i] {
			if c != 0 {
				break notEmpty
			}
		}
		mtrx = mtrx[:len(mtrx)-1]
	}

	varIndex := findNonSolvedIndexes(mtrx)
	if len(varIndex) == 0 {
		return mtrx, nil, nil
	}
	var res []int
	if knownRes != nil {
		res = slices.Clone(knownRes)
	} else {
		res = make([]int, len(mtrx[0])-1)
		for i := range res {
			res[i] = -1
		}
	}
	// solve all proportionnal values
	mtrx, res = solveProportionalLines(mtrx, res)
	return mtrx, res, varIndex
}

const MAX_TEST = 200

func iterativeSolver(mtrx matrix, res []int, varIndexes []int, minRes *int, maxSearch int) (int, bool) {
	resNdx := len(mtrx[0]) - 1
	curSum := sum(res)
	var localMin int
	if minRes == nil {
		minRes = &localMin
		*minRes = math.MaxInt
	}
	if len(varIndexes) == 0 {
		newMtrx, res, final := apply(mtrx, res)
		if final {
			return sum(res), true
		}
		if res == nil {
			return 0, false
		}
		mtrx, newRes, _v := triangularize(newMtrx, res)
		if len(_v) != 0 || newRes == nil {
			return 0, false
		}
		for m := len(mtrx) - 1; m >= 0; m-- {
			r := 0
			for i := 0; i < resNdx; i++ {
				r += newRes[i] * int(mtrx[m][i])
			}
			if r != int(mtrx[m][resNdx]) {
				return 0, false
			}
		}
		return sum(newRes), true
	}

	var found1Sol bool
	for resTest := 0; resTest <= maxSearch; resTest++ {
		if resTest+curSum > *minRes {
			return *minRes, found1Sol
		}
		res[varIndexes[0]] = resTest
		r, ok := iterativeSolver(mtrx, res, varIndexes[1:], minRes, maxSearch)
		if ok {
			found1Sol = true
			if *minRes > r {
				*minRes = r
			}
		}
		res[varIndexes[0]] = -1
	}
	return *minRes, found1Sol
}

func dynSolve(mtrx matrix, knownRes []int, varIndex []int) int {
	for v := len(varIndex) - 1; v >= 0; v-- {
		if knownRes[varIndex[v]] != -1 {
			varIndex = append(varIndex[:v], varIndex[v+1:]...)
		}
	}
	maxSearch := MAX_TEST / len(varIndex)
	for ; maxSearch <= MAX_TEST; maxSearch += 25 {
		s, ok := iterativeSolver(mtrx, knownRes, varIndex, nil, maxSearch)
		if ok {
			return s
		}
	}
	return 0
}
