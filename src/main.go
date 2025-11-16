package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gucardona/bellman-ford/src/algorithm"
	"github.com/gucardona/bellman-ford/src/generate"
	"github.com/gucardona/bellman-ford/src/graph"
	"github.com/gucardona/bellman-ford/src/utils"
)

func main() {
	infile := flag.String("i", "./src/edges.csv", "arquivo CSV de entrada (from,to,weight)")
	source := flag.String("s", "", "nó fonte (obrigatório)")
	outdir := flag.String("o", "./out", "diretório de saída para .dot, .png, .gif, .txt")
	// Flag 'skipImages' removida para focar na geração de frames
	flag.Parse()

	utils.EnsurePath()

	path, _ := exec.LookPath("dot")
	fmt.Println("dot localizado em:", path)

	if *source == "" {
		fmt.Println("Erro: informe -s <source node>")
		os.Exit(1)
	}

	// Carrega grafo
	g := graph.NewGraph()
	if err := g.LoadFromCSV(*infile); err != nil {
		fmt.Printf("Erro ao ler CSV '%s': %v\n", *infile, err)
		os.Exit(1)
	}

	// Cria diretório de saída
	if err := os.MkdirAll(*outdir, 0755); err != nil {
		fmt.Printf("Erro criando diretório de saída '%s': %v\n", *outdir, err)
		os.Exit(1)
	}

	// Detecta arestas negativas
	hasNegative := false
	for _, e := range g.Edges {
		if e.W < 0 {
			hasNegative = true
			break
		}
	}

	// === Executa Dijkstra (Modo Didático) ===
	// Executa sempre. Se houver arestas negativas, o resultado estará
	// incorreto, o que é útil para a comparação.
	fmt.Println("Executando Dijkstra...")
	if hasNegative {
		fmt.Println("AVISO: Arestas negativas detectadas. O resultado do Dijkstra pode estar incorreto.")
	}

	dijkstraSnapshots, errD := algorithm.Dijkstra(g, *source)
	if errD != nil {
		fmt.Printf("Erro Dijkstra: %v\n", errD)
	} else {
		// --- Geração do GIF do Dijkstra ---
		fmt.Println("Gerando frames do Dijkstra...")
		dijkstraFrameDir := filepath.Join(*outdir, "dijkstra_frames")
		os.MkdirAll(dijkstraFrameDir, 0755)
		dijkstraFramePaths := generateFrames(g, dijkstraSnapshots, dijkstraFrameDir, "Dijkstra")

		if len(dijkstraFramePaths) > 0 {
			gifPath := filepath.Join(*outdir, "dijkstra_animation.gif")
			fmt.Println("Montando GIF animado do Dijkstra:", gifPath)
			if err := generate.CreateGIF(gifPath, dijkstraFramePaths, "100"); err != nil { // delay 100
				fmt.Printf("Erro criando GIF do Dijkstra: %v\n", err)
			}
		}

		// --- Salva o resultado FINAL (último snapshot) ---
		finalDijkstraSnap := dijkstraSnapshots[len(dijkstraSnapshots)-1]
		dotPath := filepath.Join(*outdir, "dijkstra.dot")
		title := "Dijkstra - Árvore de caminhos"
		if hasNegative {
			title += " (Resultado PODE ESTAR INCORRETO)"
		}
		generate.WriteDOT(g, &finalDijkstraSnap, dotPath, title)
		utils.SaveDistances(finalDijkstraSnap.Dist, filepath.Join(*outdir, "dijkstra_dist.txt"))
	}

	// === Executa Bellman-Ford (Modo Animação) ===
	fmt.Println("Executando Bellman-Ford (Modo Animação)...")
	snapshots, negCycle, err := algorithm.BellmanFord(g, *source)
	if err != nil {
		fmt.Printf("Erro Bellman-Ford: %v\n", err)
	} else {
		// --- Geração do GIF do Bellman-Ford ---
		fmt.Println("Gerando frames do Bellman-Ford...")
		frameDir := filepath.Join(*outdir, "bf_frames")
		os.MkdirAll(frameDir, 0755)
		framePaths := generateFrames(g, snapshots, frameDir, "Bellman-Ford")

		if len(framePaths) > 0 {
			gifPath := filepath.Join(*outdir, "bellman_animation.gif")
			fmt.Println("Montando GIF animado do Bellman-Ford:", gifPath)
			if err := generate.CreateGIF(gifPath, framePaths, "80"); err != nil { // delay 80
				fmt.Printf("Erro criando GIF animado: %v\n", err)
			}
		}

		// --- Salva o resultado FINAL (último snapshot) ---
		finalSnap := snapshots[len(snapshots)-1]
		title := "Bellman-Ford - Árvore de caminhos"
		if negCycle {
			title += " (CICLO NEGATIVO DETECTADO)"
		}
		dotPath := filepath.Join(*outdir, "bellman.dot")
		generate.WriteDOT(g, &finalSnap, dotPath, title) // Passa o snapshot final
		utils.SaveDistances(finalSnap.Dist, filepath.Join(*outdir, "bellman_dist.txt"))
		if negCycle {
			fmt.Println("Aviso: ciclo negativo detectado.")
		}
	}

	fmt.Println("Processo concluído. Verifique o diretório:", *outdir)
}

// generateFrames itera sobre os snapshots, gera um .dot e .png para cada um,
// e retorna a lista de caminhos dos .png gerados.
func generateFrames(g *graph.Graph, snapshots []graph.Snapshot, outdir string, algName string) []string {
	framePaths := []string{}
	for i, snap := range snapshots {
		dotPath := filepath.Join(outdir, fmt.Sprintf("frame_%04d.dot", i))
		pngPath := filepath.Join(outdir, fmt.Sprintf("frame_%04d.png", i))

		// Passa o snapshot atual para o WriteDOT
		title := fmt.Sprintf("%s - %s", algName, snap.StepTitle)
		if err := generate.WriteDOT(g, &snap, dotPath, title); err != nil {
			fmt.Printf("Erro frame %d: %v\n", i, err)
			continue
		}

		// Converte .dot -> .png
		if err := generate.RenderDotToPNG(dotPath, pngPath); err != nil {
			fmt.Printf("Erro renderizando %s: %v\n", dotPath, err)
			continue
		}
		framePaths = append(framePaths, pngPath)
	}
	return framePaths
}
