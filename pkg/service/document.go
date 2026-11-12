package service

import (
	"github.com/ProductBuildersHQ/outcome-model/pkg/document"
	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/report"
)

// Document assembles the outcome model as one human-readable document with the
// sections opts names, in that order (see document.Build). st and sets may be
// nil, for a document without observations. Invalid outcome models are
// rejected by the graph build, as for every other read-only command.
func Document(doc *model.Model, st *model.Status, sets []*model.ObservationSet, opts document.Options) (*document.Document, error) {
	g, err := BuildGraph(doc)
	if err != nil {
		return nil, err
	}
	whole, err := report.BuildView(doc, g, st, report.Options{})
	if err != nil {
		return nil, err
	}
	pathways, err := report.PathwayViews(doc, g, st)
	if err != nil {
		return nil, err
	}
	return document.Build(document.Input{
		Model:        doc,
		Whole:        whole,
		Pathways:     pathways,
		Status:       st,
		Observations: sets,
	}, opts)
}
