package day11

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type node struct {
	src     string
	targets []string
}

type graph map[string]node

func loadGraph(fn string) (graph, error) {
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
	graph := make(graph, len(lines))
	for _, l := range lines {
		st := strings.Split(l, ": ")
		tgt := strings.Split(st[1], " ")
		graph[st[0]] = node{
			st[0],
			tgt,
		}
	}
	return graph, nil
}

const (
	none = iota
	fft
	dac
	all
)

type cache map[string]map[int]int
type visits map[string]bool

func pathCount(g graph, cur string, visited visits, cch cache, remainingVisits int) int {
	c := 0
	if cur == "out" {
		if cch == nil {
			return 1
		}
		if visited["fft"] && visited["dac"] {
			return 1
		}
		return 0
	}
	known, exist := cch[cur]
	if exist {
		c, exist := known[remainingVisits]
		if exist {
			return c
		}
	}
	src := g[cur]
	visited[cur] = true
	for _, tgt := range src.targets {
		if visited[tgt] {
			continue
		}
		nextRemaining := remainingVisits
		if tgt == "fft" {
			nextRemaining |= fft
		}
		if tgt == "dac" {
			nextRemaining |= dac
		}
		c += pathCount(g, tgt, visited, cch, nextRemaining)
	}
	delete(visited, cur)
	if cch != nil {
		cchCur, cached := cch[cur]
		if !cached {
			cchCur = map[int]int{}
			cch[cur] = cchCur
		}
		cchCur[remainingVisits] = c
	}
	return c
}

func TestDay11(t *testing.T) {
	g, err := loadGraph("light.txt")
	require.NoError(t, err)
	assert.Equal(t, 5, pathCount(g, "you", visits{}, nil, none))

	g, err = loadGraph("fftdac.txt")
	require.NoError(t, err)
	assert.Equal(t, 2, pathCount(g, "svr", visits{}, cache{}, none))

	g, err = loadGraph("large.txt")
	require.NoError(t, err)
	assert.Equal(t, 788, pathCount(g, "you", visits{}, nil, none))
	assert.Equal(t, 316291887968000, pathCount(g, "svr", visits{}, cache{}, none))
}
