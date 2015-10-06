// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple

import (
	"fmt"

	"github.com/gonum/graph"
)

// A Directed graph is a highly generalized MutableDirectedGraph.
//
// In most cases it's likely more desireable to use a graph specific to your
// problem domain.
type DirectedGraph struct {
	successors   map[int]map[int]graph.Edge
	predecessors map[int]map[int]graph.Edge
	nodeMap      map[int]graph.Node

	self, absent float64

	// Add/remove convenience variables
	maxID   int
	freeMap map[int]struct{}
}

func NewDirectedGraph(self, absent float64) *DirectedGraph {
	return &DirectedGraph{
		successors:   make(map[int]map[int]graph.Edge),
		predecessors: make(map[int]map[int]graph.Edge),
		nodeMap:      make(map[int]graph.Node),

		self:   self,
		absent: absent,

		maxID:   0,
		freeMap: make(map[int]struct{}),
	}
}

func (g *DirectedGraph) NewNodeID() int {
	if g.maxID != maxInt {
		g.maxID++
		return g.maxID
	}

	// Implicitly checks if len(g.freeMap) == 0
	for id := range g.freeMap {
		return id
	}

	// I cannot foresee this ever happening, but just in case
	if len(g.nodeMap) == maxInt {
		panic("cannot allocate node: graph too large")
	}

	for i := 0; i < maxInt; i++ {
		if _, ok := g.nodeMap[i]; !ok {
			return i
		}
	}

	// Should not happen.
	panic("cannot allocate node id: no free id found")
}

// Adds a node to the graph. Implementation note: if you add a node close to or at
// the max int on your machine NewNode will become slower.
func (g *DirectedGraph) AddNode(n graph.Node) {
	if _, exists := g.nodeMap[n.ID()]; exists {
		panic(fmt.Sprintf("simple: node ID collision: %d", n.ID()))
	}
	g.nodeMap[n.ID()] = n
	g.successors[n.ID()] = make(map[int]graph.Edge)
	g.predecessors[n.ID()] = make(map[int]graph.Edge)

	delete(g.freeMap, n.ID())
	g.maxID = max(g.maxID, n.ID())
}

func (g *DirectedGraph) SetEdge(e graph.Edge) {
	var (
		from = e.From()
		fid  = from.ID()
		to   = e.To()
		tid  = to.ID()
	)

	if fid == tid {
		panic("simple: adding self edge")
	}

	if !g.Has(from) {
		g.AddNode(from)
	}

	if !g.Has(to) {
		g.AddNode(to)
	}

	g.successors[fid][tid] = e
	g.predecessors[tid][fid] = e
}

func (g *DirectedGraph) RemoveNode(n graph.Node) {
	if _, ok := g.nodeMap[n.ID()]; !ok {
		return
	}
	delete(g.nodeMap, n.ID())

	for succ := range g.successors[n.ID()] {
		delete(g.predecessors[succ], n.ID())
	}
	delete(g.successors, n.ID())

	for pred := range g.predecessors[n.ID()] {
		delete(g.successors[pred], n.ID())
	}
	delete(g.predecessors, n.ID())

	g.freeMap[n.ID()] = struct{}{}
}

func (g *DirectedGraph) RemoveEdge(e graph.Edge) {
	from, to := e.From(), e.To()
	if _, ok := g.nodeMap[from.ID()]; !ok {
		return
	} else if _, ok := g.nodeMap[to.ID()]; !ok {
		return
	}

	delete(g.successors[from.ID()], to.ID())
	delete(g.predecessors[to.ID()], from.ID())
}

func (g *DirectedGraph) EmptyGraph() {
	g.successors = make(map[int]map[int]graph.Edge)
	g.predecessors = make(map[int]map[int]graph.Edge)
	g.nodeMap = make(map[int]graph.Node)
}

/* Graph implementation */

func (g *DirectedGraph) From(n graph.Node) []graph.Node {
	if _, ok := g.successors[n.ID()]; !ok {
		return nil
	}

	successors := make([]graph.Node, len(g.successors[n.ID()]))
	i := 0
	for succ := range g.successors[n.ID()] {
		successors[i] = g.nodeMap[succ]
		i++
	}

	return successors
}

func (g *DirectedGraph) HasEdgeBetween(x, y graph.Node) bool {
	xid := x.ID()
	yid := y.ID()
	if _, ok := g.nodeMap[xid]; !ok {
		return false
	}
	if _, ok := g.nodeMap[yid]; !ok {
		return false
	}
	if _, ok := g.successors[xid][yid]; ok {
		return true
	}
	_, ok := g.successors[yid][xid]
	return ok
}

func (g *DirectedGraph) Edge(u, v graph.Node) graph.Edge {
	if _, ok := g.nodeMap[u.ID()]; !ok {
		return nil
	}
	if _, ok := g.nodeMap[v.ID()]; !ok {
		return nil
	}
	edge, ok := g.successors[u.ID()][v.ID()]
	if !ok {
		return nil
	}
	return edge
}

func (g *DirectedGraph) HasEdgeFromTo(u, v graph.Node) bool {
	if _, ok := g.nodeMap[u.ID()]; !ok {
		return false
	}
	if _, ok := g.nodeMap[v.ID()]; !ok {
		return false
	}
	if _, ok := g.successors[u.ID()][v.ID()]; !ok {
		return false
	}
	return true
}

func (g *DirectedGraph) To(n graph.Node) []graph.Node {
	if _, ok := g.successors[n.ID()]; !ok {
		return nil
	}

	predecessors := make([]graph.Node, len(g.predecessors[n.ID()]))
	i := 0
	for succ := range g.predecessors[n.ID()] {
		predecessors[i] = g.nodeMap[succ]
		i++
	}

	return predecessors
}

func (g *DirectedGraph) Node(id int) graph.Node {
	return g.nodeMap[id]
}

func (g *DirectedGraph) Has(n graph.Node) bool {
	_, ok := g.nodeMap[n.ID()]

	return ok
}

func (g *DirectedGraph) Degree(n graph.Node) int {
	if _, ok := g.nodeMap[n.ID()]; !ok {
		return 0
	}

	return len(g.successors[n.ID()]) + len(g.predecessors[n.ID()])
}

func (g *DirectedGraph) InDegree(n graph.Node) int {
	if _, ok := g.nodeMap[n.ID()]; !ok {
		return 0
	}

	return len(g.predecessors[n.ID()])
}

func (g *DirectedGraph) OutDegree(n graph.Node) int {
	if _, ok := g.nodeMap[n.ID()]; !ok {
		return 0
	}

	return len(g.successors[n.ID()])
}

func (g *DirectedGraph) Nodes() []graph.Node {
	nodes := make([]graph.Node, len(g.successors))
	i := 0
	for _, n := range g.nodeMap {
		nodes[i] = n
		i++
	}

	return nodes
}

func (g *DirectedGraph) Weight(x, y graph.Node) (w float64, ok bool) {
	xid := x.ID()
	yid := y.ID()
	if xid == yid {
		return g.self, true
	}
	if s, ok := g.successors[xid]; ok {
		if e, ok := s[yid]; ok {
			return e.Weight(), true
		}
	}
	return g.absent, false
}

func (g *DirectedGraph) Edges() []graph.Edge {
	var edges []graph.Edge
	for _, u := range g.nodeMap {
		for _, e := range g.successors[u.ID()] {
			edges = append(edges, e)
		}
	}
	return edges
}
