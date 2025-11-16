package generate

import (
	"fmt"
	"os/exec"

	"github.com/gucardona/bellman-ford/src/utils"
)

// RenderDotToPNG converte um arquivo .dot em .png usando graphviz.
// Tenta usar 'sfdp', 'neato' e 'dot' em ordem de preferência.
func RenderDotToPNG(dotPath, pngPath string) error {
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

// CreateGIF usa o ImageMagick (comando 'convert') para unir múltiplos
// arquivos .png em um .gif animado.
func CreateGIF(gifPath string, framePaths []string, delay string) error {
	if len(framePaths) == 0 {
		return fmt.Errorf("nenhum frame fornecido para o GIF")
	}

	// Usa 'convert' (ImageMagick)
	args := []string{"-delay", delay, "-loop", "0"}
	args = append(args, framePaths...) // Adiciona todos os PNGs
	args = append(args, gifPath)       // Arquivo de saída

	if err := utils.RunCmd("convert", args...); err != nil {
		return fmt.Errorf("erro criando GIF com 'convert': %w", err)
	}
	return nil
}
