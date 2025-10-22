package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	LabelLogFileNamePattern = "log_file_name_pattern"
)

var (
	LogLines = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "log_lines",
		Help: "Number of lines of log processed distinguished by log filename pattern",
	},
		[]string{LabelLogFileNamePattern},
	)
)

var once sync.Once

func InitMetricsCollector() {
	once.Do(func() {
		prometheus.MustRegister(
			LogLines,
		)
	})
}
