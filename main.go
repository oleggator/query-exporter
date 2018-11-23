package main

import (
	"flag"
	"github.com/jackc/pgx"
	"github.com/oleggator/query-exporter/workerpool"
	"log"
	"runtime"
)

func main() {
	threadsCount := flag.Int("t", 1, "threads count")
	configPath := flag.String("c", "./config.yml", "config file")
	flag.Parse()

	if configPath == nil || *configPath == "" {
		log.Fatalln("empty config path")
	}

	runtime.GOMAXPROCS(*threadsCount)

	config, err := getConfig(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	workerPool := workerpool.NewWorkerPool(*threadsCount, createWorkerFunc(config))
	workerPool.Run()

	queries := workerPool.GetInputChannel()
	for _, query := range config.Queries {
		queries <- query
	}

	workerPool.Shutdown()
	workerPool.Wait()
}

func createWorkerFunc(config *Config) func(int, <-chan interface{}) {
	pgConfig := pgx.ConnConfig{
		Host:     config.Host,
		Port:     config.Port,
		Database: config.DBName,
		User:     config.User,
		Password: config.Password,
	}

	return func(id int, queries <-chan interface{}) {
		conn, err := pgx.Connect(pgConfig)
		if err != nil {
			log.Fatalln("Unable to connection to database:", err)
		}
		defer conn.Close()

		for query := range queries {
			count, err := query.(Query).Export(conn, config.OutputDir)
			if err != nil {
				log.Println(query.(Query).Name+": error:", err)
				continue
			}

			log.Println(query.(Query).Name+": successfuly exported", count, "records")
		}
	}
}
