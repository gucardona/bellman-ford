package utils

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// RunCmd executa um comando externo e encaminha stdout/stderr para o processo atual.
func RunCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// FileExists verifica se o caminho existe (arquivo).
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// SaveDistances salva o mapa de distâncias em um arquivo.
func SaveDistances(dist map[string]float64, outpath string) {
	f, err := os.Create(outpath)
	if err != nil {
		fmt.Printf("Erro escrevendo %s: %v\n", outpath, err)
		return
	}
	defer f.Close()

	// ordem determinística
	keys := make([]string, 0, len(dist))
	for k := range dist {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := dist[k]

		if math.IsInf(v, 1) { // infinito positivo
			fmt.Fprintf(f, "%s: INF\n", k)
		} else {
			fmt.Fprintf(f, "%s: %f\n", k, v)
		}
	}
}

func EnsurePath() {
	extraPaths := []string{
		"/usr/local/bin",
		"/usr/bin",
		"/bin",
		"/opt/homebrew/bin", // Apple Silicon (M1/M2/M3)
		"/opt/local/bin",
		"/usr/local/sbin",
		"/opt/homebrew/sbin",
	}

	current := os.Getenv("PATH")
	for _, p := range extraPaths {
		if !strings.Contains(current, p) {
			current += ":" + p
		}
	}
	os.Setenv("PATH", current)
}
