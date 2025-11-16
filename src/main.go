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

// main.go atualizado:
// - lê CSV
// - roda Dijkstra e Bellman-Ford
// - gera .dot e arquivos de distâncias (usa WriteDOT e saveDistances)
// - converte .dot -> .png com 'dot'
// - gera comparison.gif com 'convert'
// Não depende de um script .sh externo.

func main() {
	infile := flag.String("i", "./src/edges.csv", "arquivo CSV de entrada (from,to,weight)")
	source := flag.String("s", "", "nó fonte (obrigatório)")
	outdir := flag.String("o", "./out", "diretório de saída para .dot, .png, .gif, .txt")
	skipImages := flag.Bool("noimages", false, "se true, não gera PNG/GIF (apenas .dot/.txt)")
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

	// Executa Dijkstra (pula se detectar aresta negativa)
	hasNegative := false
	for _, e := range g.Edges {
		if e.W < 0 {
			hasNegative = true
			break
		}
	}

	if !hasNegative {
		fmt.Println("Executando Dijkstra...")
		distD, prevD, err := algorithm.Dijkstra(g, *source)
		if err != nil {
			fmt.Printf("Erro Dijkstra: %v\n", err)
		} else {
			dotPath := filepath.Join(*outdir, "dijkstra.dot")
			if err := generate.WriteDOT(g, prevD, dotPath, "Dijkstra - Árvore de caminhos"); err != nil {
				fmt.Printf("Erro escrevendo DOT Dijkstra: %v\n", err)
			} else {
				fmt.Printf("DOT Dijkstra gerado: %s\n", dotPath)
			}
			utils.SaveDistances(distD, filepath.Join(*outdir, "dijkstra_dist.txt"))
		}
	} else {
		fmt.Println("Arestas negativas detectadas — pulando Dijkstra (não é seguro para pesos negativos).")
	}

	// Executa Bellman-Ford
	fmt.Println("Executando Bellman-Ford...")
	distBF, prevBF, negCycle, err := algorithm.BellmanFord(g, *source)
	if err != nil {
		fmt.Printf("Erro Bellman-Ford: %v\n", err)
	} else {
		title := "Bellman-Ford - Árvore de caminhos"
		if negCycle {
			title += " (CICLO NEGATIVO DETECTADO)"
		}
		dotPath := filepath.Join(*outdir, "bellman.dot")
		if err := generate.WriteDOT(g, prevBF, dotPath, title); err != nil {
			fmt.Printf("Erro escrevendo DOT Bellman-Ford: %v\n", err)
		} else {
			fmt.Printf("DOT Bellman-Ford gerado: %s\n", dotPath)
		}
		utils.SaveDistances(distBF, filepath.Join(*outdir, "bellman_dist.txt"))
		if negCycle {
			fmt.Println("Aviso: ciclo negativo detectado. Distâncias podem não ser válidas.")
		}
	}

	// Se o usuário pediu para não gerar imagens, encerra aqui
	if *skipImages {
		fmt.Println("skipImages=true: finalizando sem gerar PNG/GIF.")
		return
	}

	// Gera PNGs e GIF diretamente em Go
	fmt.Println("Gerando PNGs e GIF (usando 'dot' e 'convert')...")
	if err := generate.GenerateImagesAndGIF(*outdir); err != nil {
		fmt.Printf("Erro ao gerar PNG/GIF: %v\n", err)
	} else {
		fmt.Println("PNG e GIF gerados com sucesso.")
	}

	fmt.Println("Processo concluído. Verifique o diretório:", *outdir)
}
