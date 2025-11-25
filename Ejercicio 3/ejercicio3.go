package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Edge struct {
	from, to, cost int
}

type UnionFind struct {
	parent, rank []int
}

func NewUnionFind(n int) *UnionFind {
	uf := &UnionFind{
		parent: make([]int, n),
		rank:   make([]int, n),
	}
	for i := range uf.parent {
		uf.parent[i] = i
	}
	return uf
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
	rootX, rootY := uf.Find(x), uf.Find(y)
	if rootX == rootY {
		return false
	}

	if uf.rank[rootX] < uf.rank[rootY] {
		uf.parent[rootX] = rootY
	} else if uf.rank[rootX] > uf.rank[rootY] {
		uf.parent[rootY] = rootX
	} else {
		uf.parent[rootY] = rootX
		uf.rank[rootX]++
	}
	return true
}

func primMST(n int, edges []Edge) (int, []Edge, int) {
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].cost < edges[j].cost
	})

	uf := NewUnionFind(n)
	var mstEdges []Edge
	totalCost := 0
	allCost := 0

	for _, e := range edges {
		allCost += e.cost
		if uf.Union(e.from, e.to) {
			mstEdges = append(mstEdges, e)
			totalCost += e.cost
			if len(mstEdges) == n-1 {
				break
			}
		}
	}

	for i := len(mstEdges); i < len(edges); i++ {
		allCost += edges[i].cost
	}

	return totalCost, mstEdges, allCost
}

func readUSGrid(filename string) (int, []Edge) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error al abrir archivo:", err)
		return 0, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var edges []Edge
	nodeSet := make(map[int]bool)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "%") || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		from, err1 := strconv.Atoi(parts[0])
		to, err2 := strconv.Atoi(parts[1])

		if err1 != nil || err2 != nil {
			continue
		}

		cost := 1
		if len(parts) >= 3 {
			if c, err := strconv.Atoi(parts[2]); err == nil {
				cost = c
			}
		}

		edges = append(edges, Edge{from, to, cost})
		nodeSet[from] = true
		nodeSet[to] = true
	}

	return len(nodeSet), edges
}

func main() {
	n, edges := readUSGrid("US-Grid.txt")

	if n == 0 || len(edges) == 0 {
		fmt.Println("No se pudo cargar el archivo")
		return
	}

	fmt.Printf("Edificios: %d\n", n)
	fmt.Printf("Conexiones posibles: %d\n\n", len(edges))

	totalCost, mstEdges, allCost := primMST(n, edges)

	fmt.Printf("Costo total minimo: %d\n", totalCost)
	fmt.Printf("Conexiones a instalar: %d\n\n", len(mstEdges))

	fmt.Println("Lista de conexiones:")
	for _, e := range mstEdges {
		fmt.Printf("%d - %d (costo: %d)\n", e.from, e.to, e.cost)
	}

	fmt.Printf("\nCosto si se conectaran todos contra todos: %d\n", allCost)
	fmt.Printf("Ahorro: %d\n", allCost-totalCost)
}
