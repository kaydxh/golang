package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

// M is the Metrics instance
var M = newMetrics()

const (
	MetircLabelMethod      = "method"
	MetircLabelClientIP    = "client_ip"
	MetircLabelCodeMessage = "code_message"
)

type Metrics struct {
	CalledTotal           *prometheus.CounterVec
	DurationCostHistogram *prometheus.HistogramVec
	DurationCost          *prometheus.SummaryVec
}

func newMetrics() *Metrics {

	return &Metrics{
		CalledTotal: newCounterVec(
			"http_method_called_total",
			"the total number of method called.",
			MetircLabelMethod,
			MetircLabelClientIP,
			MetircLabelCodeMessage,
		),

		DurationCost: newSummaryVec(
			"http_method_cost_duration_milliseconds",
			"the duration of method cost in milliseconds.",
			MetircLabelMethod,
		),
	}
}

//var (
/*
	cost = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "method_cost_duration_milliseconds",
			Help: "the duration of method now cost in milliseconds.",
		},
		[]string{MetircLabelMethod},
	)
*/

/*
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
*/

//)

func newCounterVec(name, help string, labels ...string) *prometheus.CounterVec {
	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		}, labels)
	prometheus.MustRegister(vec)
	return vec
}

func newGaugeVec(name, help string, labels ...string) *prometheus.GaugeVec {
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		}, labels)
	prometheus.MustRegister(vec)
	return vec
}

func newSummaryVec(name, help string, labels ...string) *prometheus.SummaryVec {
	vec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       name,
			Help:       help,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}, labels)
	prometheus.MustRegister(vec)
	return vec
}