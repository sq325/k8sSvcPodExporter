package metric

const (
	metricsPath = "/metrics"
	healthzPath = "/healthz"
)

type Labels map[string]string
type Type string

var Gauge Type = "gauge"
var Counter Type = "counter"

// Metric represents a single time series.
type Metric struct {
	Labels Labels
	Value  float64
}

// Family represents a set of metrics with the same name and help text.
type Family struct {
	Name    string
	Type    Type
	Metrics []*Metric
}

// FamilyGenerator provides everything needed to generate a metric family with a
// Kubernetes object.
// DeprecatedVersion is defined only if the metric for which this options applies is,
// in fact, deprecated.
type FamilyGenerator struct {
	Name         string
	Help         string
	Type         Type
	GenerateFunc func(obj any) *Family
}

// NewFamilyGenerator creates new FamilyGenerator instances.
func NewFamilyGenerator(name, help string, _type Type, genFunc func(obj any) *Family) *FamilyGenerator {
	return &FamilyGenerator{
		Name:         name,
		Help:         help,
		Type:         _type,
		GenerateFunc: genFunc,
	}
}
