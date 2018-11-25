package main

import (
	"flag"
	"github.com/jackc/pgx"
	"log"
	"runtime"
)

func main() {
	threadsCount := flag.Int("t", 1, "threads count")
	configPath := flag.String("c", "./config.yml", "config file")
	flag.Parse()

	if configPath == nil || *configPath == "" {
		log.Fatalln("Empty config path")
	}

	runtime.GOMAXPROCS(*threadsCount)
	log.Println("Using", *threadsCount, "threads")

	config, err := getConfig(*configPath)
	if err != nil {
		log.Fatalln("Config open error:", err)
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
		log.Fatalln("Unable to connect to database:", err)
	}
	defer conn.Close()

	workerFunc := func(id int, queries <-chan interface{}) {
		for query := range queries {
			count, err := query.(Query).Export(conn, config.OutputDir)
			if err != nil {
				log.Println(query.(Query).Name, "exporting error:", err)
				continue
			}

			log.Println(query.(Query).Name+": successfully exported", count, "records")
		}
	}

	workerPool := NewWorkerPool(*threadsCount, workerFunc)
	workerPool.Run()

	queries := workerPool.GetInputChannel()
	for _, query := range config.Queries {
		queries <- query
	}

	workerPool.Shutdown()
	workerPool.Wait()
}
