package generate

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/gucardona/bellman-ford/src/graph"
)

// WriteDOT gera um .dot que destaca o estado atual do algoritmo (snapshot).
// Convertido para usar a lógica do 'dijkstra-visualizer' (layout neato, tabela HTML).
func WriteDOT(g *graph.Graph, snap *graph.Snapshot, outpath string, title string) error {
	f, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer f.Close()

	var sb strings.Builder

	sb.WriteString("digraph G {\n")
	// Usa 'neato' para um layout estável "force-directed"
	sb.WriteString("  layout=neato;\n")
	sb.WriteString("  overlap=false;\n")
	sb.WriteString("  splines=true;\n")
	sb.WriteString("  node [shape=circle, style=filled, fontsize=16, width=0.8, height=0.8];\n")
	sb.WriteString("  edge [fontsize=14, fontcolor=\"#333333\", dir=forward];\n") // dijkstra era 'dir=none'
	sb.WriteString("  labelloc=\"t\";\n")

	// Título do Snapshot
	dotTitle := title
	if snap != nil && snap.StepTitle != "" {
		dotTitle = snap.StepTitle
	}
	sb.WriteString(fmt.Sprintf("  label=\"%s\";\n", escape(dotTitle)))
	sb.WriteString("  fontsize=20;\n")
	sb.WriteString("  fontname=\"Arial Bold\";\n\n")

	// 1. Define todos os nós COM seus estilos
	for n := range g.Nodes {
		color := getVertexColor(n, snap)
		distLabel := getDistanceLabel(snap.Dist[n])

		// Não usamos 'pos' fixas, pois o CSV é dinâmico
		sb.WriteString(fmt.Sprintf("  \"%s\" [fillcolor=\"%s\", label=\"%s\\n%s\"];\n",
			escape(n), color, escape(n), distLabel))
	}
	sb.WriteString("\n")

	// 2. Define as arestas (com estilos priorizados)
	treeEdges := map[string]struct{}{}
	if snap != nil && snap.Prev != nil {
		for v, u := range snap.Prev {
			if u == "" {
				continue
			}
			treeEdges[edgeKey(u, v)] = struct{}{}
		}
	}

	for _, e := range g.Edges {
		color := getEdgeColor(e, snap, treeEdges)
		width := getEdgeWidth(e, snap, treeEdges)
		weightStr := fmtWeight(e.W)

		sb.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\", color=\"%s\", penwidth=%d];\n",
			escape(e.From), escape(e.To), weightStr, color, width))
	}

	sb.WriteString("\n")
	// 3. Gera a tabela HTML (lógica exata do dijkstra-visualizer)
	sb.WriteString(generateTable(snap))

	sb.WriteString("}\n")

	_, err = f.Write([]byte(sb.String()))
	return err
}

