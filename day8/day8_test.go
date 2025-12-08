package day8

import (
	"maps"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildConnections(input []point) []connection {
	g := make([]connection, 0, len(input)*(len(input)-1)/2)
	for s := range input {
		for e := s + 1; e < len(input); e++ {
			g = append(g, connection{
				input[s],
				input[e],
				sqDistance(input[s], input[e]),
			})
		}
	}
	sort.SliceStable(g, func(i, j int) bool {
		return g[i].distance < g[j].distance
	})
	return g
}

type circuit map[point]bool

var checkPoints = map[point]bool{
	{84314, 61210, 54625}: true,
	{84358, 62436, 51394}: true,
	{85096, 62446, 58342}: true,
}

func mergeCircuits(circuits []circuit) []circuit {
nextCircuit:
	for i := 0; i < len(circuits); i++ {
		for p := range circuits[i] {
			for j := i + 1; j < len(circuits); j++ {
				if circuits[j][p] {
					for k := range circuits[j] {
						circuits[i][k] = true
					}
					circuits = append(circuits[:j], circuits[j+1:]...)
					i--
					continue nextCircuit
				}
			}
		}
	}
	sort.SliceStable(circuits, func(i, j int) bool {
		return len(circuits[i]) > len(circuits[j])
	})
	return circuits
}

func findCircuits(cnx []connection, snapshot int, num int) ([]circuit, connection) {
	circuits := make([]circuit, 0)
	var crcSnap []circuit
nextEdge:
	for l, c := range cnx {
		if l == snapshot {
			crcSnap = make([]circuit, len(circuits))
			for i := range circuits {
				crcSnap[i] = maps.Clone(circuits[i])
			}
			crcSnap = mergeCircuits(crcSnap)
		}
		if len(circuits) == 1 && len(circuits[0]) == num {
			return crcSnap, cnx[l-1]
		}
		for i := range circuits {
			if circuits[i][c.start] {
				if circuits[i][c.end] {
					continue nextEdge
				}
				circuits[i][c.end] = true
				circuits = mergeCircuits(circuits)
				continue nextEdge
			}
			if circuits[i][c.end] {
				if circuits[i][c.start] {
					continue nextEdge
				}
				circuits[i][c.start] = true
				continue nextEdge
			}
		}
		circuits = append(circuits, circuit{
			c.start: true,
			c.end:   true,
		})
	}
	// merge existing circuits
	return crcSnap, connection{}
}

func TestDay8(t *testing.T) {
	g := buildConnections(light)
	require.NotNil(t, g)
	crc, last := findCircuits(g, 10, len(light))
	assert.Equal(t, 40, len(crc[0])*len(crc[1])*len(crc[2]))
	assert.Equal(t, 25272, last.start.x*last.end.x)
	g = buildConnections(large)
	require.NotNil(t, g)
	crc, last = findCircuits(g, 1000, len(large))
	assert.Equal(t, 67488, len(crc[0])*len(crc[1])*len(crc[2]))
	assert.Equal(t, 3767453340, last.start.x*last.end.x)
}
