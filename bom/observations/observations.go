package observations

import (
	"bom_exporter/bom/schema"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// MetricNames is a list of metrics exported by the observations collector.
var MetricNames = []string{
	"bom_observations_temperature",
	"bom_observations_wind_speed",
	"bom_observations_humidity",
	"bom_observations_pressure",
	"bom_observations_visibility",
	"bom_observations_cloud_base",
	"bom_observations_cloud_cover",
	"bom_observations_wind_direction",
	"bom_observations_rainfall",
}

// Observations combines unmarshalled observations data and the corresponding
// Prometheus metrics.
type Observations struct {
	product         schema.Product
	temperatureDesc *prometheus.Desc
	windSpeedDesc   *prometheus.Desc
	humidityDesc    *prometheus.Desc
	pressureDesc    *prometheus.Desc
	visibilityDesc  *prometheus.Desc
	cloudBaseDesc   *prometheus.Desc
	cloudDesc       *prometheus.Desc
	windDirDesc     *prometheus.Desc
	rainfallDesc    *prometheus.Desc
}

// New creates a new observations collector.
func New(product schema.Product) *Observations {
	var o Observations

	o.product = product

	labels := prometheus.Labels{"identifier": product.Amoc.Identifier}

	o.temperatureDesc = prometheus.NewDesc(
		prometheus.BuildFQName("bom", "observations", "temperature"),
		"Temperature observation.",
		[]string{"bom_id", "wmo_id", "station_name", "latitude", "longitude", "description", "region", "index", "type", "units"},
		labels)

	o.windSpeedDesc = prometheus.NewDesc(
		prometheus.BuildFQName("bom", "observations", "wind_speed"),
		"Wind speed.",
		[]string{"bom_id", "wmo_id", "station_name", "latitude", "longitude", "description", "region", "index", "type", "units"},
		labels)

	o.humidityDesc = prometheus.NewDesc(
		prometheus.BuildFQName("bom", "observations", "humidity"),
		"Relative humidity.",
		[]string{"bom_id", "wmo_id", "station_name", "latitude", "longitude", "description", "region", "index", "units"},
		labels)

	o.pressureDesc = prometheus.NewDesc(
		prometheus.BuildFQName("bom", "observations", "pressure"),
		"Atmospheric pressure.",
		[]string{"bom_id", "wmo_id", "station_name", "latitude", "longitude", "description", "region", "index", "type", "units"},
		labels)

	o.visibilityDesc = prometheus.NewDesc(
		prometheus.BuildFQName("bom", "observations", "visibility"),
		"Visibility.",
		[]string{"bom_id", "wmo_id", "station_name", "latitude", "longitude", "description", "region", "index", "units"},
		labels)

	o.cloudBaseDesc = prometheus.NewDesc(
		prometheus.BuildFQName("bom", "observations", "cloud_base"),
		"Cloud base.",
		[]string{"bom_id", "wmo_id", "station_name", "latitude", "longitude", "description", "region", "index", "units"},
		labels)

	o.cloudDesc = prometheus.NewDesc(
		prometheus.BuildFQName("bom", "observations", "cloud_cover"),
		"Cloud cover.",
		[]string{"bom_id", "wmo_id", "station_name", "latitude", "longitude", "description", "region", "index", "units"},
		labels)

	o.windDirDesc = prometheus.NewDesc(
		prometheus.BuildFQName("bom", "observations", "wind_direction"),
		"Wind direction.",
		[]string{"bom_id", "wmo_id", "station_name", "latitude", "longitude", "description", "region", "index", "units"},
		labels)

	o.rainfallDesc = prometheus.NewDesc(
		prometheus.BuildFQName("bom", "observations", "rainfall"),
		"Rain since 9am.",
		[]string{"bom_id", "wmo_id", "station_name", "latitude", "longitude", "description", "region", "index", "units"},
		labels)

	return &o
}

var temperatureLabelMap map[string]string = map[string]string{"apparent_temp": "apparent",
	"air_temperature":         "ambient",
	"maximum_air_temperature": "maximum",
	"minimum_air_temperature": "minimum",
	"dew_point":               "dew_point",
	"delta_t":                 "delta_t",
}

var windTypeLabelMap map[string]string = map[string]string{"gust_kmh": "gust",
	"wind_gust_spd": "gust",
	"wind_spd_kmh":  "average",
	"wind_spd":      "average",
}

var pressureTypeLabelMap map[string]string = map[string]string{"pres": "absolute",
	"msl_pres": "msl",
	"qnh_pres": "qnh",
}

func (o *Observations) processPeriod(station *schema.Station, ch chan<- prometheus.Metric) {
	region := o.product.Amoc.Source.Region
	for _, e := range station.Period.Level.Element {
		log.Infof("Type: %s, Value: %s, Units: %s", e.Type, e.Value, e.Unit)
		switch e.Type {
		case "apparent_temp", "air_temperature", "maximum_air_temperature", "minimum_air_temperature", "dew_point", "delta_t":
			v, err := strconv.ParseFloat(e.Value, 64)
			if err == nil {
				ch <- prometheus.NewMetricWithTimestamp(time.Time(station.Period.TimeUTC),
					prometheus.MustNewConstMetric(
						o.temperatureDesc,
						prometheus.GaugeValue,
						v,
						station.WmoID,
						station.BomID,
						station.Name,
						fmt.Sprintf("%f", station.Latitude),
						fmt.Sprintf("%f", station.Longitude),
						station.Description,
						region,
						station.Period.Index,
						temperatureLabelMap[e.Type],
						e.Unit))
			}
		case "gust_kmh", "wind_gust_spd", "wind_spd_kmh", "wind_spd":
			v, err := strconv.ParseFloat(e.Value, 64)
			if err == nil {
				ch <- prometheus.NewMetricWithTimestamp(time.Time(station.Period.TimeUTC),
					prometheus.MustNewConstMetric(
						o.windSpeedDesc,
						prometheus.GaugeValue,
						v,
						station.WmoID,
						station.BomID,
						station.Name,
						fmt.Sprintf("%f", station.Latitude),
						fmt.Sprintf("%f", station.Longitude),
						station.Description,
						region,
						station.Period.Index,
						windTypeLabelMap[e.Type],
						e.Unit))
			}
		case "rel-humidity":
			v, err := strconv.ParseFloat(e.Value, 64)
			if err == nil {
				ch <- prometheus.NewMetricWithTimestamp(time.Time(station.Period.TimeUTC),
					prometheus.MustNewConstMetric(
						o.humidityDesc,
						prometheus.GaugeValue,
						v,
						station.WmoID,
						station.BomID,
						station.Name,
						fmt.Sprintf("%f", station.Latitude),
						fmt.Sprintf("%f", station.Longitude),
						station.Description,
						region,
						station.Period.Index,
						e.Unit))
			}
		case "pres", "msl_pres", "qnh_pres":
			v, err := strconv.ParseFloat(e.Value, 64)
			if err == nil {
				ch <- prometheus.NewMetricWithTimestamp(time.Time(station.Period.TimeUTC),
					prometheus.MustNewConstMetric(
						o.pressureDesc,
						prometheus.GaugeValue,
						v,
						station.WmoID,
						station.BomID,
						station.Name,
						fmt.Sprintf("%f", station.Latitude),
						fmt.Sprintf("%f", station.Longitude),
						station.Description,
						region,
						station.Period.Index,
						pressureTypeLabelMap[e.Type],
						e.Unit))
			}
		case "vis_km":
			v, err := strconv.ParseFloat(e.Value, 64)
			if err == nil {
				ch <- prometheus.NewMetricWithTimestamp(time.Time(station.Period.TimeUTC),
					prometheus.MustNewConstMetric(
						o.visibilityDesc,
						prometheus.GaugeValue,
						v,
						station.WmoID,
						station.BomID,
						station.Name,
						fmt.Sprintf("%f", station.Latitude),
						fmt.Sprintf("%f", station.Longitude),
						station.Description,
						region,
						station.Period.Index,
						e.Unit))
			}
		case "cloud_base_m":
			v, err := strconv.ParseFloat(e.Value, 64)
			if err == nil {
				ch <- prometheus.NewMetricWithTimestamp(time.Time(station.Period.TimeUTC),
					prometheus.MustNewConstMetric(
						o.cloudBaseDesc,
						prometheus.GaugeValue,
						v,
						station.WmoID,
						station.BomID,
						station.Name,
						fmt.Sprintf("%f", station.Latitude),
						fmt.Sprintf("%f", station.Longitude),
						station.Description,
						region,
						station.Period.Index,
						e.Unit))
			}
		case "cloud_oktas":
			v, err := strconv.ParseFloat(e.Value, 64)
			if err == nil {
				ch <- prometheus.NewMetricWithTimestamp(time.Time(station.Period.TimeUTC),
					prometheus.MustNewConstMetric(
						o.cloudDesc,
						prometheus.GaugeValue,
						v,
						station.WmoID,
						station.BomID,
						station.Name,
						fmt.Sprintf("%f", station.Latitude),
						fmt.Sprintf("%f", station.Longitude),
						station.Description,
						region,
						station.Period.Index,
						"oktas"))
			}

		case "wind_dir_deg":
			v, err := strconv.ParseFloat(e.Value, 64)
			if err == nil {
				ch <- prometheus.NewMetricWithTimestamp(time.Time(station.Period.TimeUTC),
					prometheus.MustNewConstMetric(
						o.windDirDesc,
						prometheus.GaugeValue,
						v,
						station.WmoID,
						station.BomID,
						station.Name,
						fmt.Sprintf("%f", station.Latitude),
						fmt.Sprintf("%f", station.Longitude),
						station.Description,
						region,
						station.Period.Index,
						e.Unit))
			}

		case "rainfall_24hr":
			v, err := strconv.ParseFloat(e.Value, 64)
			if err == nil {
				ch <- prometheus.NewMetricWithTimestamp(time.Time(station.Period.TimeUTC),
					prometheus.MustNewConstMetric(
						o.rainfallDesc,
						prometheus.GaugeValue,
						v,
						station.WmoID,
						station.BomID,
						station.Name,
						fmt.Sprintf("%f", station.Latitude),
						fmt.Sprintf("%f", station.Longitude),
						station.Description,
						region,
						station.Period.Index,
						e.Unit))
			}

		}
	}

}

// Describe implements the Prometheus Collector interface.
func (o *Observations) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(o, ch)
}

// Collect implements the Prometheus Collector interface.
func (o *Observations) Collect(ch chan<- prometheus.Metric) {
	for _, s := range o.product.Observations.Station {
		o.processPeriod(&s, ch)
	}
}

// Dump logs the contents of an the unmarshalled observations.
func (o *Observations) Dump() {
	p := o.product

	fmt.Printf("%+v\n", p)
}
