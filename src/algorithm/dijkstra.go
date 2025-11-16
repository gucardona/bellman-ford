package algorithm

import (
	"container/heap"
	"errors"
	"fmt"
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

// Len retorna o tamanho da heap.
func (pq *PriorityQueue) Len() int {
	return len(*pq)
}

// Less compara a prioridade de dois itens.
func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].Priority < (*pq)[j].Priority
}

// Swap troca dois itens na heap.
func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].Index = i
	(*pq)[j].Index = j
}

// Push adiciona um item à heap.
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PQItem)
	item.Index = n
	*pq = append(*pq, item)
}

// Pop remove e retorna o item de menor prioridade da heap.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// Dijkstra executa o algoritmo e retorna uma lista de 'snapshots' (frames)
// da execução, e um erro.
func Dijkstra(g *graph.Graph, source string) ([]graph.Snapshot, error) {
	dist := make(map[string]float64)
	prev := make(map[string]string)
	visited := make(map[string]bool)

	var snapshots []graph.Snapshot // Lista de frames

	for node := range g.Nodes {
		dist[node] = math.Inf(1)
		prev[node] = ""
	}
	if _, ok := g.Nodes[source]; !ok {
		return nil, errors.New("source node not in graph")
	}
	dist[source] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &PQItem{Node: source, Priority: 0})

	// Salva o estado inicial
	snapshots = append(snapshots, graph.Snapshot{
		Dist:      graph.CopyDist(dist),
		Prev:      graph.CopyPrev(prev),
		Visited:   graph.CopyVisited(visited),
		StepTitle: "Estado Inicial",
	})

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*PQItem)
		u := item.Node

		if visited[u] {
			continue
		}
		visited[u] = true

		// Snapshot para "Visitando Nó"
		snapshots = append(snapshots, graph.Snapshot{
			Dist:        graph.CopyDist(dist),
			Prev:        graph.CopyPrev(prev),
			Visited:     graph.CopyVisited(visited),
			StepTitle:   fmt.Sprintf("Visitando nó %s (dist: %.2f)", u, dist[u]),
			UpdatedNode: u, // Destaca o nó 'u' que saiu da PQ
		})

		for _, e := range g.Adj[u] {
			// Snapshot para "Checando Aresta"
			snap := graph.Snapshot{
				Dist:       graph.CopyDist(dist),
				Prev:       graph.CopyPrev(prev),
				Visited:    graph.CopyVisited(visited),
				StepTitle:  fmt.Sprintf("Checando %s->%s", e.From, e.To),
				ActiveEdge: &e,
			}

			v := e.To
			alt := dist[u] + e.W
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u
				heap.Push(pq, &PQItem{Node: v, Priority: alt})

				// Atualiza snapshot para "Relaxou"
				snap.UpdatedNode = v
				snap.StepTitle = fmt.Sprintf("Relaxou %s->%s! (nova dist: %.2f)", e.From, e.To, alt)
				snap.Dist = graph.CopyDist(dist) // Pega nova distância
				snap.Prev = graph.CopyPrev(prev) // Pega novo predecessor
			}
			snapshots = append(snapshots, snap)
		}
	}

	snapshots = append(snapshots, graph.Snapshot{
		Dist:      graph.CopyDist(dist),
		Prev:      graph.CopyPrev(prev),
		Visited:   graph.CopyVisited(visited),
		StepTitle: "Finalizado",
	})

	return snapshots, nil
}
