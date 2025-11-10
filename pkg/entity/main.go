package entity

import "encoding/xml"

type DailyValueSQL struct {
	ValuteId  string
	NumCode   string
	CharCode  string
	Nominal   int
	Name      string
	Value     string
	VunitRate string
	Day       string
}

type DailyValueXml struct {
	XMLName xml.Name `xml:"ValCurs" json:"-" swaggerignore:"true"`
	Text    string   `xml:",chardata"`
	Date    string   `xml:"Date,attr,omitempty"`
	Name    string   `xml:"name,attr,omitempty"`
	Valute  []*Value `xml:"Valute,omitempty"`
}
type Value struct {
	Text      string `xml:",chardata"`
	ID        string `xml:"ID,attr"`
	NumCode   string `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Nominal   int    `xml:"Nominal"`
	Name      string `xml:"Name"`
	Value     string `xml:"Value"`
	VunitRate string `xml:"VunitRate"`
}
