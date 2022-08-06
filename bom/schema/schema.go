package schema

import (
	"encoding/xml"
	"time"
)

// Product contains the unmarshalled product XML data.
type Product struct {
	XMLName      xml.Name      `xml:"product"`
	Version      string        `xml:"version"`
	Amoc         Amoc          `xml:"amoc"`
	Forecast     *Forecast     `xml:"forecast"`
	Observations *Observations `xml:"observations"`
}

// Amoc contains the unmarshalled AMOC XML data.
type Amoc struct {
	XMLName                 xml.Name  `xml:"amoc"`
	Source                  Source    `xml:"source"`
	Identifier              string    `xml:"identifier"`
	IssueTimeUTC            TimeField `xml:"issue-time-utc"`
	NextRoutineIssueTimeUTC TimeField `xml:"next-routine-issue-time-utc"`
	ProductType             string    `xml:"product-type"`
}

// Source contains the unmarshalled source XML data.
type Source struct {
	XMLName    xml.Name `xml:"source"`
	Sender     string   `xml:"sender"`
	Region     string   `xml:"region"`
	Office     string   `xml:"office"`
	Copyright  string   `xml:"copyright"`
	Disclaimer string   `xml:"disclaimer"`
}

// Element contains an unmarshalled element XML data instance.
type Element struct {
	XMLName xml.Name `xml:"element"`
	Type    string   `xml:"type,attr"`
	Unit    string   `xml:"units,attr"`
	Value   string   `xml:",chardata"`
}

// Text contains an unmarshalled text XML data instance.
type Text struct {
	XMLName xml.Name `xml:"text"`
	Type    string   `xml:"type,attr"`
	Value   string   `xml:",chardata"`
}

// TimeField is a wrapper type for time.Time.
type TimeField time.Time

// UnmarshalXML unmarshals an RFC3339 formatted XML string into native Time.
func (t *TimeField) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	tm, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return err
	}

	*t = TimeField(tm)
	return nil
}

// TimeFieldAttr is a wrapper type for time.Time.
type TimeFieldAttr time.Time

// UnmarshalXMLAttr unmarshals an RFC3339 formatted XML attribute into native Time.
func (t *TimeFieldAttr) UnmarshalXMLAttr(attr xml.Attr) error {
	tm, err := time.Parse(time.RFC3339, attr.Value)
	if err != nil {
		return err
	}

	*t = TimeFieldAttr(tm)
	return nil
}

// Parse unmarshals XML data into a top level Product.
func (p *Product) Parse(data []byte) error {
	return xml.Unmarshal(data, p)
}
