package bom

import (
	"bom_exporter/bom/connection/file"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"testing"
)

func TestForecastCollector(t *testing.T) {
	f := New(file.New("schema/IDS10034.xml"))
	err := f.RetrieveAndParse()
	if err != nil {
		t.Errorf("Failed to retrieve and parse '%v': %v", f, err)
	}

	count := testutil.CollectAndCount(f,
		"bom_forecast_precis",
		"bom_forecast_temperature",
		"bom_forecast_precipitation_probability",
		"bom_forecast_icon_code")
	if count != 24 {
		t.Errorf("Got %d metrics, expected %d", count, 24)
	}

	problems, err := testutil.CollectAndLint(f,
		"bom_forecast_precis",
		"bom_forecast_temperature",
		"bom_forecast_precipitation_probability",
		"bom_forecast_icon_code",
		"bom_blah")
	if err != nil {
		t.Errorf("CollectAndLint failed: %v", err)
	}

	if len(problems) > 0 {
		t.Errorf("%v", problems)
	}
}

func TestObservationsCollector(t *testing.T) {
	o := New(file.New("schema/IDS60920.xml"))
	err := o.RetrieveAndParse()
	if err != nil {
		t.Errorf("Failed to retrieve and parse '%v': %v", o, err)
	}

	count := testutil.CollectAndCount(o,
		"bom_observations_temperature",
		"bom_observations_wind_speed",
		"bom_observations_humidity",
		"bom_observations_pressure",
		"bom_observations_visibility",
		"bom_observations_cloud_base",
		"bom_observations_cloud_cover",
		"bom_observations_wind_direction",
		"bom_observations_rainfall")
	if count != 1184 {
		t.Errorf("Got %d metrics, expected %d", count, 1184)
	}

	problems, err := testutil.CollectAndLint(o,
		"bom_observations_temperature",
		"bom_observations_wind_speed",
		"bom_observations_humidity",
		"bom_observations_pressure",
		"bom_observations_visibility",
		"bom_observations_cloud_base",
		"bom_observations_cloud_cover",
		"bom_observations_wind_direction",
		"bom_observations_rainfall")
	if err != nil {
		t.Errorf("CollectAndLint failed: %v", err)
	}

	if len(problems) > 0 {
		t.Errorf("%v", problems)
	}
}
