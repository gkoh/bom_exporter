package forecast

import (
	"github.com/gkoh/bom_exporter/bom/schema"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"testing"
	"time"
)

func TestCollector(t *testing.T) {
	identifier := "a5a5a5a5"
	region := "springfield"
	aac := "bart"
	parentaac := "homer"
	description := "Evergreen Terrace"
	issuetime := time.Now()

	elements := []schema.Element{
		{Type: "air_temperature_minimum",
			Unit:  "Celsius",
			Value: "12.34"},
		{Type: "air_temperature_maximum",
			Unit:  "Celsius",
			Value: "43.21"},
		{Type: "forecast_icon_code",
			Unit:  "N/A",
			Value: "3"},
	}

	texts := []schema.Text{
		{Type: "precis",
			Value: "Very fake"},
		{Type: "probability_of_precipitation",
			Value: "55%"},
	}

	periods := []schema.ForecastPeriod{
		{Index: "0", Elements: elements, Texts: texts},
		{Index: "1", Elements: elements, Texts: texts},
	}
	areas := []schema.Area{
		{Description: description, Aac: aac, ParentAac: parentaac, Type: "location", Period: periods},
	}
	forecast := schema.Forecast{Area: areas}
	product := schema.Product{
		Amoc: schema.Amoc{
			Identifier:   identifier,
			Source:       schema.Source{Region: region},
			IssueTimeUTC: schema.TimeField(issuetime)},
		Forecast: &forecast}

	f := New(product)

	f.Dump()

	// 2 index entries, repeated
	expected := 2 * (len(elements) + len(texts))

	count := testutil.CollectAndCount(f, MetricNames...)
	if count != expected {
		t.Errorf("Got %d metrics, expected %d", count, expected)
	}

	problems, err := testutil.CollectAndLint(f, MetricNames...)
	if err != nil {
		t.Errorf("CollectAndLint failed: %v", err)
	}

	if len(problems) > 0 {
		t.Errorf("Problems found: %v", problems)
	}

}
