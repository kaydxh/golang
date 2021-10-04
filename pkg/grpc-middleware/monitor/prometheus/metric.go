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

	/*
		cost = promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "method_cost_duration_milliseconds",
				Help: "the duration of method now cost in milliseconds.",
			},
			[]string{MetircLabelMethod},
		)
	*/

	durationCostHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "method_cost_duration_histogram_milliseconds",
			Help: "Histogram of request duration",
			Buckets: []float64{
				.001,
				.002,
				.004,
				.006,
				.008,
				.01,
				.02,
				.04,
				.06,
				.08,
				.1,
				.2,
				.4,
				.6,
				.8,
				1,
			},
		},
		[]string{MetircLabelMethod},
	)

	durationCost = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "method_cost_duration_milliseconds",
			Help:       "the duration of method cost in milliseconds.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{MetircLabelMethod},
	)
)
