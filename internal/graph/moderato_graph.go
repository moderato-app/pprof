package graph

import (
	"fmt"
	"path/filepath"
)

func (i *NodeInfo) Func() string {
	return i.Name
}

func (i *NodeInfo) Line() string {

	switch {
	case i.Lineno != 0:
		s := fmt.Sprintf("%s:%d", i.File, i.Lineno)
		if i.Columnno != 0 {
			s += fmt.Sprintf(":%d", i.Columnno)
		}
		// User requested line numbers, provide what we have.
		return s
	case i.File != "":
		// User requested file name, provide it.
		return i.File
	case i.Name != "":
		// User requested function name. It was already included.
	case i.Objfile != "":
		// Only binary name is available
		return "[" + filepath.Base(i.Objfile) + "]"
	default:
		// Do not leave it empty if there is no information at all.
		return "<unknown>"
	}

	panic("unreachable")
}
