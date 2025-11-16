package generate

import (
	"fmt"
	"os"
	"strings"

	"github.com/gucardona/bellman-ford/src/graph"
)

// WriteDOT gera um .dot que destaca a árvore de caminhos (predecessores).
// Se prev map estiver vazio, gera apenas o grafo.
func WriteDOT(g *graph.Graph, prev map[string]string, outpath string, title string) error {
	f, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "digraph G {\n")
	fmt.Fprintf(f, "  label=\"%s\";\n", escape(title))
	fmt.Fprintf(f, "  labelloc=top;\n")
	fmt.Fprintf(f, "  fontsize=20;\n")
	fmt.Fprintf(f, "  rankdir=LR;\n")

	fmt.Fprintf(f, "  graph [pad=\"0.2\", nodesep=\"0.5\", ranksep=\"0.5\", overlap=false, splines=true];\n")
	fmt.Fprintf(f, "  node [shape=circle, style=filled, fillcolor=white, fontsize=12, fontname=\"Helvetica\"];\n")
	fmt.Fprintf(f, "  edge [fontsize=10, fontname=\"Helvetica\"];\n\n")

	// imprimimos nós explícitos
	for n := range g.Nodes {
		fmt.Fprintf(f, "  \"%s\";\n", escape(n))
	}

	// marcadores da árvore de caminhos
	treeEdges := map[string]struct{}{}
	for v, u := range prev {
		if u == "" {
			continue
		}
		key := edgeKey(u, v)
		treeEdges[key] = struct{}{}
	}

	for _, e := range g.Edges {
		style := "color=gray"
		if e.W < 0 {
			style = "color=red, style=dashed"
		}
		if _, ok := treeEdges[edgeKey(e.From, e.To)]; ok {
			// aresta pertence à árvore de caminhos
			style = "color=green, penwidth=2.5"
		}
		fmt.Fprintf(f, "  \"%s\" -> \"%s\" [label=\"%s\" %s];\n", escape(e.From), escape(e.To), fmtWeight(e.W), style)
	}
	fmt.Fprintf(f, "}\n")
	return nil
}

func edgeKey(a, b string) string {
	return a + "->" + b
}

func fmtWeight(w float64) string {
	// evita notação exponencial estranha
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.6f", w), "0"), ".")
}

func escape(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}
