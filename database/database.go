package database

import (
	"database/sql"
	"fmt"
	"github.com/shaad7/currency-rate-server/util"
)

type Database struct {
	db *sql.DB
}

var dataBase *Database

// GetDatabase ensures singleton pattern and returns created mysql DB
func GetDatabase(dbAddress, dbName string, password string) (*Database, error) {
	if dataBase == nil {
		var err error
		dataBase, err = createDatabase(dbAddress, dbName, password)
		if err != nil {
			return nil, err
		}
	}
	return dataBase, nil
}

// createDatabase takes mysql username and dbName as argument and create a database connection
// And store that database connection in Database struct and return it
func createDatabase(dbAddress string, dbName string, password string) (*Database, error) {

	// dataSourceAddress = $username:$passwod@tcp($dbadress:$port)/mysql_db_name
	dataSourceAddress := fmt.Sprintf("%s:%s@tcp(%s)/%s", util.DefaultMySqlUsername, password, dbAddress, dbName)
	db, err := sql.Open(util.DriverName, dataSourceAddress)
	if err != nil {
		return nil, err
	} else if err = db.Ping(); err != nil {
		panic(err)
	}
	return &Database{
		db: db,
	}, nil
}

// CreateTable functions creates a table in mysql db named ExchangeRate. Sample Table Looks like :
//+------------+--------------+-----------------+---------+
//| Date       | BaseCurrency | CurrentCurrency | Rate    |
//+------------+--------------+-----------------+---------+
//| 2022-11-07 | EUR          | AUD             |  1.5428 |
//| 2022-11-07 | EUR          | BGN             |  1.9558 |
//| 2022-11-07 | EUR          | BRL             |    5.07 |
//| 2022-11-07 | EUR          | CAD             |  1.3464 |
//+------------+--------------+-----------------+---------+

func (d *Database) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS ExchangeRate 
					    (Date DATE NOT NULL,
					     BaseCurrency VARCHAR(10) NOT NULL,
					     CurrentCurrency VARCHAR(10) NOT NULL,
					     Rate DOUBLE NOT NULL,
					    PRIMARY KEY (Date, CurrentCurrency)
				);`
	_, err := d.db.Exec(query)
	return err
}

// InsertRow functions takes a row info of ExchangeRate Table
func (d *Database) InsertRow(date string, currency string, rate float64) error {
	query := fmt.Sprintf(`INSERT IGNORE INTO 
    								ExchangeRate(Date, BaseCurrency, CurrentCurrency, Rate) 
     								VALUE('%s','EUR','%s', '%f'
     							 );`, date, currency, rate)
	_, err := d.db.Exec(query)
	return err
}

type Exchange struct {
	CurrencyName string
	Rate         float64
}

// GetExchangeRate function takes a data and return all the currency names and rate of that data. The rows are sorted according to CurrencyName.
// Then the data is added into a list of Exchange struct instance and returned
// Sample Query output is given below :
//+-----------------+----------+
//| CurrentCurrency | Rate     |
//+-----------------+----------+
//| AUD             |    1.539 |
//| BGN             |   1.9558 |
//| BRL             |   5.5654 |
//| CAD             |   1.4532 |

func (d *Database) GetExchangeRate(date string) ([]Exchange, error) {
	query := fmt.Sprintf(`SELECT CurrentCurrency,Rate 
									FROM ExchangeRate 
									WHERE Date='%s' ORDER BY CurrentCurrency`, date)
	out, err := d.db.Query(query)
	if err != nil {
		panic(err)
	}
	var exchangeList []Exchange
	for out.Next() {
		var cur Exchange
		err := out.Scan(&cur.CurrencyName, &cur.Rate)
		if err != nil {
			return nil, err
		}
		exchangeList = append(exchangeList, cur)
	}
	return exchangeList, nil
}

// GetLatestExchangeRate functions runs a SQL Query that gets all latest rate of all the currency.The rows are sorted according to CurrencyName.
// Then the data is added into a list of Exchange struct instance and returned
// Sample Query output is given below :
//+-----------------+---------+
//| CurrentCurrency | Rate    |
//+-----------------+---------+
//| AUD             |  1.5499 |
//| BGN             |  1.9558 |
//| BRL             |  5.5414 |
//| CAD             |  1.4616 |

func (d *Database) GetLatestExchangeRate() ([]Exchange, error) {
	query := fmt.Sprintf(`SELECT e.CurrentCurrency,e.Rate 
								 FROM ExchangeRate e 
    								INNER JOIN ( 
    									SELECT CurrentCurrency, MAX(Date) as MaxDate 
    									FROM ExchangeRate 
    									GROUP BY CurrentCurrency ) 
    								en ON e.CurrentCurrency = en.CurrentCurrency AND e.date = en.MaxDate ORDER BY e.CurrentCurrency;`)
	out, err := d.db.Query(query)
	if err != nil {
		panic(err)
	}
	var exchangeList []Exchange
	for out.Next() {
		var cur Exchange
		err := out.Scan(&cur.CurrencyName, &cur.Rate)
		if err != nil {
			return nil, err
		}
		exchangeList = append(exchangeList, cur)
	}
	return exchangeList, nil
}

type ExchangeAnalysis struct {
	CurrencyName string
	AverageRate  float64
	MinRate      float64
	MaxRate      float64
}

// GetExchangeAnalysis runs a SQL query that min,max and average for all currency. Then the data is extracted to ExchangeAnalysis struct and
// a list made and returned
// Sample Query output is given below :
// +-----------------+--------------------+-----------+-----------+
// | CurrentCurrency | AVG(Rate)          | MIN(Rate) | Max(Rate) |
// +-----------------+--------------------+-----------+-----------+
// | AUD             | 1.5555093749999997 |    1.5289 |    1.5972 |
// | BGN             | 1.9557999999999973 |    1.9558 |    1.9558 |
// | BRL             |  5.544087500000002 |      5.07 |    5.7758 |
// | CAD             |       1.4249140625 |    1.3464 |    1.4616 |
func (d *Database) GetExchangeAnalysis() ([]ExchangeAnalysis, error) {
	query := fmt.Sprintf(`SELECT CurrentCurrency, AVG(Rate), MIN(Rate), Max(Rate) 
								 FROM ExchangeRate 
								 GROUP BY CurrentCurrency;`)
	out, err := d.db.Query(query)
	if err != nil {
		panic(err)
	}
	var exchangeList []ExchangeAnalysis
	for out.Next() {
		var cur ExchangeAnalysis
		err := out.Scan(&cur.CurrencyName, &cur.AverageRate, &cur.MinRate, &cur.MaxRate)
		if err != nil {
			return nil, err
		}
		exchangeList = append(exchangeList, cur)
	}
	return exchangeList, nil
}
