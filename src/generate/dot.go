package generate

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/gucardona/bellman-ford/src/graph"
)

// WriteDOT gera um .dot que destaca o estado atual do algoritmo (snapshot).
func WriteDOT(g *graph.Graph, snap *graph.Snapshot, outpath string, title string, positions map[string]string, graphCenterX float64) error {
	f, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer f.Close()

	var sb strings.Builder

	sb.WriteString("digraph G {\n")
	sb.WriteString("  layout=neato;\n")
	sb.WriteString("  overlap=false;\n")
	sb.WriteString("  splines=true;\n")
	sb.WriteString("  node [shape=circle, style=filled, fontsize=16, width=0.8, height=0.8];\n")
	sb.WriteString("  edge [fontsize=14, fontcolor=\"#333333\", dir=forward];\n")
	sb.WriteString("  labelloc=\"t\";\n")

	// Título do Snapshot
	dotTitle := title
	if snap != nil && snap.StepTitle != "" {
		dotTitle = snap.StepTitle
	}
	sb.WriteString(fmt.Sprintf("  label=\"%s\";\n", escape(dotTitle)))
	sb.WriteString("  fontsize=20;\n")
	sb.WriteString("  fontname=\"Arial Bold\";\n\n")

	// 1. Define todos os nós COM seus estilos E POSIÇÕES
	for n := range g.Nodes {
		color := getVertexColor(n, snap)
		distLabel := getDistanceLabel(snap.Dist[n])

		pos, ok := positions[n]
		if !ok {
			return fmt.Errorf("posição não encontrada para o nó: %s", n)
		}

		sb.WriteString(fmt.Sprintf("  \"%s\" [fillcolor=\"%s\", label=\"%s\n%s\", pos=\"%s\"];\n",
			escape(n), color, escape(n), distLabel, pos))
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
	sb.WriteString(generateTable(snap, graphCenterX))

	sb.WriteString("}\n")

	_, err = f.Write([]byte(sb.String()))
	return err
}

// generateTable cria a tabela HTML para o .dot
func generateTable(step *graph.Snapshot, graphCenterX float64) string {
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

	// --- INÍCIO DA MUDANÇA (Bug de Espaçamento) ---
	// Diminuímos o Y para ficar mais próximo (o grafo agora começa em 120)
	const tableYPos = 0.0
	// --- FIM DA MUDANÇA ---

	sb.WriteString(fmt.Sprintf("  table [shape=plaintext, label=%s, pos=\"%f,%f!\"];\n",
		tableLabel, graphCenterX, tableYPos))

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

// CalculateLayout executa o 'neato' uma vez para obter posições estáveis dos nós.
// Retorna um mapa de [nomeDoNó] -> "x,y!" E a coordenada X central
func CalculateLayout(g *graph.Graph) (map[string]string, float64, error) {
	var sb strings.Builder
	sb.WriteString("digraph G {\n")
	sb.WriteString("  layout=neato;\n")
	sb.WriteString("  overlap=false;\n")

	// Define nós
	for n := range g.Nodes {
		sb.WriteString(fmt.Sprintf("  \"%s\";\n", escape(n)))
	}
	// Define arestas
	for _, e := range g.Edges {
		sb.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\";\n", escape(e.From), escape(e.To)))
	}
	sb.WriteString("}\n")

	// Executa 'neato -Tplain'
	cmd := exec.Command("neato", "-Tplain")
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, 0, fmt.Errorf("falha ao rodar 'neato' para calcular layout: %v\nStderr: %s", err, stderr.String())
	}

	positions := make(map[string]string)
	scanner := bufio.NewScanner(&out)

	const verticalOffset = 2.0
	var totalX float64 = 0
	var nodeCount float64 = 0

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) > 0 && parts[0] == "node" {
			if len(parts) >= 4 {
				nodeName := parts[1]
				xStr := parts[2]
				yStr := parts[3]

				x, errX := strconv.ParseFloat(xStr, 64)
				y, errY := strconv.ParseFloat(yStr, 64)

				if errX != nil || errY != nil {
					positions[nodeName] = fmt.Sprintf("%s,%s!", xStr, yStr) // fallback
					continue
				}

				// Acumula X para média
				totalX += x
				nodeCount++

				// Adiciona o offset e formata de volta para string
				yWithOffset := y + verticalOffset

				// O "!" é crucial para travar a posição
				positions[nodeName] = fmt.Sprintf("%f,%f!", x, yWithOffset)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, fmt.Errorf("erro ao ler saída do neato: %w", err)
	}

	if len(positions) == 0 {
		return nil, 0, fmt.Errorf("nenhuma posição de nó foi calculada. Verifique se 'neato' (Graphviz) está instalado")
	}

	// Calcula o X central
	graphCenterX := 0.0
	if nodeCount > 0 {
		graphCenterX = totalX / nodeCount
	}

	return positions, graphCenterX, nil
}
