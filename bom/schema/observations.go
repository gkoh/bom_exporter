package schema

import (
	"encoding/xml"
)

// Observations contains the unmarshalled observations XML data.
type Observations struct {
	XMLName xml.Name  `xml:"observations"`
	Station []Station `xml:"station"`
}

// Station contains the unmarshalled station XML data.
type Station struct {
	XMLName     xml.Name `xml:"station"`
	WmoID       string   `xml:"wmo-id,attr"`
	BomID       string   `xml:"bom-id,attr"`
	Timezone    string   `xml:"tz,attr"`
	Name        string   `xml:"stn-name,attr"`
	Height      float32  `xml:"stn-height,attr"`
	Type        string   `xml:"type,attr"`
	Latitude    float32  `xml:"lat,attr"`
	Longitude   float32  `xml:"lon,attr"`
	Description string   `xml:"description,attr"`
	Period      Period   `xml:"period"`
}

// Period contains the unmarshalled period XML data.
type Period struct {
	XMLName xml.Name      `xml:"period"`
	Index   string        `xml:"index,attr"`
	TimeUTC TimeFieldAttr `xml:"time-utc,attr"`
	Level   Level         `xml:"level"`
}

// Level contains the unmarshalled level XML data.
type Level struct {
	XMLName xml.Name  `xml:"level"`
	Index   string    `xml:"index,attr"`
	Type    string    `xml:"type,attr"`
	Element []Element `xml:"element"`
}
