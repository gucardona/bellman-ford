package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/gucardona/bellman-ford/src/graph"
)

// GenerateFrames itera sobre os snapshots, gera um .dot e .png para cada um.
// Esta função agora também renderiza o PNG, forçando o uso do 'dot'.
// Adaptado de 'dijkstra-visualizer/internal/visualizer/media.go'
func GenerateFrames(g *graph.Graph, snapshots []graph.Snapshot, outputDir string, algName string) ([]string, error) {
	framePaths := []string{}
	for i, snap := range snapshots {
		// Gera o título para este frame
		title := fmt.Sprintf("%s - %s", algName, snap.StepTitle)
		dotFile := filepath.Join(outputDir, fmt.Sprintf("frame_%04d.dot", i))
		pngFile := filepath.Join(outputDir, fmt.Sprintf("frame_%04d.png", i))

		// 1. Chama o WriteDOT (que agora escreve os atributos 'pos')
		if err := WriteDOT(g, &snap, dotFile, title); err != nil {
			fmt.Printf("Aviso: falha ao escrever .dot para frame %d: %v\n", i, err)
			continue // Pula este frame
		}

		// 2. Força o uso do 'dot' para renderizar (lógica do dijkstra-visualizer)
		// Isto respeita o 'layout=neato' escrito no .dot.
		cmd := exec.Command("dot", "-Tpng", dotFile, "-o", pngFile)
		if output, err := cmd.CombinedOutput(); err != nil {
			// Se 'dot' falhar, avisa e pula
			fmt.Printf("Aviso: falha ao renderizar .png para frame %d: %w\nOutput: %s\n", i, err, output)
			continue // Pula este frame
		}

		// 3. Adiciona o frame bem-sucedido à lista
		framePaths = append(framePaths, pngFile)
	}

	if len(framePaths) == 0 {
		return nil, fmt.Errorf("nenhum frame .png foi gerado com sucesso. Verifique se 'dot' (Graphviz) está instalado")
	}
	return framePaths, nil
}

// CreateGIF usa 'magick' ou 'convert' para criar o GIF animado.
// Adaptado de 'dijkstra-visualizer/internal/visualizer/media.go'
func CreateGIF(gifPath string, framePaths []string, delay string) error {
	if len(framePaths) == 0 {
		return fmt.Errorf("nenhum frame fornecido para o GIF")
	}

	sort.Strings(framePaths) // Garante a ordem

	args := []string{
		"-delay", delay,
		"-loop", "0",
	}
	args = append(args, framePaths...)
	args = append(args, gifPath)

	// Tenta 'magick' (ImageMagick 7+)
	cmd := exec.Command("magick", args...)
	if _, err := cmd.CombinedOutput(); err != nil {
		// Se falhar, tenta 'convert' (ImageMagick 6 e legados)
		cmd = exec.Command("convert", args...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("falha ao criar GIF com 'magick' e 'convert'. Verifique se o ImageMagick está instalado: %w\nOutput: %s", err, output)
		}
	}

	// Limpa os .dot
	dotFiles, _ := filepath.Glob(filepath.Join(filepath.Dir(framePaths[0]), "*.dot"))
	for _, file := range dotFiles {
		os.Remove(file)
	}

	return nil
}
