package graph

// Graph represents a lightweight routing graph model used for pathfinding.
type Graph struct {
	Nodes []Node
	Edges []Edge
}

type Node struct {
	ID       int
	Lat      float64
	Lng      float64
	Kind     string
	Metadata map[string]interface{}
}

type Edge struct {
	ID       int
	From     int
	To       int
	Weight   float64
	Mode     string
	Geometry []Node
	Metadata map[string]interface{}
}
