package main

import (
	"flag"
	"github.com/jackc/pgx"
	"log"
)

func main() {
	configPath := flag.String("c", "./config.yml", "config file")
	flag.Parse()

	if configPath == nil || *configPath == "" {
		log.Fatalln("Empty config path")
	}

	config, err := getConfig(*configPath)
	if err != nil {
		log.Fatalln("Config open error:", err)
	}

	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     config.Host,
		Port:     config.Port,
		Database: config.DBName,
		User:     config.User,
		Password: config.Password,
	})
	if err != nil {
		log.Fatalln("Unable to connect to database:", err)
	}
	defer conn.Close()

	err = generatePeople(conn, config.PeopleCount)
	if err != nil {
		log.Fatalln("People generation error:", err)
	}

	ids, err := generateCountries(conn, config.CountriesCount)
	if err != nil {
		log.Fatalln("Countries generation error:", err)
	}

	err = generateCities(conn, config.CitiesCount, ids)
	if err != nil {
		log.Fatalln("Cities generation error:", err)
	}
}
