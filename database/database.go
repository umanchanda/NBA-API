package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	SeasonTypeRegular  = "regular"
	SeasonTypePlayoffs = "playoffs"
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

// NBAPlayer holds season totals for a single player from basketball-reference.
type NBAPlayer struct {
	Season     string `json:"season,omitempty"`
	SeasonType string `json:"season_type,omitempty"`
	Name       string `json:"name,omitempty"`
	Team       string `json:"team,omitempty"`
	Pos        string `json:"pos,omitempty"`
	Age        string `json:"age,omitempty"`
	G          string `json:"g,omitempty"`
	GS         string `json:"gs,omitempty"`
	MP         string `json:"mp,omitempty"`
	FG         string `json:"fg,omitempty"`
	FGA        string `json:"fga,omitempty"`
	FGPct      string `json:"fg_pct,omitempty"`
	FG3        string `json:"fg3,omitempty"`
	FG3A       string `json:"fg3a,omitempty"`
	FG3Pct     string `json:"fg3_pct,omitempty"`
	FT         string `json:"ft,omitempty"`
	FTA        string `json:"fta,omitempty"`
	FTPct      string `json:"ft_pct,omitempty"`
	ORB        string `json:"orb,omitempty"`
	DRB        string `json:"drb,omitempty"`
	TRB        string `json:"trb,omitempty"`
	AST        string `json:"ast,omitempty"`
	STL        string `json:"stl,omitempty"`
	BLK        string `json:"blk,omitempty"`
	TOV        string `json:"tov,omitempty"`
	PF         string `json:"pf,omitempty"`
	PTS        string `json:"pts,omitempty"`
}

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS playerstats (
		id          SERIAL PRIMARY KEY,
		season      TEXT,
		season_type TEXT,
		name        TEXT,
		team        TEXT,
		pos         TEXT,
		age         TEXT,
		g           TEXT,
		gs          TEXT,
		mp          TEXT,
		fg          TEXT,
		fga         TEXT,
		fg_pct      TEXT,
		fg3         TEXT,
		fg3a        TEXT,
		fg3_pct     TEXT,
		ft          TEXT,
		fta         TEXT,
		ft_pct      TEXT,
		orb         TEXT,
		drb         TEXT,
		trb         TEXT,
		ast         TEXT,
		stl         TEXT,
		blk         TEXT,
		tov         TEXT,
		pf          TEXT,
		pts         TEXT
	)`)
	if err != nil {
		return err
	}

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_playerstats_name ON playerstats (LOWER(name))`,
		`CREATE INDEX IF NOT EXISTS idx_playerstats_season ON playerstats (season)`,
		`CREATE INDEX IF NOT EXISTS idx_playerstats_season_type ON playerstats (season_type)`,
	}
	for _, idx := range indexes {
		if _, err := db.Exec(idx); err != nil {
			return err
		}
	}
	return nil
}

func InsertPlayers(db *sql.DB, players []NBAPlayer) error {
	stmt := `INSERT INTO playerstats (
		season, season_type, name, team, pos, age, g, gs, mp,
		fg, fga, fg_pct, fg3, fg3a, fg3_pct,
		ft, fta, ft_pct, orb, drb, trb,
		ast, stl, blk, tov, pf, pts
	) VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,
		$10,$11,$12,$13,$14,$15,
		$16,$17,$18,$19,$20,$21,
		$22,$23,$24,$25,$26,$27
	)`

	for _, p := range players {
		_, err := db.Exec(stmt,
			p.Season, p.SeasonType, p.Name, p.Team, p.Pos, p.Age, p.G, p.GS, p.MP,
			p.FG, p.FGA, p.FGPct, p.FG3, p.FG3A, p.FG3Pct,
			p.FT, p.FTA, p.FTPct, p.ORB, p.DRB, p.TRB,
			p.AST, p.STL, p.BLK, p.TOV, p.PF, p.PTS,
		)
		if err != nil {
			return fmt.Errorf("insert failed for %s: %w", p.Name, err)
		}
	}
	return nil
}

func SeasonExists(db *sql.DB, year, seasonType string) (bool, error) {
	var count int
	err := db.QueryRow(
		`SELECT COUNT(*) FROM playerstats WHERE season = $1 AND season_type = $2`,
		year, seasonType,
	).Scan(&count)
	return count > 0, err
}
