# Projeto: Comparação Dijkstra vs Bellman-Ford (Go)

Implementação e comparação dos algoritmos **Dijkstra** e **Bellman-Ford** em Go.  
O projeto gera arquivos `.dot` (Graphviz), exporta distâncias e produz imagens (`.png`) e um `.gif` comparativo para visualização.

## Estrutura do repositório atual:
```

bellmanford/
├─ src/
│  ├─ algorithm/
│  │  ├─ dijkstra.go
│  │  └─ bellmanford.go
│  ├─ generate/
│  │  ├─ dot.go        # geração de .dot a partir do grafo / predecessor
│  │  └─ gif.go        # geração de PNG/GIF (pode usar dot/convert ou implementação Go)
│  ├─ graph/
│  │  └─ graph.go      # carregamento do CSV, estrutura do grafo
│  ├─ utils/
│  │  └─ helper.go      # funções utilitárias
│  ├─ main.go          # entrypoint do programa
│  └─ edges.csv        # arquivo de arestas (exemplos)
├─ go.mod
└─ README.md

```

---

## Requisitos

- Go 1.20+ (ou versão compatível)
- Opcional (melhor renderização/layout): Graphviz (`dot`)
- Opcional (se usar ImageMagick): `convert` ou `magick`

> Observação: a pasta `src/generate` contém código para gerar imagens. Dependendo de como `gif.go` foi implementado, o gerador pode:
> - usar diretamente `dot` e `convert` quando presentes (melhor layout), ou
> - renderizar imagens/GIF com código Go puro (sem dependências externas).
>
> Verifique o comentário no topo de `src/generate/gif.go` para o comportamento exato.

---

## Formato do `edges.csv`

Arquivo CSV com cabeçalho e colunas:

```

from,to,weight
S,A,4
A,B,-2
...

```

- `from` e `to`: identificadores dos nós (strings)
- `weight`: número real (pode ser negativo)
- Grafo dirigido

---

## Como executar

A partir da raiz do projeto (onde está `go.mod`), rode:

```

go run ./src -i ./src/edges.csv -s S -o ./out

```

Parâmetros:
- `-i` caminho para o CSV de entrada (padrão: `./src/edges.csv`)
- `-s` nó fonte (obrigatório)
- `-o` diretório de saída (padrão: `./out`)
- `-noimages` (flag opcional) se você quiser gerar apenas `.dot` e arquivos de distância sem criar imagens/GIF

Exemplos:

1. Gerar tudo (algoritmos + imagens + GIF):
```

go run ./src -i ./src/edges.csv -s S -o ./out

```

2. Gerar apenas arquivos de saída sem imagens:
```

go run ./src -i ./src/edges.csv -s S -o ./out -noimages=true

```

---

## Saída esperada

Após a execução (quando `-noimages` não está ativo), o diretório de saída conterá algo como:

```

out/
├─ dijkstra.dot        # árvore de caminhos Dijkstra (se aplicada)
├─ bellman.dot         # árvore de caminhos Bellman-Ford
├─ dijkstra.png        # PNG gerado a partir do .dot (se disponível)
├─ bellman.png         # PNG gerado a partir do .dot (se disponível)
├─ comparison.gif      # GIF comparando Dijkstra x Bellman-Ford (se ambos PNG existirem)
├─ dijkstra_dist.txt   # distâncias calculadas por Dijkstra (se executado)
└─ bellman_dist.txt    # distâncias calculadas por Bellman-Ford

```

Se houver arestas de peso negativo, o programa por padrão **pula Dijkstra** (avisa no log) — use Bellman-Ford para obter resultados corretos; se houver ciclo negativo alcançável, Bellman-Ford informará.

---

## Observações / dicas

- Se `src/generate/gif.go` utiliza `dot`/`convert`, instale Graphviz e ImageMagick para melhor renderização:
  - Debian/Ubuntu: `sudo apt-get install graphviz imagemagick`
  - macOS (Homebrew): `brew install graphviz imagemagick`
- Se preferir manter o controle manual, ou se a versão `gif.go` não estiver disponível, existe a opção de usar um script externo (`make_gif.sh`) — caso ainda tenha, ele pode ser usado do mesmo modo:
```

./make_gif.sh ./src/edges.csv S ./out

```
- Para debug visual rápido, abra os `.dot` com Graphviz GUI ou converta manualmente:
```

dot -Tpng out/bellman.dot -o out/bellman.png

```

---
