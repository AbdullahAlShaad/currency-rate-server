package main

import (
	"flag"
	"github.com/shaad7/currency-rate-server/api"
	"github.com/shaad7/currency-rate-server/controller"
	"github.com/shaad7/currency-rate-server/util"
	"log"
	"os"
)

func main() {
	listenAddress := flag.String("listenaddress", util.DefaultListenAddress, "server address")

	// Database Configuration
	dbAddress := flag.String("databaseaddress", util.DefaultDatabaseAddress, "database address")
	flag.Parse()

	databaseName := util.DatabaseName
	dbPassword := os.Getenv("DB_PASSWORD")
	util.CreateMySqlDatabase(*dbAddress, dbPassword, databaseName)

	controller, err := controller.NewController(*dbAddress, dbPassword, databaseName)
	if err != nil {
		log.Fatalln(err)
	}

	//Initialize Database with source values
	err = controller.EnsureDataLoadedToDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	// Start the server
	rateServer := api.NewServer(*listenAddress, controller)
	log.Fatal(rateServer.Start())
}
