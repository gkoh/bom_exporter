package bom

import (
	"bom_exporter/bom/connection"
	"bom_exporter/bom/forecast"
	"bom_exporter/bom/observations"
	"bom_exporter/bom/schema"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"sync"
)

// A Metric is an instance of a BoM product identifier.
type Metric struct {
	sync.Mutex
	identifier string
	conn       connection.Retriever
	product    schema.Product
}

// New creates a new Metric with the given retriever.
func New(retriever connection.Retriever) *Metric {
	return &Metric{identifier: retriever.Identifier(), conn: retriever}
}

// RetrieveAndParse gathers the data and parses it into the local
// representation.
func (m *Metric) RetrieveAndParse() error {
	data, err := m.conn.Retrieve()
	if err != nil {
		log.Warnf("Failed to retrieve: %s", err)
		return err
	}

	return m.product.Parse(data)
}

// Describe implements the Collector interface.
func (m *Metric) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(m, ch)
}

// Collect implements the Collector interface.
func (m *Metric) Collect(ch chan<- prometheus.Metric) {
	m.Lock()
	defer m.Unlock()

	if m.product.Forecast != nil {
		f := forecast.New(m.product)
		f.Collect(ch)
	} else if m.product.Observations != nil {
		o := observations.New(m.product)
		o.Collect(ch)
	}
}
