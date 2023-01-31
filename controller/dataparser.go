package controller

import (
	"encoding/xml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/shaad7/currency-rate-server/util"
	"io"
	"log"
	"net/http"
)

// Rate stores column of given data
type Rate struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
	Date     string `xml:"time,attr"`
}

type RateSliceStruct struct {
	Rates RateSlice `xml:"Cube>Cube"`
}

type RateSlice []Rate

func (rs *RateSlice) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	date := start.Attr[0].Value

	for {
		tok, err := d.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if se, ok := tok.(xml.StartElement); ok {
			rate := Rate{Date: date}
			if err := d.DecodeElement(&rate, &se); err != nil {
				return err
			}

			*rs = append(*rs, rate)
		}
	}
}

// ParseData parses data from util.DataUTL and return that data
func ParseData() *RateSliceStruct {
	log.Printf("Reading data from %s", util.DataURL)
	resp, err := http.Get(util.DataURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	}

	v := &RateSliceStruct{}
	if err := xml.Unmarshal(bodyBytes, v); err != nil {
		panic(err)
	}
	return v
}
