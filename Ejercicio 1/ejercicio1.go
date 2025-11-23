package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type graph struct {
	adyacencia map[string][]string
}

func NuevoGrafo() *graph {
	return &graph{
		adyacencia: make(map[string][]string),
	}
}

func (g *graph) AgregarArista(x, y string) { //amistades bidireccionales
	//el segundo nodo se agrega a los amigos del primero
	g.adyacencia[x] = append(g.adyacencia[x], y)
	//y viceversa (xq si uno es amigo del otro, el otro es amigo del uno tmb)
	g.adyacencia[y] = append(g.adyacencia[y], x)
}

// para encontrar la N separacion
func (g *graph) NSeparacion(origen string, n int) []string {
	if n < 0 {
		//distancia negativa no se puede
		return []string{}
	}

	vistos := make(map[string]bool)   //guardar los usuarios q ya se vieron
	distancia := make(map[string]int) //va guardando la distancia q recorre
	queue := []string{origen}         //cola FIFO para bfs

	//inicia
	vistos[origen] = true
	distancia[origen] = 0

	for len(queue) > 0 {
		actual := queue[0]
		queue = queue[1:] //el primero se saca

		if distancia[actual] == n { //si la distancia actual es igual a n
			continue //no sigue, salta al otro
		}

		for _, vecino := range g.adyacencia[actual] {
			if !vistos[vecino] { //si no vio al usuario
				vistos[vecino] = true //visto
				distancia[vecino] = distancia[actual] + 1
				queue = append(queue, vecino)
			}
		}
	}

	result := []string{}
	for nodo, dist := range distancia {
		if dist == n { //si la distancia es IGUAL a n
			result = append(result, nodo)
		}
	}
	return result
}

func CargarGrafo(filename string) (*graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error abriendo: %v", err)
	}
	defer file.Close()

	g := NuevoGrafo()
	scanner := bufio.NewScanner(file)

	idToLabel := make(map[string]string)

	var source, target string
	var idActual, labelActual string
	enNodo := false
	enArista := false

	for scanner.Scan() {
		linea := strings.TrimSpace(scanner.Text())

		if linea == "node" || strings.HasPrefix(linea, "node") {
			enNodo = true
			idActual = ""
			labelActual = ""
			continue
		}

		if enNodo {
			if strings.Contains(linea, "id") {
				parts := strings.Fields(linea)
				if len(parts) >= 2 {
					idActual = parts[1]
				}
				continue
			}
			if strings.Contains(linea, "label") {
				parts := strings.Fields(linea)
				if len(parts) >= 2 {
					labelActual = strings.Trim(parts[1], "\"")
				}
				continue
			}
			if linea == "]" {
				if idActual != "" {
					if labelActual != "" {
						idToLabel[idActual] = labelActual
					} else {
						idToLabel[idActual] = idActual
					}
				}
				enNodo = false
				continue
			}
		}

		if linea == "edge" || strings.HasPrefix(linea, "edge") {
			enArista = true
			source = ""
			target = ""
			continue
		}

		if enArista {
			if strings.Contains(linea, "id") {
				continue //IGNORA EL ID PORFAVOR AAAA
			}
			if strings.Contains(linea, "source") {
				parts := strings.Fields(linea)
				if len(parts) >= 2 {
					source = parts[1]
				}
				continue
			}

			if strings.Contains(linea, "target") {
				parts := strings.Fields(linea)
				if len(parts) >= 2 {
					target = parts[1]
				}
				continue
			}
			if linea == "]" {
				if source != "" && target != "" {
					sourceLabel := idToLabel[source]
					targetLabel := idToLabel[target]

					if sourceLabel == "" {
						sourceLabel = source
					}
					if targetLabel == "" {
						targetLabel = target
					}
					g.AgregarArista(sourceLabel, targetLabel)
				}
				enArista = false
				continue
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error leyendo archivo: %v", err)
	}

	// Verificar que se cargó algo
	if len(g.adyacencia) == 0 {
		return nil, fmt.Errorf("no se encontraron aristas en el archivo")
	}

	return g, nil
}

func (g *graph) MostrarNodos() []string {
	nodos := []string{}
	for nodo := range g.adyacencia {
		nodos = append(nodos, nodo)
	}
	sort.Strings(nodos)
	return nodos
}

func (g *graph) ValidarExiste(nodo string) bool {
	_, existe := g.adyacencia[nodo]
	return existe
}

func main() {
	graph, err := CargarGrafo("karate.gml")
	if err != nil {
		fmt.Print("Error cargando el grafo")
	}
	nodos := graph.MostrarNodos()
	fmt.Printf("Nodos en el grafo: %v\n", nodos)

	for { //pregunta infinitamente :D
		fmt.Println("_______________________________________")
		fmt.Println("Escriba 'SALIR' para terminar")
		var usuario string
		fmt.Print("Ingrese el usuario inicial: ")
		fmt.Scanln(&usuario)

		if usuario == "SALIR" {
			fmt.Println("Byee! :D")
			break
		}

		if !graph.ValidarExiste(usuario) {
			fmt.Println("Ese usuario no existe en el grafo :|")
			continue //reiniciar pregunta
		}

		var n int
		fmt.Print("Ingrese el grado de separacion (N): ")
		_, err := fmt.Scanln(&n)
		if err != nil || n < 0 {
			fmt.Printf("%v no es un numero valido >:/\n", n)
			continue
		}

		resultado := graph.NSeparacion(usuario, n)

		fmt.Println("Un momento... Listo! :D")
		fmt.Printf("Usuarios a %d separacion de %s: %d\n", n, usuario, len(resultado))
		if len(resultado) == 0 {
			fmt.Println("No hay usuarios a esa distancia")
		} else {
			fmt.Printf("Usuarios: %v\n", resultado)
		}

	}
}