// generateTable cria a tabela HTML para o .dot
func generateTable(step *graph.Snapshot) string {
	if step == nil {
		return ""
	}

	var sb strings.Builder

	// Ordena os vértices para a tabela
	vertices := make([]string, 0, len(step.Dist))
	for v := range step.Dist {
		vertices = append(vertices, v)
	}
	sort.Strings(vertices)

	tableLabel := "<<TABLE BORDER=\"2\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"12\" BGCOLOR=\"white\">"
	tableLabel += "<TR><TD BGCOLOR=\"#4A90E2\"><FONT COLOR=\"white\" POINT-SIZE=\"18\"><B>Vertex</B></FONT></TD>"
	tableLabel += "<TD BGCOLOR=\"#4A90E2\"><FONT COLOR=\"white\" POINT-SIZE=\"18\"><B>Distance</B></FONT></TD>"
	tableLabel += "<TD BGCOLOR=\"#4A90E2\"><FONT COLOR=\"white\" POINT-SIZE=\"18\"><B>Previous</B></FONT></TD>"
	tableLabel += "<TD BGCOLOR=\"#4A90E2\"><FONT COLOR=\"white\" POINT-SIZE=\"18\"><B>Visited</B></FONT></TD></TR>"

	for _, v := range vertices {
		dist := step.Dist[v]
		prev := step.Prev[v]
		visited := step.Visited != nil && step.Visited[v] // Checa se 'Visited' existe

		bgColor := "#FFFFFF"
		if v == step.UpdatedNode { // No Bellman-Ford, 'UpdatedNode' é mais útil que 'CurrentVertex'
			bgColor = "#FFD700" // Amarelo
		} else if visited {
			bgColor = "#90EE90" // Verde claro
		}

		distStr := getDistanceLabel(dist)
		prevStr := "-"
		if prev != "" {
			prevStr = prev
		}
		visitedStr := "No"
		if visited {
			visitedStr = "Yes"
		}

		tableLabel += fmt.Sprintf("<TR><TD BGCOLOR=\"%s\"><FONT POINT-SIZE=\"16\"><B>%s</B></FONT></TD>", bgColor, escapeHTML(v))
		tableLabel += fmt.Sprintf("<TD BGCOLOR=\"%s\"><FONT POINT-SIZE=\"16\">%s</FONT></TD>", bgColor, distStr)
		tableLabel += fmt.Sprintf("<TD BGCOLOR=\"%s\"><FONT POINT-SIZE=\"16\">%s</FONT></TD>", bgColor, escapeHTML(prevStr))
		tableLabel += fmt.Sprintf("<TD BGCOLOR=\"%s\"><FONT POINT-SIZE=\"16\">%s</FONT></TD></TR>", bgColor, visitedStr)
	}

	tableLabel += "</TABLE>>"

	// Define a tabela num nó 'invisível' e o posiciona (layout=neato não usa rank)
	// Posiciona a tabela abaixo do grafo
	sb.WriteString(fmt.Sprintf("  table [shape=plaintext, label=%s];\n", tableLabel))

	return sb.String()
}

// Funções auxiliares (adaptadas do dijkstra-visualizer)

func getVertexColor(vertex string, step *graph.Snapshot) string {
	if step == nil {
		return "#ADD8E6"
	} // Não visitado (azul claro)

	// No Bellman-Ford, não temos 'CurrentVertex' da mesma forma, usamos 'UpdatedNode'
	if vertex == step.UpdatedNode {
		return "#FFD700" // Amarelo
	}
	if step.Visited != nil && step.Visited[vertex] {
		return "#90EE90" // Visitado (verde claro)
	}
	return "#ADD8E6" // Não visitado
}

func getDistanceLabel(distance float64) string {
	if math.IsInf(distance, 1) {
		return "∞"
	}
	// Formata para 2 casas decimais, mas remove .00
	return strings.TrimSuffix(fmt.Sprintf("%.2f", distance), ".00")
}

func getEdgeColor(edge graph.Edge, step *graph.Snapshot, treeEdges map[string]struct{}) string {
	if step != nil && step.ActiveEdge != nil {
		if step.ActiveEdge.From == edge.From && step.ActiveEdge.To == edge.To {
			return "#FF4500" // Laranja (explorando)
		}
	}

	if _, ok := treeEdges[edgeKey(edge.From, edge.To)]; ok {
		return "#32CD32" // Verde (na árvore)
	}

	if edge.W < 0 {
		return "#FF0000" // Vermelho (negativo)
	}

	return "#888888" // Preto/Cinza padrão
}

func getEdgeWidth(edge graph.Edge, step *graph.Snapshot, treeEdges map[string]struct{}) int {
	if step != nil && step.ActiveEdge != nil {
		if step.ActiveEdge.From == edge.From && step.ActiveEdge.To == edge.To {
			return 4
		}
	}
	if _, ok := treeEdges[edgeKey(edge.From, edge.To)]; ok {
		return 3
	}
	return 2
}

func edgeKey(a, b string) string {
	return a + "->" + b
}

func fmtWeight(w float64) string {
	return strings.TrimSuffix(fmt.Sprintf("%.2f", w), ".00")
}

func escape(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
