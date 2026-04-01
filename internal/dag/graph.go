package dag

import (
	"errors"
	"fmt"
	"sort"
)

var ErrCycleDetected = errors.New("cycle detected in dependency graph")

// Node represents a single vertex in the dependency graph.
type Node struct {
	ID     string
	Table  string
	Column string
}

// Edge represents a directed edge from one node to another.
type Edge struct {
	From string
	To   string
}

// Graph is a directed acyclic graph tracking column-level dependencies.
type Graph struct {
	nodes   map[string]*Node
	edges   map[string][]string
	reverse map[string][]string
}

// New creates an empty Graph.
func New() *Graph {
	return &Graph{
		nodes:   make(map[string]*Node),
		edges:   make(map[string][]string),
		reverse: make(map[string][]string),
	}
}

// AddNode inserts a node into the graph. Returns an error if a node with the
// same ID already exists.
func (g *Graph) AddNode(id, table, column string) error {
	if _, exists := g.nodes[id]; exists {
		return fmt.Errorf("adding node %q: %w", id, errors.New("duplicate node"))
	}
	g.nodes[id] = &Node{ID: id, Table: table, Column: column}
	return nil
}

// AddEdge creates a directed edge from one node to another. Returns an error
// if either node does not exist.
func (g *Graph) AddEdge(from, to string) error {
	if _, ok := g.nodes[from]; !ok {
		return fmt.Errorf("adding edge: %w", fmt.Errorf("node %q not found", from))
	}
	if _, ok := g.nodes[to]; !ok {
		return fmt.Errorf("adding edge: %w", fmt.Errorf("node %q not found", to))
	}
	for _, existing := range g.edges[from] {
		if existing == to {
			return nil
		}
	}
	g.edges[from] = append(g.edges[from], to)
	g.reverse[to] = append(g.reverse[to], from)
	return nil
}

// Node returns the node with the given ID, if it exists.
func (g *Graph) Node(id string) (*Node, bool) {
	n, ok := g.nodes[id]
	return n, ok
}

// Nodes returns all nodes in the graph in sorted order by ID.
func (g *Graph) Nodes() []*Node {
	result := make([]*Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		result = append(result, n)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

// Dependents returns the direct downstream dependents of a node.
func (g *Graph) Dependents(id string) []string {
	out := make([]string, len(g.edges[id]))
	copy(out, g.edges[id])
	sort.Strings(out)
	return out
}

// Dependencies returns the direct upstream dependencies of a node.
func (g *Graph) Dependencies(id string) []string {
	out := make([]string, len(g.reverse[id]))
	copy(out, g.reverse[id])
	sort.Strings(out)
	return out
}

// AllDependents returns all transitive downstream dependents using BFS.
func (g *Graph) AllDependents(id string) []string {
	return g.bfs(id, g.edges)
}

// AllDependencies returns all transitive upstream dependencies using BFS.
func (g *Graph) AllDependencies(id string) []string {
	return g.bfs(id, g.reverse)
}

func (g *Graph) bfs(start string, adj map[string][]string) []string {
	visited := make(map[string]bool)
	queue := []string{start}
	visited[start] = true
	var result []string

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, neighbor := range adj[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				result = append(result, neighbor)
				queue = append(queue, neighbor)
			}
		}
	}
	sort.Strings(result)
	return result
}

// HasNode reports whether the graph contains a node with the given ID.
func (g *Graph) HasNode(id string) bool {
	_, ok := g.nodes[id]
	return ok
}

// HasEdge reports whether a directed edge exists from one node to another.
func (g *Graph) HasEdge(from, to string) bool {
	for _, t := range g.edges[from] {
		if t == to {
			return true
		}
	}
	return false
}

// RemoveNode deletes a node and all of its associated edges from the graph.
func (g *Graph) RemoveNode(id string) {
	if _, ok := g.nodes[id]; !ok {
		return
	}

	for _, to := range g.edges[id] {
		g.reverse[to] = removeFromSlice(g.reverse[to], id)
	}
	delete(g.edges, id)

	for _, from := range g.reverse[id] {
		g.edges[from] = removeFromSlice(g.edges[from], id)
	}
	delete(g.reverse, id)

	delete(g.nodes, id)
}

func removeFromSlice(s []string, val string) []string {
	result := s[:0]
	for _, v := range s {
		if v != val {
			result = append(result, v)
		}
	}
	return result
}
