package models

import "encoding/xml"

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

type ValueCalc struct {
	Name    string
	Min     float64
	MinDate string
	Max     float64
	MaxDate string
	Avg     float64
	Counter int64
}
