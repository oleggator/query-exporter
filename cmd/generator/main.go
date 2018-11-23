package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/oleggator/query-exporter/workerpool"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"runtime"
	"strings"
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

	tables := workerPool.GetInputChannel()

	for _, table := range config.Tables {
		tables <- table
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

	return func(id int, tables <-chan interface{}) {
		conn, err := pgx.Connect(pgConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		}
		defer conn.Close()

		for table := range tables {
			fieldNames := make([]string, len(table.(Table).Fields))
			for i, field := range table.(Table).Fields {
				fieldNames[i] = field.Name
			}

			fieldNamesString := strings.Join(fieldNames, ",")

			tx, _ := conn.Begin()
			batch := tx.BeginBatch()
			for i := 0; i < table.(Table).Count; i++ {
				fieldValues := make([]string, len(table.(Table).Fields))
				for i, field := range table.(Table).Fields {
					fieldValues[i] = "'" + strings.Replace(GetFakeValueByTypeName(field.Type), "'", "''", -1) + "'"
				}

				query := fmt.Sprintf("insert into %s (%s) values (%s)",
					table.(Table).Name,
					fieldNamesString,
					strings.Join(fieldValues, ","))

				batch.Queue(query, nil, nil, nil)
			}

			err = batch.Send(context.Background(), nil)
			if err != nil {
				tx.Rollback()
				log.Println(err)
			}
			_, err = batch.ExecResults()
			if err != nil {
				tx.Rollback()
				log.Println(err)
			}

			batch.Close()
			tx.Commit()
		}

	}
}
