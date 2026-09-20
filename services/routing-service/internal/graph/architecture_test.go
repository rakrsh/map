package graph

import "testing"

func TestArchitecturesApplyLiveWeights(t *testing.T) {
	input := Graph{Edges: []Edge{
		{ID: 0, Weight: 10},
		{ID: 1, Weight: 20},
		{ID: 2, Weight: 30},
	}}
	updates := map[int]float64{1: 42}

	architectures := []Architecture{
		&ContractionHierarchy{},
		&CustomizableContractionHierarchy{},
	}
	for _, architecture := range architectures {
		if err := architecture.Build(input); err != nil {
			t.Fatalf("%s build failed: %v", architecture.Name(), err)
		}
		if _, err := architecture.UpdateWeights(updates); err != nil {
			t.Fatalf("%s update failed: %v", architecture.Name(), err)
		}
		if weight, ok := architecture.Weight(1); !ok || weight != 42 {
			t.Fatalf("%s did not apply updated weight: got %v, exists %v", architecture.Name(), weight, ok)
		}
		if architecture.MemoryBytes() == 0 {
			t.Fatalf("%s reported zero memory", architecture.Name())
		}
	}
}

func TestArchitecturesRejectUnknownEdges(t *testing.T) {
	input := Graph{Edges: []Edge{{ID: 0, Weight: 1}}}
	for _, architecture := range []Architecture{&ContractionHierarchy{}, &CustomizableContractionHierarchy{}} {
		if err := architecture.Build(input); err != nil {
			t.Fatal(err)
		}
		if _, err := architecture.UpdateWeights(map[int]float64{99: 2}); err == nil {
			t.Fatalf("%s accepted an unknown edge", architecture.Name())
		}
	}
}
