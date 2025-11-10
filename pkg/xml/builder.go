package xml

import (
	"encoding/xml"
	"github.com/DarkOugi/OZON/pkg/entity"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

var header = "<?xml version=\"1.0\" encoding=\"windows-1251\"?>"

func addHeader(data []byte) []byte {
	data = append([]byte(header), data...)
	return data
}
func decodeTo1251(data []byte) ([]byte, error) {
	data1251, _, err1251 := transform.Bytes(charmap.Windows1251.NewEncoder(), data)
	if err1251 != nil {
		return nil, err1251
	}
	return data1251, nil
}

func CreateXML(name, date string, data []*entity.DailyValueSQL) ([]byte, error) {
	xmlData := entity.DailyValueXml{
		Date: date,
		Name: name,
	}
	var dateDv string

	for _, dv := range data {
		dateDv = dv.Day
		if dv.ValuteId != "" {
			xmlData.Valute = append(xmlData.Valute, &entity.Value{
				ID:        dv.ValuteId,
				NumCode:   dv.NumCode,
				CharCode:  dv.CharCode,
				Nominal:   dv.Nominal,
				Name:      dv.Name,
				Value:     dv.Value,
				VunitRate: dv.VunitRate,
			})
		}
	}

	if dateDv != "" {
		xmlData.Date = dateDv
	}

	dvXML, errXML := xml.Marshal(xmlData)
	if errXML != nil {
		return nil, errXML
	}
	dvXML = addHeader(dvXML)

	return decodeTo1251(dvXML)
}

func CreateErrorXml() ([]byte, error) {
	xmlData := entity.DailyValueXml{Text: "Error in parameters"}

	xmlM, errM := xml.Marshal(xmlData)
	if errM != nil {
		return nil, errM
	}
	xmlH := addHeader(xmlM)

	return decodeTo1251(xmlH)
}
