package graph

// A node within a graph
type Node struct {
	ID       string
	Outbound []Edge
	Inbound  []Edge
}

// An edge in the graph
type Edge struct {
	Start, End *Node
}

func (n *Node) AddEdgeTo(other *Node) {
	edge := Edge{n, other}

	n.Outbound = append(n.Outbound, edge)
	other.Inbound = append(other.Inbound, edge)
}

func NewNode(id string) (res Node) {
	res.ID = id
	res.Outbound = make([]Edge, 0)
	res.Inbound = make([]Edge, 0)
	return
}
