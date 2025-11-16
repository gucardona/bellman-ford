package generate

import (
	"fmt"
	"os"
	"strings"

	"github.com/gucardona/bellman-ford/src/graph"
)

// WriteDOT gera um .dot que destaca o estado atual do algoritmo (snapshot).
func WriteDOT(g *graph.Graph, snap *graph.Snapshot, outpath string, title string) error {
	f, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "digraph G {\n")

	// Usa o título do snapshot se ele existir, senão usa o título principal.
	dotTitle := title
	if snap != nil && snap.StepTitle != "" {
		dotTitle = snap.StepTitle
	}
	fmt.Fprintf(f, "  label=\"%s\";\n", escape(dotTitle))
	fmt.Fprintf(f, "  labelloc=top;\n")
	fmt.Fprintf(f, "  fontsize=20;\n")
	fmt.Fprintf(f, "  rankdir=LR;\n")

	fmt.Fprintf(f, "  graph [pad=\"0.2\", nodesep=\"0.5\", ranksep=\"0.5\", overlap=false, splines=true];\n")

	// 1. Define o estilo PADRÃO do nó (baseado no seu .dot original)
	fmt.Fprintf(f, "  node [shape=circle, style=filled, fillcolor=white, fontsize=12, fontname=\"Helvetica\"];\n")
	fmt.Fprintf(f, "  edge [fontsize=10, fontname=\"Helvetica\"];\n\n")

	// 2. Define todos os nós (sem estilo)
	// Isso garante que todos apareçam, mesmo se não tiverem estilo especial.
	for n := range g.Nodes {
		fmt.Fprintf(f, "  \"%s\";\n", escape(n))
	}
	fmt.Fprintf(f, "\n")

	// 3. Aplica estilos de override (cores)
	// Itera e colore os nós com base no snapshot.
	if snap != nil {
		for n := range g.Nodes {
			style := "" // O estilo agora é SÓ a cor
			if snap.Visited != nil && snap.Visited[n] {
				style = "fillcolor=lightgray"
			}
			if snap.UpdatedNode == n {
				style = "fillcolor=lightblue"
			}

			if style != "" {
				// Aplica o override de cor
				fmt.Fprintf(f, "  \"%s\" [%s];\n", escape(n), style)
			}
		}
	}
	fmt.Fprintf(f, "\n")

	// 4. Define as arestas (com seus estilos)
	// marcadores da árvore de caminhos (baseado no 'prev' do snapshot)
	treeEdges := map[string]struct{}{}
	if snap != nil && snap.Prev != nil {
		for v, u := range snap.Prev {
			if u == "" {
				continue
			}
			key := edgeKey(u, v)
			treeEdges[key] = struct{}{}
		}
	}

	for _, e := range g.Edges {
		style := "color=gray"
		if e.W < 0 {
			style = "color=red, style=dashed" // Aresta negativa (sempre)
		}

		if _, ok := treeEdges[edgeKey(e.From, e.To)]; ok {
			// Aresta pertence à árvore de caminhos
			style = "color=green, penwidth=2.5"
		}

		// Destaca a aresta ativa
		if snap != nil && snap.ActiveEdge != nil && snap.ActiveEdge.From == e.From && snap.ActiveEdge.To == e.To {
			style = "color=blue, penwidth=3.0, style=dashed"
		}

		fmt.Fprintf(f, "  \"%s\" -> \"%s\" [label=\"%s\" %s];\n", escape(e.From), escape(e.To), fmtWeight(e.W), style)
	}
	fmt.Fprintf(f, "}\n")
	return nil
}

// edgeKey cria uma chave única para uma aresta.
func edgeKey(a, b string) string {
	return a + "->" + b
}

// fmtWeight formata o peso da aresta para exibição.
func fmtWeight(w float64) string {
	// evita notação exponencial estranha
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.6f", w), "0"), ".")
}

// escape trata aspas em nomes de nós.
func escape(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}
