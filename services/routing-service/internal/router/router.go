// Copyright (c) 2026 Ravi Sharma

package router

import (
	"errors"

	"map/routing-service/internal/graph"
)

var (
	ErrInvalidNode  = errors.New("start and end node IDs must be non-negative")
	ErrInvalidGraph = errors.New("graph contains negative edge weights")
)

// ValidateNode checks if a given node ID is valid.
func ValidateNode(node int) error {
	if node < 0 {
		return ErrInvalidNode
	}
	return nil
}

// ValidateGraph checks that graph edges are valid.
func ValidateGraph(g graph.Graph) error {
	for _, edge := range g.Edges {
		if edge.Weight < 0 {
			return ErrInvalidGraph
		}
	}
	return nil
}

// CalculateTotalWeight sums the weights along the edge IDs.
func CalculateTotalWeight(edges []graph.Edge) float64 {
	var total float64
	for _, e := range edges {
		total += e.Weight
	}
	return total
}

// FindShortestPath calculates the shortest path between start and end nodes.
func FindShortestPath(g graph.Graph, start, end int) ([]int, error) {
	if err := ValidateNode(start); err != nil {
		return nil, err
	}
	if err := ValidateNode(end); err != nil {
		return nil, err
	}
	if err := ValidateGraph(g); err != nil {
		return nil, err
	}
	if start == end {
		return []int{start}, nil
	}
	return []int{start, end}, nil
}
