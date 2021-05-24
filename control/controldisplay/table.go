package controldisplay

import (
	"fmt"

	"github.com/turbot/steampipe/control/controlexecute"
)

type TableRenderer struct {
	resultTree *controlexecute.ExecutionTree

	// screen width
	width             int
	maxFailedControls int
	maxTotalControls  int
}

func NewTableRenderer(resultTree *controlexecute.ExecutionTree, width int) *TableRenderer {
	return &TableRenderer{
		resultTree:        resultTree,
		width:             width,
		maxFailedControls: resultTree.Root.Summary.Status.FailedCount(),
		maxTotalControls:  resultTree.Root.Summary.Status.TotalCount(),
	}
}

func (r TableRenderer) Render() string {
	// leading newline
	fmt.Println()
	return NewGroupRenderer(r.resultTree.Root, r.maxFailedControls, r.maxTotalControls, r.resultTree, r.width).Render()
}
