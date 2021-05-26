package modconfig

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/zclconf/go-cty/cty"
)

// Report is a struct representing the Report resource
type Report struct {
	FullName  string `cty:"name"`
	ShortName string
	Title     *string

	Reports []*Report //`hcl:"report,block"`
	Panels  []*Panel  //`hcl:"panel,block"`

	DeclRange hcl.Range
	metadata  *ResourceMetadata
}

func NewReport(block *hcl.Block) *Report {
	report := &Report{
		ShortName: block.Labels[0],
		FullName:  fmt.Sprintf("report.%s", block.Labels[0]),
		DeclRange: block.DefRange,
	}
	return report
}

// CtyValue implements HclResource
func (r *Report) CtyValue() (cty.Value, error) {
	return getCtyValue(r)
}

// Name implements HclResource
// return name in format: 'panel.<shortName>'
func (r *Report) Name() string {
	return r.FullName
}

// QualifiedName returns the name in format: '<modName>.report.<shortName>'
func (r *Report) QualifiedName() string {
	return fmt.Sprintf("%s.%s", r.metadata.ModShortName, r.FullName)
}

// OnDecoded implements HclResource
func (r *Report) OnDecoded(*hcl.Block) {}

// AddReference implements HclResource
func (r *Report) AddReference(reference string) {
	// TODO
}

// AddChild implements ReportTreeItem
func (r *Report) AddChild(child ReportTreeItem) {
	switch c := child.(type) {
	case *Panel:
		r.Panels = append(r.Panels, c)
	case *Report:
		r.Reports = append(r.Reports, c)
	}
	return nil
}

// GetPanels implements ReportTreeItem
func (r *Report) GetPanels() []*Panel {
	return r.Panels
}

// GetReports implements ReportTreeItem
func (r *Report) GetReports() []*Report {
	return r.Reports
}

// GetMetadata implements ResourceWithMetadata
func (r *Report) GetMetadata() *ResourceMetadata {
	return r.metadata
}

// SetMetadata implements ResourceWithMetadata
func (r *Report) SetMetadata(metadata *ResourceMetadata) {
	r.metadata = metadata
}

func (r *Report) Diff(new *Report) *ReportTreeItemDiffs {
	res := &ReportTreeItemDiffs{
		Item: r,
		Name: r.Name(),
	}

	if typehelpers.SafeString(r.Title) != typehelpers.SafeString(new.Title) {
		res.AddPropertyDiff("Title")
	}

	res.populateChildDiffs(r, new)
	return res
}
