package util

import "regexp"

const (
	DataURL              = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	DriverName           = "mysql"
	DefaultMySqlUsername = "root"
	DatabaseName         = "currency_exchange"
	TestDatabaseName     = "test_currency_exchange"
	DBPort               = "3306"

	BaseCurrency            = "EUR"
	DBTableName             = "ExchangeRate"
	DBDateColumn            = "Date"
	DBBaseCurrencyColumn    = "BaseCurrency"
	DBCurrentCurrencyColumn = "CurrentCurrency"
	DBRateColumn            = "Rate"

	DefaultListenAddress   = "127.0.0.1:8080"
	DefaultDatabaseAddress = "127.0.0.1:3306"

	ErrorMethodNotAllowed = "ErrorMethodNotAllowed"
	ErrorForbidden        = "ErrorForbidden"
	ErrorNotFound         = "ErrorNotFound"
)

var DateAPIPattern = regexp.MustCompile(`^\/rates[\/]\d{4}\-(0?[1-9]|1[012])\-(0?[1-9]|[12][0-9]|3[01])$`)
