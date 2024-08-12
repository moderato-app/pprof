package report

import (
	"github.com/moderato-app/pprof/internal/graph"
)

// ModeratoItem holds a single text report entry.
type ModeratoItem struct {
	Func        string
	Line        string
	InlineLabel string // Not empty if inlined
	Flat        int64  // cpu in ns, inuse_space in byte
	Cum         int64  // Raw values
}

func ModeratoItems(rpt *Report) ([]ModeratoItem, []string) {
	g, origCount, droppedNodes, _ := rpt.newTrimmedGraph()
	rpt.selectOutputUnit(g)
	labels := reportLabels(rpt, graphTotal(g), len(g.Nodes), origCount, droppedNodes, 0, false)

	var items []ModeratoItem
	var flatSum int64
	for _, n := range g.Nodes {
		flat := n.FlatValue()

		flatSum += flat
		items = append(items, ModeratoItem{
			Func:        n.Info.Func(),
			Line:        n.Info.Line(),
			InlineLabel: inlineLabel(n),
			Flat:        flat,
			Cum:         n.CumValue(),
		})
	}
	return items, labels
}

func inlineLabel(n *graph.Node) string {
	var inline, noinline bool
	for _, e := range n.In {
		if e.Inline {
			inline = true
		} else {
			noinline = true
		}
	}

	var inl string
	if inline {
		if noinline {
			inl = "(partial-inline)"
		} else {
			inl = "(inline)"
		}
	}
	return inl
}
