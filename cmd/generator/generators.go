package main

import (
	"context"
	"fmt"
	"github.com/icrowley/fake"
	"github.com/jackc/pgx"
	"log"
	"math"
	"math/rand"
	"time"
)

func generatePeople(conn *pgx.Conn, count int) (err error) {
	_, err = conn.Prepare("insert_people",
		`insert into people (name, lastname, birthday, some_flag) values ($1, $2, $3, $4)`)
	if err != nil {
		return err
	}

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	batch := tx.BeginBatch()

	var (
		name     string
		lastname string
		birthday string
		someFlag int
	)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := 0; i < count; i++ {
		name = fake.FirstName()
		lastname = fake.LastName()
		birthday = fmt.Sprintf("%d-%02d-%02d",
			fake.Year(1970, 2000),
			fake.MonthNum(),
			r.Intn(28)+1)
		someFlag = r.Intn(math.MaxInt32)
		batch.Queue("insert_people",
			[]interface{}{name, lastname, birthday, someFlag},
			nil,
			nil)
	}

	err = batch.Send(context.Background(), nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = batch.Close()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func generateCountries(conn *pgx.Conn, count int) (countryIDs []int, err error) {
	_, err = conn.Prepare("insert_country",
		`insert into countries (name) values ($1) returning id`)
	if err != nil {
		return nil, err
	}

	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}

	batch := tx.BeginBatch()
	var name string
	for i := 0; i < count; i++ {
		name = fake.Country()
		batch.Queue("insert_country", []interface{}{name}, nil, []int16{pgx.BinaryFormatCode})
	}

	err = batch.Send(context.Background(), nil)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	countryIDs = make([]int, count)
	for i := 0; i < count; i++ {
		batch.QueryRowResults().Scan(&countryIDs[i])
	}

	err = batch.Close()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return countryIDs, nil
}

func generateCities(conn *pgx.Conn, count int, countryIDs []int) (err error) {
	_, err = conn.Prepare("insert_city",
		`insert into cities (name, country_id) values ($1, $2)`)
	if err != nil {
		log.Fatalln(err)
	}

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	batch := tx.BeginBatch()

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	var (
		name      string
		countryID int
	)
	for i := 0; i < count; i++ {
		name = fake.City()
		countryID = countryIDs[r.Intn(len(countryIDs))]
		batch.Queue("insert_city", []interface{}{name, countryID}, nil, nil)
	}

	err = batch.Send(context.Background(), nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = batch.Close()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
