package metrics

import (
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	CreatedPVZ = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "created_pvz_total",
			Help: "Total number of created PVZ.",
		},
		[]string{"city"},
	)

	CreatedProducts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "created_products_total",
			Help: "Total number of created products.",
		},
		[]string{"pvz_id"},
	)

	CreatedReceptions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "created_receptions_total",
			Help: "Total number of created receptions.",
		},
		[]string{"pvz_id"},
	)
)

func IncCreatedPVZ(city string) {
	CreatedPVZ.WithLabelValues(city).Inc()
}

func IncCreatedProducts(pvzId uuid.UUID) {
	CreatedProducts.WithLabelValues(pvzId.String()).Inc()
}

func IncCreatedReceptions(pvzId uuid.UUID) {
	CreatedReceptions.WithLabelValues(pvzId.String()).Inc()
}
