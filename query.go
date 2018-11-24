package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"os"
	"path"
	"time"
)

// Query struct
type Query struct {
	QueryString   string `yaml:"query"`
	Name          string `yaml:"name"`
	MaxLinesCount int    `yaml:"max_lines"`
}

// QueryExecutor for struct with Query method
type QueryExecutor interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
}

// Export exports query result to csv files
func (query Query) Export(conn QueryExecutor, outputDir string) (int, error) {
	rows, err := conn.Query(query.QueryString)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	fieldDescriptions := rows.FieldDescriptions()
	headers := make([]string, len(fieldDescriptions))
	for i, field := range fieldDescriptions {
		headers[i] = field.Name
	}

	dirPath := path.Join(outputDir, query.Name)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return 0, err
	}

	var csvFile *CSVFile
	fileIndex := 0
	i := 0
	for ; rows.Next(); i++ {
		if i%query.MaxLinesCount == 0 {
			if csvFile != nil {
				err = csvFile.Close()
				if err != nil {
					return i, err
				}
			}

			filePath := path.Join(dirPath, fmt.Sprintf("%03d.csv", fileIndex))
			csvFile, err = NewCSVFile(filePath, headers)
			if err != nil {
				return i, err
			}

			fileIndex++
		}

		rawValues, err := rows.Values()
		if err != nil {
			return i, err
		}

		values := make([]string, len(rawValues))
		for i, rawValue := range rawValues {
			if fieldDescriptions[i].DataTypeName == "timestamp" || fieldDescriptions[i].DataTypeName == "date" {
				values[i] = fmt.Sprintf("%d", rawValue.(time.Time).Unix())
			} else {
				values[i] = rawValue.(string)
			}
		}

		csvFile.Write(values)
		if err != nil {
			return i, err
		}
	}

	csvFile.Close()

	return i, nil
}
