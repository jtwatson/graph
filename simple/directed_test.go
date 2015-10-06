// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple

import (
	"math"
	"testing"

	"github.com/gonum/graph"
)

var _ graph.Graph = &DirectedGraph{}
var _ graph.Directed = &DirectedGraph{}
var _ graph.Directed = &DirectedGraph{}

// Tests Issue #27
func TestEdgeOvercounting(t *testing.T) {
	g := generateDummyGraph()

	if neigh := g.From(Node(Node(2))); len(neigh) != 2 {
		t.Errorf("Node 2 has incorrect number of neighbors got neighbors %v (count %d), expected 2 neighbors {0,1}", neigh, len(neigh))
	}
}

func generateDummyGraph() *DirectedGraph {
	nodes := [4]struct{ srcId, targetId int }{
		{2, 1},
		{1, 0},
		{2, 0},
		{0, 2},
	}

	g := NewDirectedGraph(0, math.Inf(1))

	for _, n := range nodes {
		g.SetEdge(Edge{F: Node(n.srcId), T: Node(n.targetId), W: 1})
	}

	return g
}

// Test for issue #123 https://github.com/gonum/graph/issues/123
func TestIssue123DirectedGraph(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
	g := NewDirectedGraph(0, math.Inf(1))

	n0 := Node(g.NewNodeID())
	g.AddNode(n0)

	n1 := Node(g.NewNodeID())
	g.AddNode(n1)

	g.RemoveNode(n0)

	n2 := Node(g.NewNodeID())
	g.AddNode(n2)
}

func TestDirectedDegree(t *testing.T) {
	g := generateDummyGraph()

	if degree := g.Degree(Node(0)); degree != 3 {
		t.Errorf("Node 0 has incorrect Degree, got Degree of %d, expected 3\n", degree)
	}

	if degree := g.Degree(Node(1)); degree != 2 {
		t.Errorf("Node 1 has incorrect Degree, got Degree of %d, expected 2\n", degree)
	}

	if degree := g.Degree(Node(2)); degree != 3 {
		t.Errorf("Node 2 has incorrect Degree, got Degree of %d, expected 3\n", degree)
	}

	if degree := g.Degree(Node(3)); degree != 0 {
		t.Errorf("Node 3 has incorrect Degree, got Degree of %d, expected 0\n", degree)
	}

}

func TestInDegree(t *testing.T) {
	g := generateDummyGraph()

	if degree := g.InDegree(Node(0)); degree != 2 {
		t.Errorf("Node 0 has incorrect InDegree, got InDegree of %d, expected 2\n", degree)
	}

	if degree := g.InDegree(Node(1)); degree != 1 {
		t.Errorf("Node 1 has incorrect InDegree, got InDegree of %d, expected 1\n", degree)
	}

	if degree := g.InDegree(Node(2)); degree != 1 {
		t.Errorf("Node 2 has incorrect InDegree, got InDegree of %d, expected 1\n", degree)
	}

	if degree := g.InDegree(Node(3)); degree != 0 {
		t.Errorf("Node 3 has incorrect InDegree, got InDegree of %d, expected 0\n", degree)
	}

}

func TestOutDegree(t *testing.T) {
	g := generateDummyGraph()

	if degree := g.OutDegree(Node(0)); degree != 1 {
		t.Errorf("Node 0 has incorrect OutDegree, got OutDegree of %d, expected 1\n", degree)
	}

	if degree := g.OutDegree(Node(1)); degree != 1 {
		t.Errorf("Node 1 has incorrect OutDegree, got OutDegree of %d, expected 1\n", degree)
	}

	if degree := g.OutDegree(Node(2)); degree != 2 {
		t.Errorf("Node 2 has incorrect OutDegree, got OutDegree of %d, expected 2\n", degree)
	}

	if degree := g.OutDegree(Node(3)); degree != 0 {
		t.Errorf("Node 3 has incorrect OutDegree, got OutDegree of %d, expected 0\n", degree)
	}

}
