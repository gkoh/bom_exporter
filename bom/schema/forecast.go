package schema

import (
	"encoding/xml"
)

// Forecast contains the unmarshalled forecast XML data.
type Forecast struct {
	XMLName xml.Name `xml:"forecast"`
	Area    []Area   `xml:"area"`
}

// Area contains the unmarshalled Area XML data.
type Area struct {
	XMLName     xml.Name         `xml:"area"`
	Aac         string           `xml:"aac,attr"`
	Description string           `xml:"description,attr"`
	Type        string           `xml:"type,attr"`
	ParentAac   string           `xml:"parent-aac,attr"`
	Period      []ForecastPeriod `xml:"forecast-period"`
}

// ForecastPeriod contains the unmarshalled forecast period XML data.
type ForecastPeriod struct {
	XMLName        xml.Name      `xml:"forecast-period"`
	Index          string        `xml:"index,attr"`
	StartTimeLocal TimeFieldAttr `xml:"start-time-local,attr"`
	EndTimeLocal   TimeFieldAttr `xml:"end-time-local,attr"`
	StartTimeUTC   TimeFieldAttr `xml:"start-time-utc,attr"`
	EndTimeUTC     TimeFieldAttr `xml:"end-time-utc,attr"`
	Elements       []Element     `xml:"element"`
	Texts          []Text        `xml:"text"`
}
