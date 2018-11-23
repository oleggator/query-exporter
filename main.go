package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/oleggator/query-exporter/workerpool"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

func main() {
	threadsCount := flag.Int("t", 1, "threads count")
	configPath := flag.String("c", "", "config file")
	flag.Parse()

	if configPath == nil || *configPath == "" {
		log.Fatal("empty config path")
	}

	runtime.GOMAXPROCS(*threadsCount)

	configFile, err := os.Open(*configPath)
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	configDecoder := yaml.NewDecoder(bufio.NewReader(configFile))
	config := &Config{}
	err = configDecoder.Decode(config)
	if err != nil {
		log.Fatal(err)
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
			fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		}
		defer conn.Close()

		for query := range queries {
			rows, err := conn.Query(query.(Query).QueryString)
			if err != nil {
				log.Println(err)
			}

			fieldDescriptions := rows.FieldDescriptions()
			headers := make([]string, len(fieldDescriptions))
			for i, field := range fieldDescriptions {
				headers[i] = field.Name
			}

			dirPath := path.Join(config.OutputDir, query.(Query).Name)
			err = os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				log.Println(err)
			}

			fileIndex := 0
			filePath := path.Join(dirPath, fmt.Sprintf("%03d.csv", fileIndex))
			csvFile, err := NewCSVFile(filePath, headers)
			if err != nil {
				log.Println(err)
			}

			for i := 0; rows.Next(); i++ {
				rawValues, err := rows.Values()
				if err != nil {
					log.Println(err)
				}

				values := make([]string, len(rawValues))
				for i, rawValue := range rawValues {
					if fieldDescriptions[i].DataTypeName == "timestamp" {
						values[i] = fmt.Sprintf("%d", rawValue.(time.Time).Unix())
					} else {
						values[i] = rawValue.(string)
					}
				}

				csvFile.Write(values)

				if i != 0 && i%query.(Query).MaxLinesCount == 0 {
					fileIndex++

					csvFile.Close()

					filePath := path.Join(dirPath, fmt.Sprintf("%03d.csv", fileIndex))
					csvFile, err = NewCSVFile(filePath, headers)
					if err != nil {
						log.Println(err)
					}
				}
			}

			csvFile.Close()
		}
	}
}
