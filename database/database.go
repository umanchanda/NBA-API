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

func connStr() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	if port == "" {
		port = "5432"
	}
	return "postgres://" + username + ":" + password + "@" + host + ":" + port + "/" + dbName + "?sslmode=require"
}

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

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sql.DB) error {
	createSQLStatement := `CREATE TABLE IF NOT EXISTS playerstats (
		id SERIAL PRIMARY KEY,
		Name TEXT,
		Height TEXT,
		Weight TEXT,
		Team TEXT,
		Age TEXT,
		Salary TEXT,
		Points TEXT,
		Blocks TEXT,
		Steals TEXT,
		Assists TEXT,
		Rebounds TEXT,
		FT TEXT,
		FTA TEXT,
		FG3 TEXT,
		FG3A TEXT,
		FG TEXT,
		FGA TEXT,
		MP TEXT,
		G TEXT,
		PER TEXT,
		OWS TEXT,
		DWS TEXT,
		WS TEXT,
		WS48 TEXT,
		USG TEXT,
		BPM TEXT,
		VORP TEXT
	)`

	_, err := db.Exec(createSQLStatement)
	return err
}

func InsertData(db *sql.DB, filename string) error {
	nbaCSVFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer nbaCSVFile.Close()

	insertSQLStatement := `INSERT INTO playerstats (
		Name, Height, Weight, Team, Age, Salary, Points, Blocks, Steals,
		Assists, Rebounds, FT, FTA, FG3, FG3A, FG, FGA, MP, G,
		PER, OWS, DWS, WS, WS48, USG, BPM, VORP
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
		$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27
	)`

	reader := csv.NewReader(bufio.NewReader(nbaCSVFile))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		_, err = db.Exec(insertSQLStatement,
			line[0], line[1], line[2], line[3], line[4], line[5], line[6],
			line[7], line[8], line[9], line[10], line[11], line[12], line[13],
			line[14], line[15], line[16], line[17], line[18], line[19], line[20],
			line[21], line[22], line[23], line[24], line[25], line[26],
		)
		if err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}
	}
	return nil
}

func DatabaseFunctions() error {
	db, err := ConnectToDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := CreateTable(db); err != nil {
		return err
	}
	return InsertData(db, "database/nbastats2018-2019.csv")
}
