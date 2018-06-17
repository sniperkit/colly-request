package http_stats

import (
	"net/http"

	"github.com/segmentio/stats"
	"github.com/segmentio/stats/httpstats"
)

var (
	statsEngine        *stats.Engine
	httpStatsHandler   http.Handler      = httpstats.NewHandler()
	httpStatsTransport http.RoundTripper = httpstats.NewTransport()
)
