package day8

type point struct {
	x, y, z int
}

func sqDistance(p1, p2 point) int {
	return (p2.x-p1.x)*(p2.x-p1.x) +
		(p2.y-p1.y)*(p2.y-p1.y) +
		(p2.z-p1.z)*(p2.z-p1.z)
}

type connection struct {
	start, end point
	distance   int
}
