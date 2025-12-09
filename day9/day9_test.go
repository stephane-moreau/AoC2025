package day9

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func largestArea(tiles []point) int {
	res := 0
	for i := range tiles {
		for j := i + 1; j < len(tiles); j++ {
			if area(tiles[i], tiles[j]) > res {
				res = area(tiles[i], tiles[j])
			}
		}
	}
	return res
}

func isIn(p point, corner1, corner2 point) bool {
	topLeft := point{min(corner1.x, corner2.x), min(corner1.y, corner2.y)}
	bottomRight := point{max(corner1.x, corner2.x), max(corner1.y, corner2.y)}
	if p.x >= topLeft.x && p.x <= bottomRight.x && p.y >= topLeft.y && p.y <= bottomRight.y {
		return true
	}
	return false
}

type pointSet map[point]bool

func largestInnerArea(tiles []point, extPoints pointSet) int {
	res := 0
	for i := range tiles {
	nextTile:
		for j := i + 1; j < len(tiles); j++ {
			for p := range extPoints {
				if isIn(p, tiles[i], tiles[j]) {
					continue nextTile
				}
			}
			newArea := area(tiles[i], tiles[j])
			if newArea > res {
				res = newArea
			}
		}
	}
	return res
}

func dump(topLeft, bottomRight, cur point, ps, ext pointSet) {
	if bottomRight.x-topLeft.x < 50 {
		println("")
		for y := topLeft.x - 1; y < bottomRight.y+2; y++ {
			s := ""
			for x := topLeft.x - 1; x < bottomRight.x+2; x++ {
				pt := point{x, y}
				if pt == cur {
					s += "O"
				} else if ext[pt] {
					s += "X"
				} else if ps[pt] {
					s += "#"
				} else {
					s += "."
				}
			}
			println(s)
		}
	}
}

var cardDirections = []point{point{0, -1}, point{0, 1}, point{-1, 0}, point{1, 0}}

// rotations:
// x' = cos(θ) * x - sin(θ) * y
// y' = sin(θ) * x + cos(θ) * y
//
// rotation clockwise
func rotCW(vector point) point {
	return point{-vector.y, vector.x}
}

// rotation counter clockwise
func rotCCW(vector point) point {
	return point{vector.y, -vector.x}
}

func externalPoints(tiles []point) pointSet {
	ps := make(pointSet, 2*len(tiles))
	var bottomRight, topLeft point
	for i := range len(tiles) {
		tile := tiles[i]
		ps[tile] = true
		var nextTile point
		if i < len(tiles)-1 {
			nextTile = tiles[i+1]
		} else {
			nextTile = tiles[0]
		}
		if tile.x == nextTile.x {
			yMin := min(tile.y, nextTile.y)
			yMax := max(tile.y, nextTile.y)
			for y := yMin + 1; y < yMax; y++ {
				ps[point{tile.x, y}] = true
			}
		} else if tile.y == nextTile.y {
			xMin := min(tile.x, nextTile.x)
			xMax := max(tile.x, nextTile.x)
			for x := xMin + 1; x < xMax; x++ {
				ps[point{x, tile.y}] = true
			}
		} else {
			panic("diagonal ?")
		}
		bottomRight.x = max(bottomRight.x, tile.x)
		bottomRight.y = max(bottomRight.y, tile.y)
		topLeft.x = min(topLeft.x, tile.x)
		topLeft.y = min(topLeft.y, tile.y)
	}
	dump(topLeft, bottomRight, point{}, ps, nil)
	ext := make(pointSet)
	var extStart []point
	for _, tile := range tiles {
		extPoints := []point{}
		for _, cardinal := range cardDirections {
			ptCard := point{tile.x + cardinal.x, tile.y + cardinal.y}
			numFree := 0
		cardinalChecks:
			for _, sc := range cardDirections {
				testPoint := ptCard
				if ps[testPoint] {
					continue
				}
				move := 0
				for testPoint.x <= bottomRight.x && testPoint.x >= topLeft.x && testPoint.y <= bottomRight.y && testPoint.y >= topLeft.y {
					move++
					testPoint = point{ptCard.x + move*sc.x, ptCard.y + move*sc.y}
					if ps[testPoint] {
						continue cardinalChecks
					}
				}
				numFree++
			}
			if numFree != 0 {
				extPoints = append(extPoints, ptCard)
			}
		}
		if len(extPoints) > 1 {
			extStart = extPoints
			break
		}
	}
	// ext is cornder to start external perimeter
	// free direction is a perpendicular of a perimeter start
	for _, p := range extStart {
		ext[p] = true
	}
	ptCorner1 := point{extStart[0].x, extStart[1].y}
	ptCorner2 := point{extStart[1].x, extStart[0].y}
	var ptCornerStart, ptCorner point
	if ps[ptCorner1] {
		ext[ptCorner2] = true
		ptCornerStart = ptCorner2
		ptCorner = ptCorner1
	} else if ps[ptCorner2] {
		ext[ptCorner1] = true
		ptCornerStart = ptCorner1
		ptCorner = ptCorner2
	} else {
		panic("what's this corner")
	}
	var moving, zoneDir point
	for _, dir := range cardDirections {
		next := point{ptCornerStart.x + dir.x, ptCornerStart.y + dir.y}
		if ext[next] {
			moving = dir
			zoneDir = point{ptCorner.x - next.x, ptCorner.y - next.y}
			break
		}
	}
	next := point{ptCornerStart.x + moving.x, ptCornerStart.y + moving.y}
	for next != ptCornerStart {
		next = point{next.x + moving.x, next.y + moving.y}
		dump(topLeft, bottomRight, next, ps, ext)
		neighbor := point{next.x + zoneDir.x, next.y + zoneDir.y}
		var newMoving point
		if ps[next] {
			next = point{next.x - moving.x, next.y - moving.y}
			if zoneDir == rotCW(moving) {
				newMoving = rotCCW(moving)
			}
			if zoneDir == rotCCW(moving) {
				newMoving = rotCW(moving)
			}
		}
		ext[next] = true
		if !ps[neighbor] {
			newMoving = zoneDir
		}
		if newMoving.x != 0 || newMoving.y != 0 {
			if newMoving == rotCW(moving) {
				zoneDir = rotCW(zoneDir)
			}
			if newMoving == rotCCW(moving) {
				zoneDir = rotCCW(zoneDir)
			}
			moving = newMoving
		}
	}
	return ext
}

func TestDay9(t *testing.T) {
	assert.Equal(t, 50, largestArea(light))
	assert.Equal(t, 4771508457, largestArea(large))

	assert.Equal(t, 24, largestInnerArea(light, externalPoints(light)))
	assert.Equal(t, 25, largestInnerArea(lightH, externalPoints(lightH)))
	// 4694419923 too high
	assert.Equal(t, 1539809693, largestInnerArea(large, externalPoints(large)))
}
