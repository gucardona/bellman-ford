package generate

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/gucardona/bellman-ford/src/utils"
)

func GenerateImagesAndGIF(outdir string) error {
	dijkstraDot := filepath.Join(outdir, "dijkstra.dot")
	bellmanDot := filepath.Join(outdir, "bellman.dot")
	dijkstraPng := filepath.Join(outdir, "dijkstra.png")
	bellmanPng := filepath.Join(outdir, "bellman.png")
	gifFile := filepath.Join(outdir, "comparison.gif")

	// Converte DOT -> PNG usando Graphviz
	if utils.FileExists(dijkstraDot) {
		if err := renderDotToPNG(dijkstraDot, dijkstraPng); err != nil {
			return fmt.Errorf("erro convertendo %s: %w", dijkstraDot, err)
		}
		fmt.Println("Gerado:", dijkstraPng)
	}

	if utils.FileExists(bellmanDot) {
		if err := renderDotToPNG(bellmanDot, bellmanPng); err != nil {
			return fmt.Errorf("erro convertendo %s: %w", bellmanDot, err)
		}
		fmt.Println("Gerado:", bellmanPng)
	}

	// Se ambos os PNG existirem, gera GIF
	if utils.FileExists(dijkstraPng) && utils.FileExists(bellmanPng) {
		if err := utils.RunCmd("convert", "-delay", "120", "-loop", "0", dijkstraPng, bellmanPng, gifFile); err != nil {
			return fmt.Errorf("erro criando GIF: %w", err)
		}
		fmt.Println("GIF criado:", gifFile)
	} else {
		fmt.Println("PNG(s) faltando, GIF não foi gerado.")
	}

	return nil
}

func renderDotToPNG(dotPath, pngPath string) error {
	// prioridade: sfdp (bom para grafos maiores), depois neato, depois dot
	tryCmds := [][]string{
		// sfdp: bom para layouts em força em grafos maiores
		{"sfdp", "-Tpng", "-Goverlap=false", "-Gsplines=true", "-Gdpi=150", dotPath, "-o", pngPath},
		// neato: bom para layouts force-directed em grafos menores
		{"neato", "-Tpng", "-Goverlap=false", "-Gsplines=true", "-Gdpi=150", dotPath, "-o", pngPath},
		// fallback dot com flags para melhorar visual
		{"dot", "-Tpng", "-Gdpi=150", "-Nfontsize=12", "-Efontsize=10", dotPath, "-o", pngPath},
	}

	for _, cmdArgs := range tryCmds {
		name := cmdArgs[0]
		if _, err := exec.LookPath(name); err != nil {
			// comando não disponível, tenta próximo
			continue
		}
		// executar comando
		if err := utils.RunCmd(cmdArgs[0], cmdArgs[1:]...); err != nil {
			// se falhou ao executar, continua para tentar o próximo
			fmt.Printf("tentativa com %s falhou: %v\n", name, err)
			continue
		}
		// sucesso
		return nil
	}
	return fmt.Errorf("nenhum renderer disponível (tentadas: sfdp, neato, dot) ou todas falharam")
}
