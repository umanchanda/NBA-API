package database

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	_ "github.com/lib/pq"
)

var username = "postgres"
var password = "d7OY7hwTmyIbEO7Ppd2M"
var host = "nba-api.celg0gvjzujb.us-east-1.rds.amazonaws.com"
var port = "5432"
var dbName = "nba-api"

// NBAPlayer contains fields for various nba stats
type NBAPlayer struct {
	Name     string `json:"name,omitempty"`
	Height   string `json:"height,omitempty"`
	Weight   string `json:"weight,omitempty"`
	Team     string `json:"team,omitempty"`
	Age      string `json:"age,omitempty"`
	Salary   string `json:"salary,omitempty"`
	Points   string `json:"points,omitempty"`
	Blocks   string `json:"blocks,omitempty"`
	Steals   string `json:"steals,omitempty"`
	Assists  string `json:"assists,omitempty"`
	Rebounds string `json:"rebounds,omitempty"`
	FT       string `json:"ft,omitempty"`
	FTA      string `json:"fta,omitempty"`
	FG3      string `json:"fg3,omitempty"`
	FG3A     string `json:"fg3a,omitempty"`
	FG       string `json:"fg,omitempty"`
	FGA      string `json:"fga,omitempty"`
	MP       string `json:"mp,omitempty"`
	G        string `json:"g,omitempty"`
	PER      string `json:"per,omitempty"`
	OWS      string `json:"ows,omitempty"`
	DWS      string `json:"dws,omitempty"`
	WS       string `json:"ws,omitempty"`
	WS48     string `json:"ws48,omitempty"`
	USG      string `json:"usg,omitempty"`
	BPM      string `json:"bpm,omitempty"`
	VORP     string `json:"vorp,omitempty"`
}

func ConnectToDB() *sql.DB {
	connStr := "postgres://" + username + ":" + password + "@" + host + ":" + port + "/" + dbName

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	return db
}

func CreateTable(db *sql.DB) {
	createSQLStatement := `CREATE TABLE IF NOT EXISTS playerstats (
		Name string,
		Height string,
		Weight string,
		Team string,
		Age string,
		Salary string,
		Points string,
		Blocks string,
		Steals string,
		Assists string,
		Rebounds string,
		FT string,
		FTA string,
		FG3 string,
		FG3A string,
		FG string,
		FGA string,
		MP string,
		G string,
		PER string,
		OWS string,
		DWS string,
		WS string,
		WS48 string,
		USG string,
		BPM string,
		VORP string
	)`

	_, err := db.Exec(createSQLStatement)
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()
}

func InsertData(db *sql.DB, filename string) {
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

func DatabaseFunctions() {
	db := ConnectToDB()
	CreateTable(db)
	InsertData(db, "database/nbastats2018-2019.csv")
}