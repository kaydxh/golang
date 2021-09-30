package interceptorprometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	MetircLabelMethod   = "method"
	MetircLabelClientIP = "client_ip"
)

var (
	cost = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "method_cost_duration_milliseconds",
			Help: "the duration of method now cost in milliseconds.",
		},
		[]string{MetircLabelMethod},
	)
)
