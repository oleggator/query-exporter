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

	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     config.Host,
			Port:     config.Port,
			Database: config.DBName,
			User:     config.User,
			Password: config.Password,
		},
		MaxConnections: 10,
	})
	if err != nil {
		log.Fatalln("Unable to connection to database:", err)
	}
	defer conn.Close()

	workerFunc := func(id int, queries <-chan interface{}) {
		for query := range queries {
			count, err := query.(Query).Export(conn, config.OutputDir)
			if err != nil {
				log.Println(query.(Query).Name+": error:", err)
				continue
			}

			log.Println(query.(Query).Name+": successfuly exported", count, "records")
		}
	}

	workerPool := workerpool.NewWorkerPool(*threadsCount, workerFunc)
	workerPool.Run()

	queries := workerPool.GetInputChannel()
	for _, query := range config.Queries {
		queries <- query
	}

	workerPool.Shutdown()
	workerPool.Wait()
}
