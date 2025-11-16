package graph

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// Edge representa uma aresta do grafo (origem -> destino) com peso.
type Edge struct {
	From string
	To   string
	W    float64
}

type Graph struct {
	Nodes map[string]struct{}
	Adj   map[string][]Edge // adj list keyed by node name
	Edges []Edge
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]struct{}),
		Adj:   make(map[string][]Edge),
		Edges: []Edge{},
	}
}

// LoadFromCSV espera um CSV com cabeçalho: from,to,weight
// Os nós podem ser strings (ex: A, B, 1, nodeX)
func (g *Graph) LoadFromCSV(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return err
	}
	start := 0
	if len(records) == 0 {
		return fmt.Errorf("CSV vazio")
	}
	// detecta cabeçalho simples
	if len(records[0]) >= 3 {
		if records[0][0] == "from" && records[0][1] == "to" {
			start = 1
		}
	}
	for i := start; i < len(records); i++ {
		row := records[i]
		if len(row) < 3 {
			continue
		}
		w, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return err
		}
		e := Edge{From: row[0], To: row[1], W: w}
		g.AddEdge(e)
	}
	return nil
}

func (g *Graph) AddEdge(e Edge) {
	g.Nodes[e.From] = struct{}{}
	g.Nodes[e.To] = struct{}{}
	g.Adj[e.From] = append(g.Adj[e.From], e)
	g.Edges = append(g.Edges, e)
}
