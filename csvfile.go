package main

import (
	"encoding/csv"
	"os"
)

type CSVFile struct {
	file      *os.File
	csvWriter *csv.Writer
}

func NewCSVFile(filePath string, headers []string) (*CSVFile, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	csvWriter := csv.NewWriter(file)
	err = csvWriter.Write(headers)
	if err != nil {
		return nil, err
	}

	return &CSVFile{
		file:      file,
		csvWriter: csvWriter,
	}, nil
}

func (csvFile *CSVFile) Write(record []string) error {
	return csvFile.csvWriter.Write(record)
}

func (csvFile *CSVFile) Close() (err error) {
	csvFile.csvWriter.Flush()

	err = csvFile.file.Sync()
	if err != nil {
		return err
	}

	err = csvFile.file.Close()
	if err != nil {
		return err
	}

	return nil
}
