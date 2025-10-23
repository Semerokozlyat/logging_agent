package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	nodeNameLabel           = "node_name"
	logFileNamePatternLabel = "log_file_name_pattern"
)

var (
	LogLines = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "log_lines",
		Help: "Number of lines of log processed distinguished by log filename pattern",
	},
		[]string{logFileNamePatternLabel, nodeNameLabel},
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

func MakeLabelsForLogLine(fileNamePattern, nodeName string) prometheus.Labels {
	return map[string]string{
		logFileNamePatternLabel: fileNamePattern,
		nodeNameLabel:           nodeName,
	}
}
