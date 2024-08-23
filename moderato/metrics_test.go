package moderato

import (
	"fmt"
	"strings"
	"testing"

	"github.com/moderato-app/pprof/internal/measurement"
	"github.com/stretchr/testify/assert"
)

func TestItems(t *testing.T) {

	tt := []struct {
		filename     string
		mustContains []string
	}{
		{
			"./test.pprof.samples.inuse_space.025.pb.gz",
			[]string{"2788862 total",
				"1212696 43.48% 43.48%    1212696 43.48%  go.etcd.io/etcd/server/v3/storage/wal.newEncoder",
				"524296 18.80%   100%     524296 18.80%  github.com/prometheus/client_golang/prometheus.MakeLabelPairs",
			},
		},
		{
			"./test.pprof.samples.cpu.025.pb.gz",
			[]string{"150000000 total",
				"60000000 40.00% 40.00%   60000000 40.00%  runtime.pthread_cond_signal",
				"10000000  6.67%   100%   10000000  6.67%  syscall.syscall",
			},
		},
		{
			"./test.goroutine.019.pb.gz",
			[]string{"37 total",
				"34 91.89% 91.89%         34 91.89%  runtime.gopark",
				"1  2.70% 94.59%          1  2.70%  runtime.goroutineProfileWithLabels",
			},
		},
		{
			"./test.pprof.alloc_objects.alloc_space.inuse_objects.inuse_space.015.pb.gz",
			[]string{"160278937496 total",
				"67097020760 41.86% 41.86% 67097020760 41.86%  go.etcd.io/raft/v3.(*MemoryStorage).Append",
				"52698228258 32.88% 74.74% 52698228258 32.88%  go.etcd.io/raft/v3.(*MemoryStorage).Compact",
			},
		},
	}
	for _, tc := range tt {
		t.Run("", func(t *testing.T) {
			metrics, err := GetMetricsFromFile(tc.filename)
			if err != nil {
				t.Fatalf("perf: %v", err)
			}
			s, err := getText(metrics)
			if err != nil {
				t.Fatalf("perf: %v", err)
			}
			fmt.Println(s)
			for _, v := range tc.mustContains {
				assert.Contains(t, s, v)
			}
		})
	}
}

func getText(metrics *Metrics) (string, error) {
	b := new(strings.Builder)

	fmt.Fprintln(b, strings.Join(metrics.Labels, "\n"))
	fmt.Fprintf(b, "%10s %5s%% %5s%% %10s %5s%%\n",
		"flat", "flat", "sum", "cum", "cum")
	var flatSum int64
	for _, item := range metrics.Items {
		inl := item.InlineLabel
		if inl != "" {
			inl = " " + inl
		}
		flatSum += item.Flat
		total := metrics.Total
		fmt.Fprintf(b, "%10d %s %s %10d %s  %s%s\n",
			item.Flat, measurement.Percentage(item.Flat, total),
			measurement.Percentage(flatSum, total),
			item.Cum, measurement.Percentage(item.Cum, total),
			item.Func, inl)
	}
	return b.String(), nil
}
