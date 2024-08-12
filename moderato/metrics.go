package moderato

import (
	"io"
	"os"
	"reflect"

	"github.com/moderato-app/pprof/internal/report"
	"github.com/moderato-app/pprof/profile"
)

type Metrics struct {
	items  []report.ModeratoItem
	labels []string
	total  int64
}

// GetMetricsFromFile returns the metrics data of either inuse_space or cpu
func GetMetricsFromFile(path string) (*Metrics, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return GetMetrics(file)
}

// GetMetrics returns the metrics data of either inuse_space or cpu
//
// reader contains bytes in the format of .pb or pg.gz
func GetMetrics(reader io.Reader) (*Metrics, error) {

	prof, err := profile.Parse(reader)

	if err != nil {
		return nil, err
	}

	for _, loc := range prof.Location {
		loc.Address = 0
	}

	index, err := prof.SampleIndexByName(prof.DefaultSampleType)
	if err != nil {
		return nil, err
	}

	o := report.Options{
		OutputFormat: report.Text,
		NodeFraction: 0.05,
		SampleType:   prof.DefaultSampleType,
		SampleValue: func(s []int64) int64 {
			return s[index]
		},
	}

	rpt := report.New(prof, &o)

	items, labels := report.ModeratoItems(rpt)

	v := reflect.ValueOf(*rpt)
	total := v.FieldByName("total").Int()

	return &Metrics{
		items:  items,
		labels: labels,
		total:  total,
	}, nil

}
