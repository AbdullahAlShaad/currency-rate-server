# currency-rate-server


## Database information

A `mysql` instance needs to be running to run the api server. The server does database operations using mysql user `root`. A database named `currency_exchange` to 
do database operations and `test_currency_exchange` to run the tests is created by the server. Database password is required to connect to `mysql`. Let's export the password
for root user as environmental variable.

```bash
export DB_PASSWORD="pqLUQuIRmZorQBSCY"
```

## Run server

Navigate to the project directory and run the following command : 
```bash
 make build-and-run
```
The server will listen to default listen address and write data to default mysql database address which are `localhost`. Default values are given below : 
```
listenaddress = 127.0.0.1:8080
databaseaddress = 127.0.0.1:3306
```

To override default listen address and/or default mysql database address pass the addresses using argument.
```bash
# make build-and-run ARGS="-listenaddress $LISTEN_ADDRESS -databaseaddress $DATABASE_ADDRESS"
make build-and-run ARGS="-listenaddress 127.0.0.1:4000 -databaseaddress 127.0.0.1:3306"
```

## Endpoints

Now you can perfrom GET verbs on /rates/latest, rates/analyze and rates/{yyyy-mm-dd} endpoints . Sample API test using `curl`
```bash
$ curl http://localhost:8080/rates/analyze
{"base":"EUR","rates_analyze":{"AUD":{"avg":1.5555093749999997,"max":1.5972,"min":1.5289},"BGN":{"avg":1.9557999999999973,"max":1.9558,"min":1.9558},"BRL":{"avg":5.544087500000002,"max":5.7758,"min":5.07},"CAD":{"avg":1.4249140625,"max":1.4616,"min":1.3464},"CHF":{"avg":0.9895109375,"max":1.0056,"min":0.9751},"CNY":{"avg":7.3516359374999976,"max":7.5326,"min":7.2045},"CZK":{"avg":24.159593750000003,"max":24.399,"min":23.725},"DKK":{"avg":7.438075000000004,"max":7.4443,"min":7.4364},"GBP":{"avg":0.8747398437499997,"max":0.89289,"min":0.85715},"HKD":{"avg":8.263484375,"max":8.6183,"min":7.8128},"HRK":{"avg":7.544556410256411,"max":7.5563,"min":7.5365},"HUF":{"avg":402.53281250000003,"max":417.66,"min":386.58},"IDR":{"avg":16351.0853125,"max":16765.93,"min":15615.6},"ILS":{"avg":3.6530625000000003,"max":3.7786,"min":3.5255},"INR":{"avg":86.7145828125,"max":90.3015,"min":81.3058},"ISK":{"avg":151.128125,"max":157.1,"min":145.7},"JPY":{"avg":142.57875000000007,"max":146.82,"min":137.93},"KRW":{"avg":1362.1921874999996,"max":1397.42,"min":1330.37},"MXN":{"avg":20.410542187499995,"max":21.0634,"min":19.4395},"MYR":{"avg":4.67601875,"max":4.77,"min":4.605},"NOK":{"avg":10.536426562499999,"max":10.9783,"min":10.2495},"NZD":{"avg":1.6781187499999999,"max":1.7033,"min":1.6446},"PHP":{"avg":58.985468749999995,"max":59.678,"min":57.793},"PLN":{"avg":4.692374999999999,"max":4.7195,"min":4.6423},"RON":{"avg":4.920570312499998,"max":4.9495,"min":4.8818},"SEK":{"avg":11.04753125,"max":11.3587,"min":10.7241},"SGD":{"avg":1.4268359374999997,"max":1.4413,"min":1.3963},"THB":{"avg":36.51207812499998,"max":37.423,"min":35.455},"TRY":{"avg":19.781896875000005,"max":20.6766,"min":18.51},"USD":{"avg":1.0578906249999998,"max":1.0988,"min":0.9954},"ZAR":{"avg":18.230370312499993,"max":18.9223,"min":17.5768}}}

$ curl http://localhost:8080/rates/latest
{"base":"EUR","rates":{"AUD":1.5499,"BGN":1.9558,"BRL":5.5414,"CAD":1.4616,"CHF":0.9989,"CNY":7.3689,"CZK":23.725,"DKK":7.4443,"GBP":0.8925,"HKD":8.5802,"HRK":7.5365,"HUF":386.58,"IDR":16312.7,"ILS":3.7207,"INR":89.592,"ISK":153.7,"JPY":140.45,"KRW":1346.17,"MXN":20.4625,"MYR":4.657,"NOK":10.9783,"NZD":1.6886,"PHP":58.721,"PLN":4.692,"RON":4.902,"SEK":11.3323,"SGD":1.4331,"THB":36.114,"TRY":20.5806,"USD":1.0937,"ZAR":18.7624}}

$ curl http://localhost:8080/rates/2023-02-03
{"base":"EUR","rates":{"AUD":1.5499,"BGN":1.9558,"BRL":5.5414,"CAD":1.4616,"CHF":0.9989,"CNY":7.3689,"CZK":23.725,"DKK":7.4443,"GBP":0.8925,"HKD":8.5802,"HUF":386.58,"IDR":16312.7,"ILS":3.7207,"INR":89.592,"ISK":153.7,"JPY":140.45,"KRW":1346.17,"MXN":20.4625,"MYR":4.657,"NOK":10.9783,"NZD":1.6886,"PHP":58.721,"PLN":4.692,"RON":4.902,"SEK":11.3323,"SGD":1.4331,"THB":36.114,"TRY":20.5806,"USD":1.0937,"ZAR":18.7624}}

$ curl http://localhost:8080/rates/
{"status":500,"error":"ErrorForbidden"}
```

## Tests

To run the tests, please ensure that the mysql database address should be `127.0.0.1:3306`. 
Run the following command to perform tests : 
```bash
$ go test api/server.go api/server_test.go
ok  	command-line-arguments	0.019s
```