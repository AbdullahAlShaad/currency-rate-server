package controller

import (
	"github.com/shaad7/currency-rate-server/database"
	"github.com/shaad7/currency-rate-server/util"
	"log"
	"strconv"
	"strings"
)

type Controller struct {
	sqlDB *database.Database
}

// NewController takes mysql username, password and dbName as argument and gets a database.Database object which contains connection to SQL db
// It initializes a Controller instance with database.Database object and returns it
func NewController(dbAddress string, password string, mysqlDBName string) (*Controller, error) {
	db, err := database.GetDatabase(dbAddress, mysqlDBName, password)
	if err != nil {
		return nil, err
	}
	return &Controller{
		sqlDB: db,
	}, nil
}

// EnsureDataLoadedToDatabase functions parse data from https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml and load that's data to mysql db
func (c *Controller) EnsureDataLoadedToDatabase() error {
	data := ParseData()
	return c.LoadData(data)
}

// LoadData creates mysql table ExchangeRate and inserts data into the rows in the table
func (c *Controller) LoadData(data *RateSliceStruct) error {
	err := c.sqlDB.CreateTable()
	if err != nil {
		return err
	}

	for _, rate := range data.Rates {
		curRate, err := strconv.ParseFloat(strings.TrimSpace(rate.Rate), 64)
		if err != nil {
			log.Printf("Can't convert %s to float\n", rate.Rate)
			return err
		}
		err = c.sqlDB.InsertRow(rate.Date, rate.Currency, curRate)
		if err != nil {
			return err
		}
	}
	log.Println("Successfully Loaded data to mysql database")
	return nil
}

// GetLatestExchangeRate query to mysql for the latest exchange rate for all currency and return that data in map[string]interface{} in following format
//
//	{
//		"base": "EUR",
//		"rates": {
//			"AUD": 1.5499,
//			"BGN": 1.9558,
//			"BRL": 5.5414,
//	}
func (c *Controller) GetLatestExchangeRate() (map[string]interface{}, error) {
	latestRate, err := c.sqlDB.GetLatestExchangeRate()
	if err != nil {
		return nil, err
	}
	var dataMap = make(map[string]float64)
	for _, data := range latestRate {
		dataMap[data.CurrencyName] = data.Rate
	}
	return map[string]interface{}{
		"base":  util.BaseCurrency,
		"rates": dataMap,
	}, nil
}

// GetExchangeRate query to mysql for the exchange rate of specified date and return that data in map[string]interface{} in following format
//
//	{
//		"base": "EUR",
//		"rates": {
//			"AUD": 1.5499,
//			"BGN": 1.9558,
//			"BRL": 5.5414,
//	}
func (c *Controller) GetExchangeRate(date string) (map[string]interface{}, error) {
	dataList, err := c.sqlDB.GetExchangeRate(date)
	if err != nil {
		return nil, err
	}
	var dataMap = make(map[string]float64)
	for _, data := range dataList {
		dataMap[data.CurrencyName] = data.Rate
	}
	return map[string]interface{}{
		"base":  util.BaseCurrency,
		"rates": dataMap,
	}, nil
}

// GetAnalysis gets min,max and average of each currency and return that data in map[string]interface{} in following format
//
//	{
//		"base": "EUR",
//		"rates_analyze": {
//			"AUD": {
//			"avg": 1.5555093749999997,
//			"max": 1.5972,
//			"min": 1.5289
//		},
//	}
func (c *Controller) GetAnalysis() (map[string]interface{}, error) {
	analysis, err := c.sqlDB.GetExchangeAnalysis()
	if err != nil {
		return nil, err
	}
	var dataMap = make(map[string]interface{})
	for _, data := range analysis {
		var tempMap = map[string]float64{
			"min": data.MinRate,
			"max": data.MaxRate,
			"avg": data.AverageRate,
		}
		dataMap[data.CurrencyName] = tempMap
	}
	return map[string]interface{}{
		"base":          util.BaseCurrency,
		"rates_analyze": dataMap,
	}, nil
}
