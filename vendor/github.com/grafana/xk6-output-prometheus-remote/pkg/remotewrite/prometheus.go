package remotewrite

import (
	"sort"

	"github.com/mstoykov/atlas"
	prompb "go.buf.build/grpc/go/prometheus/prometheus"
	"go.k6.io/k6/metrics"
)

const namelbl = "__name__"

// MapTagSet converts a k6 tag set into
// the equivalent set of Labels as expected from the
// Prometheus' data model.
func MapTagSet(t *metrics.TagSet) []*prompb.Label {
	n := (*atlas.Node)(t)
	if n.Len() < 1 {
		return nil
	}
	labels := make([]*prompb.Label, 0, n.Len())
	for !n.IsRoot() {
		prev, key, value := n.Data()
		labels = append(labels, &prompb.Label{Name: key, Value: value})
		n = prev
	}
	return labels
}

// MapSeries converts a k6 time series into
// the equivalent set of Labels (name+tags) as expected from the
// Prometheus' data model.
//
// The labels are lexicographic sorted as required
// from the Remote write's specification.
func MapSeries(series metrics.TimeSeries, suffix string) []*prompb.Label {
	v := defaultMetricPrefix + series.Metric.Name
	if suffix != "" {
		v += "_" + suffix
	}
	lbls := append(MapTagSet(series.Tags), &prompb.Label{
		Name:  namelbl,
		Value: v,
	})
	sort.Slice(lbls, func(i int, j int) bool {
		return lbls[i].Name < lbls[j].Name
	})
	return lbls
}
