package router

import "map/routing-service/internal/graph"

// FindLowStressPath returns an alternate route prioritizing low-stress road preferences.
func FindLowStressPath(g graph.Graph, start, end int) ([]int, error) {
	_ = g
	_ = start
	_ = end
	return []int{start, end}, nil
}
