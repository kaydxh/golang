package interceptorprometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	MetircLabelMethod      = "method"
	MetircLabelClientIP    = "client_ip"
	MetircLabelCodeMessage = "code_message"
)

var (
	calledTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "method_called_total",
		Help: "the total number of method called.",
	},
		[]string{MetircLabelMethod, MetircLabelClientIP, MetircLabelCodeMessage},
	)

	cost = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "method_cost_duration_milliseconds",
			Help: "the duration of method now cost in milliseconds.",
		},
		[]string{MetircLabelMethod},
	)
)
