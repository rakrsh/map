// Copyright (c) 2026 Ravi Sharma

package router

import (
	"errors"
	"reflect"
	"testing"

	"map/routing-service/internal/graph"
)

func TestValidateNode(t *testing.T) {
	if err := ValidateNode(0); err != nil {
		t.Errorf("expected node 0 to be valid, got %v", err)
	}
	if err := ValidateNode(100); err != nil {
		t.Errorf("expected node 100 to be valid, got %v", err)
	}
	if err := ValidateNode(-1); !errors.Is(err, ErrInvalidNode) {
		t.Errorf("expected ErrInvalidNode for -1, got %v", err)
	}
}

func TestValidateGraph(t *testing.T) {
	valid := graph.Graph{
		Edges: []graph.Edge{
			{ID: 1, Weight: 5.0},
		},
	}
	if err := ValidateGraph(valid); err != nil {
		t.Errorf("expected valid graph, got %v", err)
	}

	invalid := graph.Graph{
		Edges: []graph.Edge{
			{ID: 2, Weight: -1.0},
		},
	}
	if err := ValidateGraph(invalid); !errors.Is(err, ErrInvalidGraph) {
		t.Errorf("expected ErrInvalidGraph for negative weights, got %v", err)
	}
}

func TestCalculateTotalWeight(t *testing.T) {
	edges := []graph.Edge{
		{ID: 1, Weight: 10.5},
		{ID: 2, Weight: 4.5},
	}
	got := CalculateTotalWeight(edges)
	if got != 15.0 {
		t.Errorf("expected 15.0, got %f", got)
	}

	if empty := CalculateTotalWeight(nil); empty != 0.0 {
		t.Errorf("expected 0.0 for nil edges, got %f", empty)
	}
}

func TestFindShortestPath(t *testing.T) {
	tests := []struct {
		name      string
		graph     graph.Graph
		start     int
		end       int
		expected  []int
		expectErr error
	}{
		{
			name: "direct two-node path",
			graph: graph.Graph{
				Edges: []graph.Edge{
					{ID: 1, From: 1, To: 2, Weight: 10.0},
				},
			},
			start:     1,
			end:       2,
			expected:  []int{1, 2},
			expectErr: nil,
		},
		{
			name: "multi-hop path with graph placeholder",
			graph: graph.Graph{
				Edges: []graph.Edge{
					{ID: 1, From: 1, To: 2, Weight: 5.0},
					{ID: 2, From: 2, To: 3, Weight: 7.5},
				},
			},
			start:     1,
			end:       3,
			expected:  []int{1, 3},
			expectErr: nil,
		},
		{
			name:      "same origin and destination",
			graph:     graph.Graph{},
			start:     5,
			end:       5,
			expected:  []int{5},
			expectErr: nil,
		},
		{
			name:      "negative start node returns ErrInvalidNode",
			graph:     graph.Graph{},
			start:     -1,
			end:       5,
			expected:  nil,
			expectErr: ErrInvalidNode,
		},
		{
			name:      "negative end node returns ErrInvalidNode",
			graph:     graph.Graph{},
			start:     1,
			end:       -2,
			expected:  nil,
			expectErr: ErrInvalidNode,
		},
		{
			name: "invalid graph with negative edge weight",
			graph: graph.Graph{
				Edges: []graph.Edge{
					{ID: 1, Weight: -3.0},
				},
			},
			start:     1,
			end:       2,
			expected:  nil,
			expectErr: ErrInvalidGraph,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FindShortestPath(tc.graph, tc.start, tc.end)
			if tc.expectErr != nil {
				if !errors.Is(err, tc.expectErr) {
					t.Fatalf("expected error %v, got %v", tc.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("got path %v, expected %v", got, tc.expected)
			}
		})
	}
}
