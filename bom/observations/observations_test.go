package observations

import (
	"bom_exporter/bom/schema"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"testing"
	"time"
)

func TestCollector(t *testing.T) {
	identifier := "a5a5a5a5"
	region := "springfield"
	latitude := float32(123.4)
	longitude := float32(-123.4)
	bomid := "111"
	wmoid := "222"
	name := "EVERGREEN"
	description := "Evergreen Terrace"
	issuetime := time.Now()

	elements := []schema.Element{
		{
			Type:  "apparent_temp",
			Unit:  "Celsius",
			Value: "12.34"},
		{
			Type:  "air_temperature",
			Unit:  "Celsius",
			Value: "23.45",
		},
		{
			Type:  "minimum_air_temperature",
			Unit:  "Celsius",
			Value: "34.56",
		},
		{
			Type:  "maximum_air_temperature",
			Unit:  "Celsius",
			Value: "45.67",
		},
		{
			Type:  "dew_point",
			Unit:  "Celsius",
			Value: "1.1",
		},
		{
			Type:  "delta_t",
			Unit:  "Celsius",
			Value: "2.2",
		},
		{
			Type:  "gust_kmh",
			Unit:  "kmh",
			Value: "2",
		},
		{
			Type:  "wind_gust_spd",
			Unit:  "knots",
			Value: "3",
		},
		{
			Type:  "wind_spd_kmh",
			Unit:  "kmh",
			Value: "4",
		},
		{
			Type:  "wind_spd",
			Unit:  "knots",
			Value: "5",
		},
		{
			Type:  "rel-humidity",
			Unit:  "%",
			Value: "55",
		},
		{
			Type:  "pres",
			Unit:  "hPa",
			Value: "1024",
		},
		{
			Type:  "msl_pres",
			Unit:  "hPa",
			Value: "1024",
		},
		{
			Type:  "qnh_pres",
			Unit:  "hPa",
			Value: "1024",
		},
		{
			Type:  "vis_km",
			Unit:  "km",
			Value: "33",
		},
		{
			Type:  "cloud_base_m",
			Unit:  "m",
			Value: "960",
		},
		{
			Type:  "cloud_oktas",
			Unit:  "oktas",
			Value: "6",
		},
		{
			Type:  "wind_dir_deg",
			Unit:  "deg",
			Value: "270",
		},
		{
			Type:  "rainfall_24hr",
			Unit:  "mm",
			Value: "4.6",
		},
	}

	period := schema.Period{
		Index: "0",
		Level: schema.Level{Element: elements},
	}

	stations := []schema.Station{
		{WmoID: wmoid,
			BomID:       bomid,
			Name:        name,
			Description: description,
			Latitude:    latitude,
			Longitude:   longitude,
			Period:      period},
		{WmoID: wmoid + "b",
			BomID:       bomid + "b",
			Name:        name + "b",
			Description: description + "b",
			Latitude:    latitude,
			Longitude:   longitude,
			Period:      period},
	}
	observations := schema.Observations{Station: stations}
	product := schema.Product{
		Amoc: schema.Amoc{
			Identifier:   identifier,
			Source:       schema.Source{Region: region},
			IssueTimeUTC: schema.TimeField(issuetime)},
		Observations: &observations}

	o := New(product)

	o.Dump()

	// 2 stations, elements repeated
	expected := 2 * len(elements)

	count := testutil.CollectAndCount(o, MetricNames...)
	if count != expected {
		t.Errorf("Got %d metrics, expected %d", count, expected)
	}

	problems, err := testutil.CollectAndLint(o, MetricNames...)
	if err != nil {
		t.Errorf("CollectAndLint failed: %v", err)
	}

	if len(problems) > 0 {
		t.Errorf("Problems found: %v", problems)
	}

}
