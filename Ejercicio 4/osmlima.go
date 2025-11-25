package main

import (
“context”
“encoding/csv”
“fmt”
“math”
“os”

```
"github.com/paulmach/osm"
"github.com/paulmach/osm/osmpbf"
```

)

type Node struct {
ID  int64
Lat float64
Lon float64
}

type Edge struct {
From     int64
To       int64
Distance float64
}

// Calcular distancia haversine entre dos coordenadas
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
const R = 6371000 // Radio de la Tierra en metros

```
dLat := (lat2 - lat1) * math.Pi / 180
dLon := (lon2 - lon1) * math.Pi / 180

a := math.Sin(dLat/2)*math.Sin(dLat/2) +
	math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
	math.Sin(dLon/2)*math.Sin(dLon/2)

c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

return R * c
```

}

func main() {
// Archivo OSM de Perú descargado desde Geofabrik
// https://download.geofabrik.de/south-america/peru-latest.osm.pbf

```
filename := "peru-latest.osm.pbf"
file, err := os.Open(filename)
if err != nil {
	fmt.Printf("Error: No se pudo abrir %s\n", filename)
	fmt.Println("\nDescarga el archivo desde:")
	fmt.Println("https://download.geofabrik.de/south-america/peru-latest.osm.pbf")
	return
}
defer file.Close()

scanner := osmpbf.New(context.Background(), file, 3)
defer scanner.Close()

nodes := make(map[int64]Node)
var edges []Edge

// Bounding box para Lima (ajusta según el distrito que necesites)
// Estas coordenadas cubren aproximadamente el distrito de Miraflores/San Isidro
const (
	minLat = -12.15
	maxLat = -12.05
	minLon = -77.05
	maxLon = -76.95
)

fmt.Println("=== EXTRACTOR DE CALLES DE LIMA - OSM ===\n")
fmt.Printf("Área seleccionada:\n")
fmt.Printf("  Latitud:  %.4f° a %.4f°\n", minLat, maxLat)
fmt.Printf("  Longitud: %.4f° a %.4f°\n\n", minLon, maxLon)
fmt.Println("Procesando archivo OSM... (esto puede tomar varios minutos)")

nodeCount := 0
wayCount := 0

for scanner.Scan() {
	switch obj := scanner.Object().(type) {
	case *osm.Node:
		// Filtrar solo nodos dentro del bounding box de Lima
		if obj.Lat >= minLat && obj.Lat <= maxLat &&
			obj.Lon >= minLon && obj.Lon <= maxLon {
			nodes[obj.ID] = Node{
				ID:  obj.ID,
				Lat: obj.Lat,
				Lon: obj.Lon,
			}
			nodeCount++
			if nodeCount%1000 == 0 {
				fmt.Printf("\rNodos procesados: %d", nodeCount)
			}
		}

	case *osm.Way:
		// Filtrar solo calles (highway tag)
		if highway, ok := obj.Tags.Find("highway"); ok {
			// Tipos de vías que nos interesan
			validTypes := map[string]bool{
				"motorway":     true,
				"trunk":        true,
				"primary":      true,
				"secondary":    true,
				"tertiary":     true,
				"residential":  true,
				"unclassified": true,
				"service":      true,
				"living_street": true,
			}

			if validTypes[highway] {
				// Crear aristas entre nodos consecutivos
				for i := 0; i < len(obj.Nodes)-1; i++ {
					fromID := obj.Nodes[i].ID
					toID := obj.Nodes[i+1].ID

					// Solo si ambos nodos están en nuestra área
					from, fromOk := nodes[fromID]
					to, toOk := nodes[toID]

					if fromOk && toOk {
						distance := haversine(from.Lat, from.Lon, to.Lat, to.Lon)
						edges = append(edges, Edge{
							From:     fromID,
							To:       toID,
							Distance: distance,
						})

						// Si no es one-way, agregar en ambas direcciones
						if oneway, ok := obj.Tags.Find("oneway"); !ok || oneway != "yes" {
							edges = append(edges, Edge{
								From:     toID,
								To:       fromID,
								Distance: distance,
							})
						}
					}
				}
				wayCount++
				if wayCount%500 == 0 {
					fmt.Printf("\rNodos: %d | Calles procesadas: %d", nodeCount, wayCount)
				}
			}
		}
	}
}

if err := scanner.Err(); err != nil {
	fmt.Printf("\nError durante procesamiento: %v\n", err)
	return
}

fmt.Printf("\n\n✅ Extracción completada:\n")
fmt.Printf("   • Nodos (intersecciones): %d\n", len(nodes))
fmt.Printf("   • Aristas (calles): %d\n", len(edges))
fmt.Printf("   • Calles procesadas: %d\n\n", wayCount)

// Exportar a CSV
fmt.Println("Exportando a archivos CSV...")
if err := exportNodesToCSV(nodes, "lima_nodes.csv"); err != nil {
	fmt.Printf("Error exportando nodos: %v\n", err)
	return
}
if err := exportEdgesToCSV(edges, "lima_edges.csv"); err != nil {
	fmt.Printf("Error exportando aristas: %v\n", err)
	return
}

fmt.Println("\n✅ Archivos exportados exitosamente:")
fmt.Println("   📄 lima_nodes.csv")
fmt.Println("   📄 lima_edges.csv")
fmt.Println("\n💡 Ahora puedes ejecutar el programa de navegación GPS:")
fmt.Println("   go run main.go")
```

}

func exportNodesToCSV(nodes map[int64]Node, filename string) error {
file, err := os.Create(filename)
if err != nil {
return err
}
defer file.Close()

```
writer := csv.NewWriter(file)
defer writer.Flush()

// Header
writer.Write([]string{"id", "lat", "lon"})

count := 0
for _, node := range nodes {
	writer.Write([]string{
		fmt.Sprintf("%d", node.ID),
		fmt.Sprintf("%.6f", node.Lat),
		fmt.Sprintf("%.6f", node.Lon),
	})
	count++
}

fmt.Printf("   • %d nodos exportados a %s\n", count, filename)
return nil
```

}

func exportEdgesToCSV(edges []Edge, filename string) error {
file, err := os.Create(filename)
if err != nil {
return err
}
defer file.Close()

```
writer := csv.NewWriter(file)
defer writer.Flush()

// Header
writer.Write([]string{"from", "to", "distance"})

for _, edge := range edges {
	writer.Write([]string{
		fmt.Sprintf("%d", edge.From),
		fmt.Sprintf("%d", edge.To),
		fmt.Sprintf("%.2f", edge.Distance),
	})
}

fmt.Printf("   • %d aristas exportadas a %s\n", len(edges), filename)
return nil
```

}