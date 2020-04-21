package database

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

var username = ""
var password = ""
var host = ""
var port = "5432"
var dbName = ""

func connectToDB() *sql.DB {
	connStr := "postgres://" + username + ":" + password + "@" + host + ":" + port + "/" + dbName

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	return db
}

func createTable(db *sql.DB) {
	createSQLStatement := `CREATE TABLE IF NOT EXISTS playerstats (
		Name string,
		Height string,
		Weight string,
		Team string,
		Age numeric,
		Salary numeric,
		Points numeric,
		Blocks numeric,
		Steals numeric,
		Assists numeric,
		Rebounds numeric,
		FT numeric,
		FTA numeric,
		FG3 numeric,
		FG3A numeric,
		FG numeric,
		FGA numeric,
		MP numeric,
		G numeric,
		PER numeric,
		OWS numeric,
		DWS numeric,
		WS numeric,
		WS48 numeric,
		USG numeric,
		BPM numeric,
		VORP numeric
	)`

	_, err := db.Exec(createSQLStatement)
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()
}

func insertData(db *sql.DB, filename string) {
	nbaCSVFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(bufio.NewReader(nbaCSVFile))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}

		insertSQLStatement := `INSERT INTO playerstats (
			Name,
			Height,
			Weight,
			Team,
			Age,
			Salary,
			Points,
			Blocks,
			Steals,
			Assists,
			Rebounds,
			FT,
			FTA,
			FG3,
			FG3,
			FG,
			FGA,
			MP,
			G,
			PER,
			OWS,
			DWS,
			WS,
			WS48,
			USG,
			BPM,
			VORP
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)
		RETURNING id`
		id := 0

		err = db.QueryRow(insertSQLStatement, line[0], line[1], line[2], line[3], line[4], line[5], line[6], line[7], line[8], line[9], line[10], line[11], line[12], line[13],
			line[14], line[15], line[16], line[17], line[18], line[19], line[20], line[21], line[22], line[23], line[24], line[25], line[26], line[27]).Scan(&id)
		if err != nil {
			fmt.Println(err)
		}
	}
}
