package day9

type point struct {
	x, y int
}

func length(p1, p2 int) int {
	if p1 >= p2 {
		return p1 - p2 + 1
	}
	return p2 - p1 + 1
}

func min(x1, x2 int) int {
	if x1 < x2 {
		return x1
	}
	return x2
}

func max(x1, x2 int) int {
	if x1 > x2 {
		return x1
	}
	return x2
}

func area(p1, p2 point) int {
	return length(p2.x, p1.x) * length(p2.y, p1.y)
}
