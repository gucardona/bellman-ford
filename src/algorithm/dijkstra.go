package algorithm

import (
	"container/heap"
	"errors"
	"math"

	"github.com/gucardona/bellman-ford/src/graph"
)

// Item para priority queue
type PQItem struct {
	Node     string
	Priority float64
	Index    int
}

// PriorityQueue implementa heap.Interface
type PriorityQueue []*PQItem

func (pq *PriorityQueue) Len() int {
	return len(*pq)
}
func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].Priority < (*pq)[j].Priority
}
func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].Index = i
	(*pq)[j].Index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PQItem)
	item.Index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// Dijkstra retorna dist e predecessor.
// Observação: Dijkstra não é correto com arestas de peso negativo.
// Aqui, não fazemos checagem ativa — usuário deve fornecer grafos sem pesos negativos
func Dijkstra(g *graph.Graph, source string) (map[string]float64, map[string]string, error) {
	dist := make(map[string]float64)
	prev := make(map[string]string)

	for node := range g.Nodes {
		dist[node] = math.Inf(1)
		prev[node] = ""
	}
	if _, ok := g.Nodes[source]; !ok {
		return nil, nil, errors.New("source node not in graph")
	}
	dist[source] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &PQItem{Node: source, Priority: 0})

	visited := make(map[string]bool)

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*PQItem)
		u := item.Node
		if visited[u] {
			continue
		}
		visited[u] = true
		for _, e := range g.Adj[u] {
			v := e.To
			alt := dist[u] + e.W
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u
				heap.Push(pq, &PQItem{Node: v, Priority: alt})
			}
		}
	}
	return dist, prev, nil
}
