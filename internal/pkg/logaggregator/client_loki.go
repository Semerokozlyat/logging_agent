package logaggregator

import (
	"os"

	kitlog "github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"

	lokiClient "github.com/grafana/loki/v3/clients/pkg/promtail/client"

	"github.com/Semerokozlyat/logging_agent/internal/config"
)

const (
	maxStreams = 10
)

func NewLokiClient(cfg lokiClient.Config, collectionCfg config.Collection) (lokiClient.Client, error) {
	metrics := lokiClient.NewMetrics(prometheus.DefaultRegisterer)

	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "component", "loki_client")

	return lokiClient.New(metrics, cfg, maxStreams, collectionCfg.MaxLineLength, false, logger)
}
