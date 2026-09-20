package main

import (
	"flag"
	"fmt"
	"log"
	"map/routing-service/internal/graph"
	"runtime"
	"time"
)

func main() {
	edges := flag.Int("edges", 100000, "number of graph edges")
	updates := flag.Int("updates", 1000, "number of live edge-weight updates")
	flag.Parse()

	if *edges <= 0 || *updates <= 0 || *updates > *edges {
		log.Fatal("edges and updates must be positive, and updates cannot exceed edges")
	}

	input := syntheticGraph(*edges)
	updateSet := make(map[int]float64, *updates)
	for edgeID := 0; edgeID < *updates; edgeID++ {
		updateSet[edgeID] = 10 + float64(edgeID%90)
	}

	fmt.Printf("navigation graph benchmark\nedges=%d\nupdates=%d\n", *edges, *updates)
	for _, architecture := range []graph.Architecture{
		&graph.ContractionHierarchy{},
		&graph.CustomizableContractionHierarchy{},
	} {
		buildStart := time.Now()
		if err := architecture.Build(input); err != nil {
			log.Fatalf("%s build failed: %v", architecture.Name(), err)
		}
		buildDuration := time.Since(buildStart)

		before := runtime.MemStats{}
		runtime.ReadMemStats(&before)
		updateDuration, err := architecture.UpdateWeights(updateSet)
		if err != nil {
			log.Fatalf("%s update failed: %v", architecture.Name(), err)
		}
		after := runtime.MemStats{}
		runtime.ReadMemStats(&after)

		fmt.Printf("\narchitecture=%s\nbuild_ms=%.3f\nupdate_ms=%.3f\nupdate_target_met=%t\nestimated_graph_bytes=%d\nheap_delta_bytes=%d\n",
			architecture.Name(),
			float64(buildDuration.Microseconds())/1000,
			float64(updateDuration.Microseconds())/1000,
			updateDuration < 100*time.Millisecond,
			architecture.MemoryBytes(),
			after.HeapAlloc-before.HeapAlloc,
		)
	}
}

func syntheticGraph(edgeCount int) graph.Graph {
	edges := make([]graph.Edge, edgeCount)
	for edgeID := range edges {
		edges[edgeID] = graph.Edge{
			ID:     edgeID,
			From:   edgeID,
			To:     (edgeID + 1) % edgeCount,
			Weight: 1 + float64(edgeID%60),
			Mode:   "car",
		}
	}
	return graph.Graph{Edges: edges}
}
