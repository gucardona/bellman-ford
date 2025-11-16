package algorithm

import (
	"errors"
	"math"

	"github.com/gucardona/bellman-ford/src/graph"
)

// BellmanFord retorna (dist, prev, hasNegCycle)
// dist infinito aponta que não alcançável.
// Se detectar ciclo negativo retorna hasNegCycle = true
func BellmanFord(g *graph.Graph, source string) (map[string]float64, map[string]string, bool, error) {
	dist := make(map[string]float64)
	prev := make(map[string]string)

	for node := range g.Nodes {
		dist[node] = math.Inf(1)
		prev[node] = ""
	}
	if _, ok := g.Nodes[source]; !ok {
		return nil, nil, false, errors.New("source node not in graph")
	}
	dist[source] = 0

	n := len(g.Nodes)
	// relax edges n-1 vezes
	for i := 0; i < n-1; i++ {
		changed := false
		for _, e := range g.Edges {
			if dist[e.From]+e.W < dist[e.To] {
				dist[e.To] = dist[e.From] + e.W
				prev[e.To] = e.From
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	// checar ciclo negativo
	for _, e := range g.Edges {
		if dist[e.From]+e.W < dist[e.To] {
			return dist, prev, true, nil
		}
	}
	return dist, prev, false, nil
}
