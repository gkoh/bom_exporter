package schema

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
	"time"
)

func TestTimeField(t *testing.T) {
	l0 := time.FixedZone("ADST", 10.5*60*60)
	v := []struct {
		input    []byte
		expected time.Time
	}{
		{[]byte(`<test><time>2022-03-29T06:02:13Z</time></test>`), time.Date(2022, time.March, 29, 6, 2, 13, 0, time.UTC)},
		{[]byte(`<test><time>2022-03-31T13:30:00Z</time></test>`), time.Date(2022, time.March, 31, 13, 30, 0, 0, time.UTC)},
		{[]byte(`<test><time>2022-04-21T07:30:00+10:30</time></test>`), time.Date(2022, time.April, 21, 7, 30, 0, 0, l0)},
		{[]byte(`<test><time>2022-07-21T06:35:00+00:00</time></test>`), time.Date(2022, time.July, 21, 6, 35, 0, 0, time.UTC)},
	}

	for _, k := range v {
		tt := struct {
			TTime TimeField `xml:"time"`
		}{}
		err := xml.Unmarshal(k.input, &tt)
		if err != nil {
			t.Errorf("Failed to unmarshal '%v': %s", k.input, err)
		}
		if !time.Time(tt.TTime).Equal(k.expected) {
			t.Errorf("Parse mismatch; %+v != %+v", time.Time(tt.TTime), k.expected)
		}
	}
}

func TestTimeFieldAttr(t *testing.T) {
	l0 := time.FixedZone("ADST", 10.5*60*60)
	v := []struct {
		input    []byte
		expected time.Time
	}{
		{[]byte(`<test time="2022-03-29T06:02:13Z"></test>`), time.Date(2022, time.March, 29, 6, 2, 13, 0, time.UTC)},
		{[]byte(`<test time="2022-03-31T13:30:00Z"></test>`), time.Date(2022, time.March, 31, 13, 30, 0, 0, time.UTC)},
		{[]byte(`<test time="2022-04-21T07:30:00+10:30"></test>`), time.Date(2022, time.April, 21, 7, 30, 0, 0, l0)},
		{[]byte(`<test time="2022-07-21T06:35:00+00:00"></test>`), time.Date(2022, time.July, 21, 6, 35, 0, 0, time.UTC)},
	}

	for _, k := range v {
		tt := struct {
			TTime TimeFieldAttr `xml:"time,attr"`
		}{}
		err := xml.Unmarshal(k.input, &tt)
		if err != nil {
			t.Errorf("Failed to unmarshal '%v': %s", k.input, err)
		}
		if !time.Time(tt.TTime).Equal(k.expected) {
			t.Errorf("Parse mismatch; %+v != %+v", time.Time(tt.TTime), k.expected)
		}
	}
}

func TestForecastStruct(t *testing.T) {
	inputs := []struct {
		file  string
		areas int
	}{
		{file: "IDS10034.xml",
			areas: 7},
		{file: "IDS10044.xml",
			areas: 75}}

	for _, x := range inputs {
		data, err := ioutil.ReadFile(x.file)
		if err != nil {
			t.Errorf("Failed to open '%s': %s", x.file, err)
		}

		var p Product
		err = p.Parse(data)
		if err != nil {
			t.Errorf("Failed to unmarshal '%s': %s", x.file, err)
		}

		if p.Forecast == nil {
			t.Errorf("Failed to properly unmarshal Forecast in '%s'", x.file)
		}

		if len(p.Forecast.Area) != x.areas {
			t.Errorf("Failed to parse all Areas, got %d, expected %d", len(p.Forecast.Area), x.areas)
		}

		for _, area := range p.Forecast.Area {
			for _, period := range area.Period {
				DumpPeriod(period)
			}
		}

	}

}

func TestObservationsStruct(t *testing.T) {
	inputs := []struct {
		file     string
		stations int
	}{
		{file: "IDS60920.xml",
			stations: 81},
		{file: "IDT60920.xml",
			stations: 65}}

	for _, x := range inputs {
		data, err := ioutil.ReadFile(x.file)
		if err != nil {
			t.Errorf("Failed to open '%s': %s", x.file, err)
		}

		var p Product
		err = p.Parse(data)
		if err != nil {
			t.Errorf("Failed to unmarshal '%s': %s", x.file, err)
		}

		if p.Observations == nil {
			t.Errorf("Failed to properly unmarshal Observations in '%s'", x.file)
		}

		if len(p.Observations.Station) != x.stations {
			t.Errorf("Failed to parse all Stations, got %d, expected %d", len(p.Observations.Station), x.stations)
		}

		for _, station := range p.Observations.Station {
			DumpElements(station.Period.Level.Element)
		}
	}

}
