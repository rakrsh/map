// Copyright (c) 2026 Ravi Sharma

package graph

import (
	"fmt"
	"time"
)

// Architecture is the minimal contract needed to compare live traffic updates.
type Architecture interface {
	Name() string
	Build(Graph) error
	UpdateWeights(map[int]float64) (time.Duration, error)
	MemoryBytes() uint64
	Weight(edgeID int) (float64, bool)
}

// ContractionHierarchy models a preprocessed hierarchy whose live updates must
// invalidate and refresh dependent shortcut weights.
type ContractionHierarchy struct {
	weights   map[int]float64
	shortcuts map[int][]int
}

func (engine *ContractionHierarchy) Name() string { return "CH" }

func (engine *ContractionHierarchy) Build(input Graph) error {
	engine.weights = make(map[int]float64, len(input.Edges))
	engine.shortcuts = make(map[int][]int, len(input.Edges))
	for _, edge := range input.Edges {
		if edge.ID < 0 {
			return fmt.Errorf("edge ID must be non-negative: %d", edge.ID)
		}
		engine.weights[edge.ID] = edge.Weight
		// A deterministic synthetic shortcut dependency. A production CH build
		// would populate this from the contracted topology.
		engine.shortcuts[edge.ID] = []int{edge.ID, (edge.ID + 1) % max(1, len(input.Edges))}
	}
	return nil
}

func (engine *ContractionHierarchy) UpdateWeights(updates map[int]float64) (time.Duration, error) {
	start := time.Now()
	for edgeID, weight := range updates {
		if _, exists := engine.weights[edgeID]; !exists {
			return 0, fmt.Errorf("unknown edge ID: %d", edgeID)
		}
		engine.weights[edgeID] = weight
		for _, dependentID := range engine.shortcuts[edgeID] {
			engine.weights[dependentID] = weight
		}
	}
	return time.Since(start), nil
}

func (engine *ContractionHierarchy) MemoryBytes() uint64 {
	return uint64(len(engine.weights))*24 + uint64(len(engine.shortcuts))*32
}

func (engine *ContractionHierarchy) Weight(edgeID int) (float64, bool) {
	weight, exists := engine.weights[edgeID]
	return weight, exists
}

// CustomizableContractionHierarchy models a topology that is built once and
// customized by applying live edge weights without rebuilding the hierarchy.
type CustomizableContractionHierarchy struct {
	weights []float64
}

func (engine *CustomizableContractionHierarchy) Name() string { return "CCH" }

func (engine *CustomizableContractionHierarchy) Build(input Graph) error {
	engine.weights = make([]float64, 0, len(input.Edges))
	for _, edge := range input.Edges {
		if edge.ID < 0 {
			return fmt.Errorf("edge ID must be non-negative: %d", edge.ID)
		}
		if edge.ID >= len(engine.weights) {
			padding := make([]float64, edge.ID-len(engine.weights)+1)
			engine.weights = append(engine.weights, padding...)
		}
		engine.weights[edge.ID] = edge.Weight
	}
	return nil
}

func (engine *CustomizableContractionHierarchy) UpdateWeights(updates map[int]float64) (time.Duration, error) {
	start := time.Now()
	for edgeID, weight := range updates {
		if edgeID < 0 || edgeID >= len(engine.weights) {
			return 0, fmt.Errorf("unknown edge ID: %d", edgeID)
		}
		engine.weights[edgeID] = weight
	}
	return time.Since(start), nil
}

func (engine *CustomizableContractionHierarchy) MemoryBytes() uint64 {
	return uint64(len(engine.weights)) * 8
}

func (engine *CustomizableContractionHierarchy) Weight(edgeID int) (float64, bool) {
	if edgeID < 0 || edgeID >= len(engine.weights) {
		return 0, false
	}
	return engine.weights[edgeID], true
}

func max(left, right int) int {
	if left > right {
		return left
	}
	return right
}
