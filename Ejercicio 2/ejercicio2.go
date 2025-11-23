package main

import "fmt"

// detectarCiclo detecta si un grafo dirigido contiene ciclos usando DFS con backtracking
func detectarCiclo(numNodos int, aristas [][2]int) (bool, []int) {
	// Construir lista de adyacencia
	grafo := make(map[int][]int)
	for _, arista := range aristas {
		u, v := arista[0], arista[1]
		grafo[u] = append(grafo[u], v)
	}
	
	// Estados: 0 = no visitado, 1 = en proceso (gris), 2 = completado (negro)
	estado := make(map[int]int)
	padre := make(map[int]int)
	
	// Inicializar todos los nodos como no visitados
	for i := 0; i < numNodos; i++ {
		estado[i] = 0
	}
	
	var ciclo []int
	var hayCiclo bool
	
	// Función DFS recursiva
	var dfs func(nodo int) bool
	dfs = func(nodo int) bool {
		// Marcar nodo como "en proceso"
		estado[nodo] = 1
		
		// Explorar vecinos
		for _, vecino := range grafo[nodo] {
			if estado[vecino] == 0 {
				// Nodo no visitado
				padre[vecino] = nodo
				if dfs(vecino) {
					return true
				}
			} else if estado[vecino] == 1 {
				// Nodo en proceso = encontramos un ciclo
				// Reconstruir el ciclo
				ciclo = []int{vecino}
				actual := nodo
				for actual != vecino {
					ciclo = append([]int{actual}, ciclo...)
					actual = padre[actual]
				}
				ciclo = append([]int{vecino}, ciclo...)
				return true
			}
			// Si estado[vecino] == 2, ya fue procesado, no hay ciclo por esa rama
		}
		
		// Marcar nodo como completado
		estado[nodo] = 2
		return false
	}
	
	// Ejecutar DFS desde cada nodo no visitado
	for i := 0; i < numNodos; i++ {
		if estado[i] == 0 {
			if dfs(i) {
				hayCiclo = true
				break
			}
		}
	}
	
	return hayCiclo, ciclo
}

func main() {
	// Ejemplo 1: Grafo con ciclo
	fmt.Println("--- Ejemplo 1: Grafo con ciclo ---")
	numNodos1 := 4
	aristas1 := [][2]int{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 1}, // Ciclo: 1 -> 2 -> 3 -> 1
	}
	
	hayCiclo1, ciclo1 := detectarCiclo(numNodos1, aristas1)
	fmt.Printf("¿Hay ciclo? %v\n", hayCiclo1)
	if hayCiclo1 {
		fmt.Printf("Nodos del ciclo: %v\n", ciclo1)
	}
	
	fmt.Println()
	
	
	fmt.Println("--- Ejemplo 2: Grafo sin ciclo (DAG) ---")
	numNodos2 := 4
	aristas2 := [][2]int{
		{0, 1},
		{0, 2},
		{1, 3},
		{2, 3},
	}
	
	hayCiclo2, ciclo2 := detectarCiclo(numNodos2, aristas2)
	fmt.Printf("¿Hay ciclo? %v\n", hayCiclo2)
	if hayCiclo2 {
		fmt.Printf("Nodos del ciclo: %v\n", ciclo2)
	}
	
	fmt.Println()
	
	
	fmt.Println("--- Ejemplo 3: Grafo con self-loop ---")
	numNodos3 := 3
	aristas3 := [][2]int{
		{0, 1},
		{1, 1}, 
		{1, 2},
	}
	
	hayCiclo3, ciclo3 := detectarCiclo(numNodos3, aristas3)
	fmt.Printf("¿Hay ciclo? %v\n", hayCiclo3)
	if hayCiclo3 {
		fmt.Printf("Nodos del ciclo: %v\n", ciclo3)
	}
}
