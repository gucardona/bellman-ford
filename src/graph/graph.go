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

// Snapshot armazena o estado do algoritmo em um ponto específico no tempo.
// Usado para gerar frames de uma animação.
type Snapshot struct {
	Dist      map[string]float64 // Distâncias atuais
	Prev      map[string]string  // Predecessores atuais
	StepTitle string             // Ex: "Iteração 2", "Relaxando A->B"

	// Para destacar no .dot
	ActiveEdge  *Edge           // Aresta sendo processada (azul, tracejada)
	UpdatedNode string          // Nó que teve a distância atualizada (preenchido)
	Visited     map[string]bool // Para Dijkstra: nós já visitados (cinza)
}

// CopyDist cria uma cópia profunda de um mapa de distâncias.
func CopyDist(m map[string]float64) map[string]float64 {
	newMap := make(map[string]float64)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

// CopyPrev cria uma cópia profunda de um mapa de predecessores.
func CopyPrev(m map[string]string) map[string]string {
	newMap := make(map[string]string)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

// CopyVisited cria uma cópia profunda de um mapa de nós visitados.
func CopyVisited(m map[string]bool) map[string]bool {
	newMap := make(map[string]bool)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
