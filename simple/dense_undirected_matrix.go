// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple

import (
	"github.com/gonum/graph"
	"github.com/gonum/matrix/mat64"
)

// UndirectedDenseGraph represents a graph such that all IDs are in a contiguous
// block from 0 to n-1.
type UndirectedDenseGraph struct {
	self   float64
	absent float64
	mat    *mat64.SymDense
}

// NewUndirectedDenseGraph creates an undirected dense graph with n nodes.
// If passable is true all pairs of nodes will be connected by an edge
// with unit cost, otherwise every node will start unconnected with
// the cost specified by absent. The self parameter specifies the cost
// of self connection.
func NewUndirectedDenseGraph(n int, passable bool, self, absent float64) *UndirectedDenseGraph {
	mat := make([]float64, n*n)
	v := 1.
	if !passable {
		v = absent
	}
	for i := range mat {
		mat[i] = v
	}
	return &UndirectedDenseGraph{
		mat:    mat64.NewSymDense(n, mat),
		self:   self,
		absent: absent,
	}
}

func (g *UndirectedDenseGraph) Has(n graph.Node) bool {
	return g.has(n.ID())
}

func (g *UndirectedDenseGraph) has(id int) bool {
	r := g.mat.Symmetric()
	return 0 <= id && id < r
}

func (g *UndirectedDenseGraph) Nodes() []graph.Node {
	r := g.mat.Symmetric()
	nodes := make([]graph.Node, r)
	for i := 0; i < r; i++ {
		nodes[i] = Node(i)
	}
	return nodes
}

func (g *UndirectedDenseGraph) Edges() []graph.Edge {
	var edges []graph.Edge
	r, _ := g.mat.Dims()
	for i := 0; i < r; i++ {
		for j := i + 1; j < r; j++ {
			if w := g.mat.At(i, j); !isSame(w, g.absent) {
				edges = append(edges, Edge{F: Node(i), T: Node(j), W: w})
			}
		}
	}
	return edges
}

func (g *UndirectedDenseGraph) Degree(n graph.Node) int {
	id := n.ID()
	var deg int
	r := g.mat.Symmetric()
	for i := 0; i < r; i++ {
		if i == id {
			continue
		}
		if !isSame(g.mat.At(id, i), g.absent) {
			deg++
		}
	}
	return deg
}

func (g *UndirectedDenseGraph) From(n graph.Node) []graph.Node {
	var neighbors []graph.Node
	id := n.ID()
	r := g.mat.Symmetric()
	for i := 0; i < r; i++ {
		if i == id {
			continue
		}
		if !isSame(g.mat.At(id, i), g.absent) {
			neighbors = append(neighbors, Node(i))
		}
	}
	return neighbors
}

func (g *UndirectedDenseGraph) HasEdgeBetween(u, v graph.Node) bool {
	uid := u.ID()
	vid := v.ID()
	return uid != vid && !isSame(g.mat.At(uid, vid), g.absent)
}

func (g *UndirectedDenseGraph) Edge(u, v graph.Node) graph.Edge {
	return g.EdgeBetween(u, v)
}

func (g *UndirectedDenseGraph) EdgeBetween(u, v graph.Node) graph.Edge {
	if g.HasEdgeBetween(u, v) {
		return Edge{F: u, T: v, W: g.mat.At(u.ID(), v.ID())}
	}
	return nil
}

func (g *UndirectedDenseGraph) Weight(x, y graph.Node) (w float64, ok bool) {
	xid := x.ID()
	yid := y.ID()
	if xid == yid {
		return g.self, true
	}
	if g.has(xid) && g.has(yid) {
		return g.mat.At(xid, yid), true
	}
	return g.absent, false
}

func (g *UndirectedDenseGraph) SetEdgeWeight(e graph.Edge) {
	fid := e.From().ID()
	tid := e.To().ID()
	if fid == tid {
		panic("simple: set edge cost of illegal edge")
	}
	g.mat.SetSym(fid, tid, e.Weight())
}

func (g *UndirectedDenseGraph) RemoveEdge(e graph.Edge) {
	g.mat.SetSym(e.From().ID(), e.To().ID(), g.absent)
}

func (g *UndirectedDenseGraph) Matrix() mat64.Matrix {
	// Prevent alteration of dimensions of the returned matrix.
	m := *g.mat
	return &m
}
