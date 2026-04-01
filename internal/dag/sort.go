package dag

import (
	"fmt"
	"sort"
)

// Sort returns node IDs in topological order using Kahn's algorithm.
// Returns ErrCycleDetected if the graph contains a cycle.
func (g *Graph) Sort() ([]string, error) {
	inDegree := make(map[string]int, len(g.nodes))
	for id := range g.nodes {
		inDegree[id] = 0
	}
	for _, targets := range g.edges {
		for _, to := range targets {
			inDegree[to]++
		}
	}

	var queue []string
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}
	sort.Strings(queue)

	var order []string
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		order = append(order, current)

		neighbors := make([]string, len(g.edges[current]))
		copy(neighbors, g.edges[current])
		sort.Strings(neighbors)

		for _, to := range neighbors {
			inDegree[to]--
			if inDegree[to] == 0 {
				queue = append(queue, to)
				sort.Strings(queue)
			}
		}
	}

	if len(order) != len(g.nodes) {
		return nil, fmt.Errorf("topological sort: %w", ErrCycleDetected)
	}
	return order, nil
}

// DetectCycle reports whether the graph contains a cycle.
func (g *Graph) DetectCycle() bool {
	_, err := g.Sort()
	return err != nil
}

// Validate checks that the graph is a valid DAG with no cycles and no orphan edges.
func (g *Graph) Validate() error {
	for from, targets := range g.edges {
		if _, ok := g.nodes[from]; !ok {
			return fmt.Errorf("validating graph: orphan edge from non-existent node %q", from)
		}
		for _, to := range targets {
			if _, ok := g.nodes[to]; !ok {
				return fmt.Errorf("validating graph: orphan edge to non-existent node %q", to)
			}
		}
	}

	if _, err := g.Sort(); err != nil {
		return fmt.Errorf("validating graph: %w", err)
	}

	return nil
}
