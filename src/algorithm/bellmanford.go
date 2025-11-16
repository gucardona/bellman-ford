package algorithm

import (
	"errors"
	"fmt"
	"math"

	"github.com/gucardona/bellman-ford/src/graph"
)

// BellmanFord executa o algoritmo e retorna uma lista de 'snapshots' (frames)
// da execução, um booleano indicando se há ciclo negativo, e um erro.
func BellmanFord(g *graph.Graph, source string) ([]graph.Snapshot, bool, error) {
	dist := make(map[string]float64)
	prev := make(map[string]string)

	// A lista de "cenas" para o nosso GIF
	var snapshots []graph.Snapshot

	for node := range g.Nodes {
		dist[node] = math.Inf(1)
		prev[node] = ""
	}
	if _, ok := g.Nodes[source]; !ok {
		return nil, false, errors.New("source node not in graph")
	}
	dist[source] = 0

	// Salva o estado inicial
	snapshots = append(snapshots, graph.Snapshot{
		Dist:      graph.CopyDist(dist),
		Prev:      graph.CopyPrev(prev),
		StepTitle: "Estado Inicial",
	})

	n := len(g.Nodes)
	// relax edges n-1 vezes
	for i := 0; i < n-1; i++ {
		// Snapshot no início de cada iteração
		snapshots = append(snapshots, graph.Snapshot{
			Dist:      graph.CopyDist(dist),
			Prev:      graph.CopyPrev(prev),
			StepTitle: fmt.Sprintf("Início da Iteração %d / %d", i+1, n-1),
		})

		changed := false
		for _, e := range g.Edges {
			// Grava snapshot ANTES de tentar relaxar
			snap := graph.Snapshot{
				Dist:       graph.CopyDist(dist),
				Prev:       graph.CopyPrev(prev),
				StepTitle:  fmt.Sprintf("Iter. %d: Checando %s->%s", i+1, e.From, e.To),
				ActiveEdge: &e, // Destaca a aresta
			}

			if dist[e.From]+e.W < dist[e.To] {
				dist[e.To] = dist[e.From] + e.W
				prev[e.To] = e.From
				changed = true

				// Atualiza o snapshot para mostrar a MUDANÇA
				snap.UpdatedNode = e.To // Destaca o nó
				snap.StepTitle = fmt.Sprintf("Iter. %d: Relaxou %s->%s!", i+1, e.From, e.To)
				// Salva os mapas *depois* da mudança
				snap.Dist = graph.CopyDist(dist)
				snap.Prev = graph.CopyPrev(prev)
			}
			snapshots = append(snapshots, snap) // Adiciona o frame
		}
		if !changed {
			break
		}
	}

	// Checa ciclo negativo (e grava snapshots)
	hasNegCycle := false
	for _, e := range g.Edges {
		snap := graph.Snapshot{
			Dist:       graph.CopyDist(dist),
			Prev:       graph.CopyPrev(prev),
			StepTitle:  fmt.Sprintf("Checando ciclo: %s->%s", e.From, e.To),
			ActiveEdge: &e,
		}
		if dist[e.From]+e.W < dist[e.To] {
			hasNegCycle = true
			snap.StepTitle = "CICLO NEGATIVO DETECTADO!"
			snap.UpdatedNode = e.To // Destaca o nó que "quebrou"
			snapshots = append(snapshots, snap)
			break // Para no primeiro ciclo
		}
		snapshots = append(snapshots, snap)
	}

	snapshots = append(snapshots, graph.Snapshot{
		Dist:      graph.CopyDist(dist),
		Prev:      graph.CopyPrev(prev),
		StepTitle: "Finalizado",
	})

	return snapshots, hasNegCycle, nil
}
