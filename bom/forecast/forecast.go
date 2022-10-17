package forecast

import (
	"github.com/gkoh/bom_exporter/bom/schema"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

// MetricNames is the list of metrics exported by the forecast collector.
var MetricNames = []string{
	"bom_forecast_precis",
	"bom_forecast_air_temperature",
	"bom_forecast_precipitation_probability",
	"bom_forecast_icon_code",
}

// Forecast combines the unmarshalled forecast data and the corresponding
// Prometheus output metrics.
type Forecast struct {
	product            schema.Product
	precisDesc         *prometheus.Desc
	precipitationDesc  *prometheus.Desc
	airTemperatureDesc *prometheus.Desc
	iconCodeDesc       *prometheus.Desc
}

// New creates an exporter instance based on an unmarshalled Product.
func New(product schema.Product) *Forecast {
	var f Forecast

	f.product = product

	labels := prometheus.Labels{"identifier": product.Amoc.Identifier}

	f.precisDesc = prometheus.NewDesc(
		prometheus.BuildFQName("bom", "forecast", "precis"),
		"Precis forecast text.",
		[]string{"aac", "parent_aac", "description", "region", "index", "precis"}, labels)

	f.precipitationDesc = prometheus.NewDesc(
		prometheus.BuildFQName("bom", "forecast", "precipitation_probability"),
		"Probability of precipitation forecast in percentage.",
		[]string{"aac", "parent_aac", "description", "region", "index"}, labels)

	f.airTemperatureDesc = prometheus.NewDesc(
		prometheus.BuildFQName("bom", "forecast", "air_temperature"),
		"Temperature forecast in Celsius.",
		[]string{"aac", "parent_aac", "description", "region", "index", "units", "type"}, labels)

	f.iconCodeDesc = prometheus.NewDesc(
		prometheus.BuildFQName("bom", "forecast", "icon_code"),
		"Forecast icon code",
		[]string{"aac", "parent_aac", "description", "region", "index"}, labels)

	return &f
}

var temperatureLabelMap map[string]string = map[string]string{"air_temperature_maximum": "maximum",
	"air_temperature_minimum": "minimum",
}

func (f *Forecast) processPeriod(area *schema.Area, period *schema.ForecastPeriod, ch chan<- prometheus.Metric) {
	region := f.product.Amoc.Source.Region
	// Process Text entries
	for _, t := range period.Texts {
		switch t.Type {
		case "precis":
			ch <- prometheus.NewMetricWithTimestamp(time.Time(f.product.Amoc.IssueTimeUTC),
				prometheus.MustNewConstMetric(
					f.precisDesc,
					prometheus.GaugeValue,
					1.0,
					area.Aac,
					area.ParentAac,
					area.Description,
					region,
					period.Index,
					t.Value))

		case "probability_of_precipitation":
			v, err := strconv.Atoi(strings.TrimSuffix(t.Value, "%"))
			if err == nil {
				ch <- prometheus.NewMetricWithTimestamp(time.Time(f.product.Amoc.IssueTimeUTC),
					prometheus.MustNewConstMetric(
						f.precipitationDesc,
						prometheus.GaugeValue,
						float64(v),
						area.Aac,
						area.ParentAac,
						area.Description,
						region,
						period.Index))
			}
		}

		log.Debugf("%s: %v\n", t.Type, t.Value)
	}

	// Process Element entries
	for _, e := range period.Elements {
		v, err := strconv.ParseFloat(e.Value, 64)
		if err == nil {
			switch e.Type {
			case "air_temperature_minimum", "air_temperature_maximum":
				ch <- prometheus.NewMetricWithTimestamp(time.Time(f.product.Amoc.IssueTimeUTC),
					prometheus.MustNewConstMetric(
						f.airTemperatureDesc,
						prometheus.GaugeValue,
						v,
						area.Aac,
						area.ParentAac,
						area.Description,
						region,
						period.Index,
						e.Unit,
						temperatureLabelMap[e.Type]))

			case "forecast_icon_code":
				ch <- prometheus.NewMetricWithTimestamp(time.Time(f.product.Amoc.IssueTimeUTC),
					prometheus.MustNewConstMetric(
						f.iconCodeDesc,
						prometheus.GaugeValue,
						v,
						area.Aac,
						area.ParentAac,
						area.Description,
						region,
						period.Index))
			}
		}

		log.Debugf("%s (%s): %v", e.Type, e.Unit, e.Value)
	}
}

// Describe implements the Prometheus Collector interface.
func (f *Forecast) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(f, ch)
}

// Collect implements the Prometheus Collector interface.
func (f *Forecast) Collect(ch chan<- prometheus.Metric) {
	p := f.product

	for _, area := range p.Forecast.Area {
		if area.Type == "location" {
			log.Debugf("=== %s, %s ===\n", area.Description, p.Amoc.Source.Region)
			for i, period := range area.Period {
				log.Debugf("[%d] %v\n", i, period)
				f.processPeriod(&area, &period, ch)
			}
		}
	}
}

// Dump logs the contents of an unmarshalled Forecast.
func (f *Forecast) Dump() {
	log.Debugf("Dumping")
	p := f.product

	for _, d := range p.Forecast.Area {
		if d.Type == "location" {
			log.Debugf("=== %s, %s ===\n", d.Description, p.Amoc.Source.Region)
			for i, period := range d.Period {
				log.Debugf("[%d]\n", i)
				schema.DumpPeriod(period)
			}
		}
	}
}
